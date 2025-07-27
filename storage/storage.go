package storage

import (
	"database/sql"
	"sync"
	"time"
	"game-news/scraper"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// ArticleWithContent 扩展文章结构以包含详细内容
type ArticleWithContent struct {
	scraper.Article
	Content string
}

// Storage handles storage of news articles
type Storage struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewStorage creates a new Storage instance
func NewStorage() *Storage {
	// 首先尝试连接PostgreSQL
	db, err := sql.Open("pq", "host=db port=5432 user=game_news password=game_news_password dbname=game_news sslmode=disable")
	if err != nil {
		// 如果PostgreSQL连接失败，回退到SQLite
		db, err = sql.Open("sqlite3", "./game_news.db")
		if err != nil {
			panic(err)
		}
	} else {
		// 测试PostgreSQL连接
		if err := db.Ping(); err != nil {
			// 如果连接测试失败，回退到SQLite
			db.Close()
			db, err = sql.Open("sqlite3", "./game_news.db")
			if err != nil {
				panic(err)
			}
		}
	}
	
	// 创建表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS articles (
		id TEXT PRIMARY KEY,
		title TEXT,
		url TEXT,
		image_url TEXT,
		summary TEXT,
		source TEXT,
		published_at TIMESTAMP,
		content TEXT
	);
	
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE,
		password_hash TEXT,
		created_at TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS bookmarks (
		id SERIAL PRIMARY KEY,
		user_id INTEGER,
		article_id TEXT,
		created_at TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (article_id) REFERENCES articles(id)
	);
	
	CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles(published_at);
	CREATE INDEX IF NOT EXISTS idx_articles_source ON articles(source);
	CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);
	CREATE INDEX IF NOT EXISTS idx_bookmarks_article_id ON bookmarks(article_id);
	`
	
	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
	
	return &Storage{
		db: db,
	}
}

// AddArticle adds a new article to storage
func (s *Storage) AddArticle(article scraper.Article, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	insertSQL := `
	INSERT INTO articles 
	(id, title, url, image_url, summary, source, published_at, content)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (id) DO UPDATE SET
	title = EXCLUDED.title,
	url = EXCLUDED.url,
	image_url = EXCLUDED.image_url,
	summary = EXCLUDED.summary,
	source = EXCLUDED.source,
	published_at = EXCLUDED.published_at,
	content = EXCLUDED.content
	`
	
	_, err := s.db.Exec(insertSQL, article.ID, article.Title, article.URL, article.ImageURL, article.Summary, article.Source, article.PublishedAt, content)
	return err
}

// AddArticles adds multiple articles to storage
func (s *Storage) AddArticles(articles []scraper.Article, scraper *scraper.Scraper) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	
	insertSQL := `
	INSERT INTO articles 
	(id, title, url, image_url, summary, source, published_at, content)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (id) DO UPDATE SET
	title = EXCLUDED.title,
	url = EXCLUDED.url,
	image_url = EXCLUDED.image_url,
	summary = EXCLUDED.summary,
	source = EXCLUDED.source,
	published_at = EXCLUDED.published_at,
	content = EXCLUDED.content
	`
	
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	
	for _, article := range articles {
		content, _ := scraper.ScrapeGameDetails(article.URL)
		_, err := stmt.Exec(article.ID, article.Title, article.URL, article.ImageURL, article.Summary, article.Source, article.PublishedAt, content)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	
	return tx.Commit()
}

// GetArticles returns all articles
func (s *Storage) GetArticles() ([]ArticleWithContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	query := `
	SELECT id, title, url, image_url, summary, source, published_at, content
	FROM articles
	ORDER BY published_at DESC
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	articles := make([]ArticleWithContent, 0)
	for rows.Next() {
		var article ArticleWithContent
		var publishedAt time.Time
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.URL,
			&article.ImageURL,
			&article.Summary,
			&article.Source,
			&publishedAt,
			&article.Content,
		)
		if err != nil {
			return nil, err
		}
		
		article.PublishedAt = publishedAt
		articles = append(articles, article)
	}
	
	return articles, nil
}

// GetArticleByID returns a specific article by ID
func (s *Storage) GetArticleByID(id string) (ArticleWithContent, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	query := `
	SELECT id, title, url, image_url, summary, source, published_at, content
	FROM articles
	WHERE id = $1
	`
	
	var article ArticleWithContent
	var publishedAt time.Time
	
	err := s.db.QueryRow(query, id).Scan(
		&article.ID,
		&article.Title,
		&article.URL,
		&article.ImageURL,
		&article.Summary,
		&article.Source,
		&publishedAt,
		&article.Content,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return article, false, nil
		}
		return article, false, err
	}
	
	article.PublishedAt = publishedAt
	return article, true, nil
}

