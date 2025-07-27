package storage

import (
	"database/sql"
	"sync"
	"time"
	"game-news/scraper"
	"log"
	"os"
	"github.com/lib/pq"
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
	var db *sql.DB
	var err error
	
	// 检查环境变量以确定使用哪种数据库
	dbHost := os.Getenv("DB_HOST")
	if dbHost != "" {
		// 使用PostgreSQL
		dbUser := getEnvOrDefault("DB_USER", "game_news")
		dbPassword := getEnvOrDefault("DB_PASSWORD", "game_news_password")
		dbName := getEnvOrDefault("DB_NAME", "game_news")
		dbPort := getEnvOrDefault("DB_PORT", "5432")
		
		connStr := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal("Failed to connect to PostgreSQL: ", err)
		}
		
		// 测试连接
		if err := db.Ping(); err != nil {
			log.Fatal("Failed to ping PostgreSQL: ", err)
		}
		log.Println("Connected to PostgreSQL database")
	} else {
		// 使用SQLite
		db, err = sql.Open("sqlite3", "./game_news.db")
		if err != nil {
			log.Fatal("Failed to connect to SQLite: ", err)
		}
		log.Println("Connected to SQLite database")
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
		log.Fatal("Failed to create tables: ", err)
	}
	
	return &Storage{
		db: db,
	}
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// isPostgreSQL 检查是否使用PostgreSQL数据库
func (s *Storage) isPostgreSQL() bool {
	// 检查驱动名称
	return s.db.Driver() == &pq.Driver{}
}

// AddArticle adds a new article to storage
func (s *Storage) AddArticle(article scraper.Article, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	var insertSQL string
	if s.isPostgreSQL() {
		insertSQL = `
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
	} else {
		insertSQL = `
		INSERT OR REPLACE INTO articles 
		(id, title, url, image_url, summary, source, published_at, content)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`
	}
	
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
	
	var insertSQL string
	if s.isPostgreSQL() {
		insertSQL = `
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
	} else {
		insertSQL = `
		INSERT OR REPLACE INTO articles 
		(id, title, url, image_url, summary, source, published_at, content)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`
	}
	
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
	
	// 对于SQLite，需要使用?占位符
	if !s.isPostgreSQL() {
		query = `
		SELECT id, title, url, image_url, summary, source, published_at, content
		FROM articles
		WHERE id = ?
		`
	}
	
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
	
	// 对于SQLite，需要使用LIKE和?占位符
	if !s.isPostgreSQL() {
		searchSQL = `
		SELECT id, title, url, image_url, summary, source, published_at, content
		FROM articles
		WHERE title LIKE ? OR summary LIKE ? OR content LIKE ?
		ORDER BY published_at DESC
		`
	}
	
	searchTerm := "%" + query + "%"
	var rows *sql.Rows
	var err error
	
	if s.isPostgreSQL() {
		rows, err = s.db.Query(searchSQL, searchTerm)
	} else {
		rows, err = s.db.Query(searchSQL, searchTerm, searchTerm, searchTerm)
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
	
	// 对于SQLite，需要使用?占位符
	if !s.isPostgreSQL() {
		query = `
		SELECT id, title, url, image_url, summary, source, published_at, content
		FROM articles
		WHERE source = ?
		ORDER BY published_at DESC
		`
	}
	
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
		if s.isPostgreSQL() {
			query += " LIMIT $1"
		} else {
			query += " LIMIT ?"
		}
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
	
	var result sql.Result
	var err error
	
	if s.isPostgreSQL() {
		result, err = s.db.Exec("DELETE FROM articles WHERE published_at < $1", cutoff)
	} else {
		result, err = s.db.Exec("DELETE FROM articles WHERE published_at < ?", cutoff)
	}
	
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
	var err error
	
	if s.isPostgreSQL() {
		err = s.db.QueryRow(
			"INSERT INTO users (username, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id",
			username, passwordHash, time.Now(),
		).Scan(&id)
	} else {
		result, err := s.db.Exec(
			"INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)",
			username, passwordHash, time.Now(),
		)
		if err != nil {
			return 0, err
		}
		id, err = result.LastInsertId()
	}
	
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
	
	var err error
	if s.isPostgreSQL() {
		err = s.db.QueryRow(
			"SELECT id, password_hash FROM users WHERE username = $1",
			username,
		).Scan(&id, &passwordHash)
	} else {
		err = s.db.QueryRow(
			"SELECT id, password_hash FROM users WHERE username = ?",
			username,
		).Scan(&id, &passwordHash)
	}
	
	if err != nil {
		return 0, "", err
	}
	
	return id, passwordHash, nil
}

// AddBookmark 添加书签
func (s *Storage) AddBookmark(userID int64, articleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	var err error
	if s.isPostgreSQL() {
		_, err = s.db.Exec(
			"INSERT INTO bookmarks (user_id, article_id, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			userID, articleID, time.Now(),
		)
	} else {
		_, err = s.db.Exec(
			"INSERT OR IGNORE INTO bookmarks (user_id, article_id, created_at) VALUES (?, ?, ?)",
			userID, articleID, time.Now(),
		)
	}
	
	return err
}

// RemoveBookmark 删除书签
func (s *Storage) RemoveBookmark(userID int64, articleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	var err error
	if s.isPostgreSQL() {
		_, err = s.db.Exec(
			"DELETE FROM bookmarks WHERE user_id = $1 AND article_id = $2",
			userID, articleID,
		)
	} else {
		_, err = s.db.Exec(
			"DELETE FROM bookmarks WHERE user_id = ? AND article_id = ?",
			userID, articleID,
		)
	}
	
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
	
	// 对于SQLite，需要使用?占位符
	if !s.isPostgreSQL() {
		query = `
		SELECT a.id, a.title, a.url, a.image_url, a.summary, a.source, a.published_at, a.content
		FROM articles a
		JOIN bookmarks b ON a.id = b.article_id
		WHERE b.user_id = ?
		ORDER BY b.created_at DESC
		`
	}
	
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