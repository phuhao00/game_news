package storage

import (
	"database/sql"
	"sync"
	"time"
	"game-news/scraper"
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
	// 初始化SQLite数据库
	db, err := sql.Open("sqlite3", "./game_news.db")
	if err != nil {
		panic(err)
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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password_hash TEXT,
		created_at TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		article_id TEXT,
		created_at TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (article_id) REFERENCES articles(id)
	);
	
	CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles(published_at);
	CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);
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
	INSERT OR REPLACE INTO articles 
	(id, title, url, image_url, summary, source, published_at, content)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
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
	INSERT OR REPLACE INTO articles 
	(id, title, url, image_url, summary, source, published_at, content)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
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
		var publishedAt []byte
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
		
		// 解析时间
		if t, err := time.Parse("2006-01-02 15:04:05", string(publishedAt)); err == nil {
			article.PublishedAt = t
		}
		
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
	WHERE id = ?
	`
	
	var article ArticleWithContent
	var publishedAt []byte
	
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
	
	// 解析时间
	if t, err := time.Parse("2006-01-02 15:04:05", string(publishedAt)); err == nil {
		article.PublishedAt = t
	}
	
	return article, true, nil
}

// SearchArticles 搜索文章
func (s *Storage) SearchArticles(query string) ([]ArticleWithContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	searchSQL := `
	SELECT id, title, url, image_url, summary, source, published_at, content
	FROM articles
	WHERE title LIKE ? OR summary LIKE ? OR content LIKE ?
	ORDER BY published_at DESC
	`
	
	searchTerm := "%" + query + "%"
	rows, err := s.db.Query(searchSQL, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	articles := make([]ArticleWithContent, 0)
	for rows.Next() {
		var article ArticleWithContent
		var publishedAt []byte
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
		
		// 解析时间
		if t, err := time.Parse("2006-01-02 15:04:05", string(publishedAt)); err == nil {
			article.PublishedAt = t
		}
		
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
	WHERE source = ?
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
		var publishedAt []byte
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
		
		// 解析时间
		if t, err := time.Parse("2006-01-02 15:04:05", string(publishedAt)); err == nil {
			article.PublishedAt = t
		}
		
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
		query += " LIMIT ?"
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
		var publishedAt []byte
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
		
		// 解析时间
		if t, err := time.Parse("2006-01-02 15:04:05", string(publishedAt)); err == nil {
			article.PublishedAt = t
		}
		
		articles = append(articles, article)
	}
	
	return articles, nil
}

// Cleanup removes articles older than the specified duration
func (s *Storage) Cleanup(olderThan time.Duration) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	cutoff := time.Now().Add(-olderThan)
	
	result, err := s.db.Exec("DELETE FROM articles WHERE published_at < ?", cutoff)
	if err != nil {
		return 0, err
	}
	
	return result.RowsAffected()
}

// CreateUser 创建新用户
func (s *Storage) CreateUser(username, passwordHash string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	result, err := s.db.Exec(
		"INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)",
		username, passwordHash, time.Now(),
	)
	if err != nil {
		return 0, err
	}
	
	return result.LastInsertId()
}

// GetUserByUsername 根据用户名获取用户
func (s *Storage) GetUserByUsername(username string) (int64, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var id int64
	var passwordHash string
	
	err := s.db.QueryRow(
		"SELECT id, password_hash FROM users WHERE username = ?",
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
		"INSERT OR IGNORE INTO bookmarks (user_id, article_id, created_at) VALUES (?, ?, ?)",
		userID, articleID, time.Now(),
	)
	
	return err
}

// RemoveBookmark 删除书签
func (s *Storage) RemoveBookmark(userID int64, articleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	_, err := s.db.Exec(
		"DELETE FROM bookmarks WHERE user_id = ? AND article_id = ?",
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
	WHERE b.user_id = ?
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
		var publishedAt []byte
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
		
		// 解析时间
		if t, err := time.Parse("2006-01-02 15:04:05", string(publishedAt)); err == nil {
			article.PublishedAt = t
		}
		
		articles = append(articles, article)
	}
	
	return articles, nil
}