// SearchArticles 搜索文章
func (s *Storage) SearchArticles(query string) ([]ArticleWithContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	searchSQL := `
	SELECT id, title, url, image_url, summary, source, published_at, content
	FROM articles
	WHERE title ILIKE $1 OR summary ILIKE $1 OR content ILIKE $1
	ORDER BY published_at DESC
	`
	
	searchTerm := "%" + query + "%"
	rows, err := s.db.Query(searchSQL, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	articles := make([]ArticleWithContent, 0)
	for rows.Next() {
		var article ArticleWithContent
		var publishedAt time.Time
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.URL,
			&article.ImageURL,
			&article.Summary,
			&article.Source,
			&publishedAt,
			&article.Content,
		)
		if err != nil {
			return nil, err
		}
		
		article.PublishedAt = publishedAt
		articles = append(articles, article)
	}
	
	return articles, nil
}

// FilterArticlesBySource 按来源过滤文章
func (s *Storage) FilterArticlesBySource(source string) ([]ArticleWithContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	query := `
	SELECT id, title, url, image_url, summary, source, published_at, content
	FROM articles
	WHERE source = $1
	ORDER BY published_at DESC
	`
	
	rows, err := s.db.Query(query, source)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	articles := make([]ArticleWithContent, 0)
	for rows.Next() {
		var article ArticleWithContent
		var publishedAt time.Time
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.URL,
			&article.ImageURL,
			&article.Summary,
			&article.Source,
			&publishedAt,
			&article.Content,
		)
		if err != nil {
			return nil, err
		}
		
		article.PublishedAt = publishedAt
		articles = append(articles, article)
	}
	
	return articles, nil
}

// GetRecentArticles returns the most recent articles
func (s *Storage) GetRecentArticles(limit int) ([]ArticleWithContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	query := `
	SELECT id, title, url, image_url, summary, source, published_at, content
	FROM articles
	ORDER BY published_at DESC
	`
	
	if limit > 0 {
		query += " LIMIT $1"
	}
	
	var rows *sql.Rows
	var err error
	
	if limit > 0 {
		rows, err = s.db.Query(query, limit)
	} else {
		rows, err = s.db.Query(query)
	}
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	articles := make([]ArticleWithContent, 0)
	for rows.Next() {
		var article ArticleWithContent
		var publishedAt time.Time
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.URL,
			&article.ImageURL,
			&article.Summary,
			&article.Source,
			&publishedAt,
			&article.Content,
		)
		if err != nil {
			return nil, err
		}
		
		article.PublishedAt = publishedAt
		articles = append(articles, article)
	}
	
	return articles, nil
}

// Cleanup removes articles older than the specified duration
func (s *Storage) Cleanup(olderThan time.Duration) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	cutoff := time.Now().Add(-olderThan)
	
	result, err := s.db.Exec("DELETE FROM articles WHERE published_at < $1", cutoff)
	if err != nil {
		return 0, err
	}
	
	return result.RowsAffected()
}

// CreateUser 创建新用户
func (s *Storage) CreateUser(username, passwordHash string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	var id int64
	err := s.db.QueryRow(
		"INSERT INTO users (username, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id",
		username, passwordHash, time.Now(),
	).Scan(&id)
	
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *Storage) GetUserByUsername(username string) (int64, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var id int64
	var passwordHash string
	
	err := s.db.QueryRow(
		"SELECT id, password_hash FROM users WHERE username = $1",
		username,
	).Scan(&id, &passwordHash)
	
	if err != nil {
		return 0, "", err
	}
	
	return id, passwordHash, nil
}

// AddBookmark 添加书签
func (s *Storage) AddBookmark(userID int64, articleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	_, err := s.db.Exec(
		"INSERT INTO bookmarks (user_id, article_id, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		userID, articleID, time.Now(),
	)
	
	return err
}

// RemoveBookmark 删除书签
func (s *Storage) RemoveBookmark(userID int64, articleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	_, err := s.db.Exec(
		"DELETE FROM bookmarks WHERE user_id = $1 AND article_id = $2",
		userID, articleID,
	)
	
	return err
}

// GetBookmarks 获取用户书签
func (s *Storage) GetBookmarks(userID int64) ([]ArticleWithContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	query := `
	SELECT a.id, a.title, a.url, a.image_url, a.summary, a.source, a.published_at, a.content
	FROM articles a
	JOIN bookmarks b ON a.id = b.article_id
	WHERE b.user_id = $1
	ORDER BY b.created_at DESC
	`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	articles := make([]ArticleWithContent, 0)
	for rows.Next() {
		var article ArticleWithContent
		var publishedAt time.Time
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.URL,
			&article.ImageURL,
			&article.Summary,
			&article.Source,
			&publishedAt,
			&article.Content,
		)
		if err != nil {
			return nil, err
		}
		
		article.PublishedAt = publishedAt
		articles = append(articles, article)
	}
	
	return articles, nil
}