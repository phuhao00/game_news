package storage

import (
	"context"
	"strings"
	"sync"
	"time"
	"game-news/scraper"
	"log"
	"os"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

// ArticleWithContent 扩展文章结构以包含详细内容
type ArticleWithContent struct {
	ID          string    `bson:"id"`
	Title       string    `bson:"title"`
	URL         string    `bson:"url"`
	ImageURL    string    `bson:"image_url"`
	Summary     string    `bson:"summary"`
	Source      string    `bson:"source"`
	PublishedAt time.Time `bson:"published_at"`
	Content     string    `bson:"content"`
}

// User represents a user in the system
type User struct {
	ID           int64     `bson:"id"`
	Username     string    `bson:"username"`
	PasswordHash string    `bson:"password_hash"`
	CreatedAt    time.Time `bson:"created_at"`
}

// Bookmark represents a user bookmark
type Bookmark struct {
	ID        int64     `bson:"id"`
	UserID    int64     `bson:"user_id"`
	ArticleID string    `bson:"article_id"`
	CreatedAt time.Time `bson:"created_at"`
}

// Storage handles storage of news articles
type Storage struct {
	client    *mongo.Client
	database  *mongo.Database
	articles  *mongo.Collection
	users     *mongo.Collection
	bookmarks *mongo.Collection
	mu        sync.RWMutex
	
	// In-memory storage for when no database is available
	inMemoryArticles map[string]ArticleWithContent
	inMemoryUsers    map[string]User
	inMemoryBookmarks map[int64][]string
	useInMemory      bool
}

// NewStorage creates a new Storage instance
func NewStorage() *Storage {
	// Default to in-memory storage
	storage := &Storage{
		inMemoryArticles:  make(map[string]ArticleWithContent),
		inMemoryUsers:     make(map[string]User),
		inMemoryBookmarks: make(map[int64][]string),
		useInMemory:       true,
	}
	
	// Check if MongoDB is configured
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Println("MONGO_URI not set, using in-memory storage")
		return storage
	}
	
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v, using in-memory storage", err)
		return storage
	}
	
	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v, using in-memory storage", err)
		return storage
	}
	
	log.Println("Connected to MongoDB")
	
	// Initialize collections
	database := client.Database("game_news")
	storage.client = client
	storage.database = database
	storage.articles = database.Collection("articles")
	storage.users = database.Collection("users")
	storage.bookmarks = database.Collection("bookmarks")
	storage.useInMemory = false
	
	// Create indexes
	storage.createIndexes()
	
	return storage
}

// createIndexes creates necessary indexes for collections
func (s *Storage) createIndexes() {
	if s.useInMemory {
		return
	}
	
	ctx := context.Background()
	
	// Articles indexes
	s.articles.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"id", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{"published_at", -1}},
		},
		{
			Keys: bson.D{{"source", 1}},
		},
	})
	
	// Users indexes
	s.users.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"username", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{"id", 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	
	// Bookmarks indexes
	s.bookmarks.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"user_id", 1}},
		},
		{
			Keys: bson.D{{"article_id", 1}},
		},
	})
}

// AddArticle adds a new article to storage
func (s *Storage) AddArticle(article scraper.Article, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// If using in-memory storage
	if s.useInMemory {
		articleWithContent := ArticleWithContent{
			ID:          article.ID,
			Title:       article.Title,
			URL:         article.URL,
			ImageURL:    article.ImageURL,
			Summary:     article.Summary,
			Source:      article.Source,
			PublishedAt: article.PublishedAt,
			Content:     content,
		}
		s.inMemoryArticles[article.ID] = articleWithContent
		return nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	articleWithContent := ArticleWithContent{
		ID:          article.ID,
		Title:       article.Title,
		URL:         article.URL,
		ImageURL:    article.ImageURL,
		Summary:     article.Summary,
		Source:      article.Source,
		PublishedAt: article.PublishedAt,
		Content:     content,
	}
	
	_, err := s.articles.UpdateOne(
		ctx,
		bson.M{"id": article.ID},
		bson.M{"$set": articleWithContent},
		options.Update().SetUpsert(true),
	)
	
	return err
}

// AddArticles adds multiple articles to storage
func (s *Storage) AddArticles(articles []scraper.Article, scraperInstance *scraper.Scraper) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// If using in-memory storage
	if s.useInMemory {
		for _, article := range articles {
			content, _ := scraperInstance.ScrapeGameDetails(article.URL)
			articleWithContent := ArticleWithContent{
				ID:          article.ID,
				Title:       article.Title,
				URL:         article.URL,
				ImageURL:    article.ImageURL,
				Summary:     article.Summary,
				Source:      article.Source,
				PublishedAt: article.PublishedAt,
				Content:     content,
			}
			s.inMemoryArticles[article.ID] = articleWithContent
		}
		return nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	var models []mongo.WriteModel
	for _, article := range articles {
		content, _ := scraperInstance.ScrapeGameDetails(article.URL)
		articleWithContent := ArticleWithContent{
			ID:          article.ID,
			Title:       article.Title,
			URL:         article.URL,
			ImageURL:    article.ImageURL,
			Summary:     article.Summary,
			Source:      article.Source,
			PublishedAt: article.PublishedAt,
			Content:     content,
		}
		
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"id": article.ID}).
			SetUpdate(bson.M{"$set": articleWithContent}).
			SetUpsert(true)
		
		models = append(models, model)
	}
	
	_, err := s.articles.BulkWrite(ctx, models)
	return err
}

// GetArticles returns all articles
func (s *Storage) GetArticles() ([]ArticleWithContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// If using in-memory storage
	if s.useInMemory {
		articles := make([]ArticleWithContent, 0, len(s.inMemoryArticles))
		for _, article := range s.inMemoryArticles {
			articles = append(articles, article)
		}
		
		// Sort by published date (newest first)
		for i := 0; i < len(articles)-1; i++ {
			for j := 0; j < len(articles)-i-1; j++ {
				if articles[j].PublishedAt.Before(articles[j+1].PublishedAt) {
					articles[j], articles[j+1] = articles[j+1], articles[j]
				}
			}
		}
		
		return articles, nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	cursor, err := s.articles.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{"published_at", -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var articles []ArticleWithContent
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	
	return articles, nil
}

// GetArticleByID returns a specific article by ID
func (s *Storage) GetArticleByID(id string) (ArticleWithContent, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// If using in-memory storage
	if s.useInMemory {
		article, exists := s.inMemoryArticles[id]
		return article, exists, nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	var article ArticleWithContent
	err := s.articles.FindOne(ctx, bson.M{"id": id}).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return article, false, nil
		}
		return article, false, err
	}
	
	return article, true, nil
}

// SearchArticles searches articles by query
func (s *Storage) SearchArticles(query string) ([]ArticleWithContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// If using in-memory storage
	if s.useInMemory {
		articles := make([]ArticleWithContent, 0)
		searchTerm := strings.ToLower(query)
		
		for _, article := range s.inMemoryArticles {
			if strings.Contains(strings.ToLower(article.Title), searchTerm) ||
				strings.Contains(strings.ToLower(article.Summary), searchTerm) ||
				strings.Contains(strings.ToLower(article.Content), searchTerm) {
				articles = append(articles, article)
			}
		}
		
		// Sort by published date (newest first)
		for i := 0; i < len(articles)-1; i++ {
			for j := 0; j < len(articles)-i-1; j++ {
				if articles[j].PublishedAt.Before(articles[j+1].PublishedAt) {
					articles[j], articles[j+1] = articles[j+1], articles[j]
				}
			}
		}
		
		return articles, nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	searchRegex := bson.M{"$regex": query, "$options": "i"}
	filter := bson.M{
		"$or": []bson.M{
			{"title": searchRegex},
			{"summary": searchRegex},
			{"content": searchRegex},
		},
	}
	
	cursor, err := s.articles.Find(ctx, filter, options.Find().SetSort(bson.D{{"published_at", -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var articles []ArticleWithContent
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	
	return articles, nil
}

// FilterArticlesBySource filters articles by source
func (s *Storage) FilterArticlesBySource(source string) ([]ArticleWithContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// If using in-memory storage
	if s.useInMemory {
		articles := make([]ArticleWithContent, 0)
		
		for _, article := range s.inMemoryArticles {
			if article.Source == source {
				articles = append(articles, article)
			}
		}
		
		// Sort by published date (newest first)
		for i := 0; i < len(articles)-1; i++ {
			for j := 0; j < len(articles)-i-1; j++ {
				if articles[j].PublishedAt.Before(articles[j+1].PublishedAt) {
					articles[j], articles[j+1] = articles[j+1], articles[j]
				}
			}
		}
		
		return articles, nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	cursor, err := s.articles.Find(ctx, bson.M{"source": source}, options.Find().SetSort(bson.D{{"published_at", -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var articles []ArticleWithContent
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	
	return articles, nil
}

// GetRecentArticles returns the most recent articles
func (s *Storage) GetRecentArticles(limit int) ([]ArticleWithContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// If using in-memory storage
	if s.useInMemory {
		articles := make([]ArticleWithContent, 0, len(s.inMemoryArticles))
		for _, article := range s.inMemoryArticles {
			articles = append(articles, article)
		}
		
		// Sort by published date (newest first)
		for i := 0; i < len(articles)-1; i++ {
			for j := 0; j < len(articles)-i-1; j++ {
				if articles[j].PublishedAt.Before(articles[j+1].PublishedAt) {
					articles[j], articles[j+1] = articles[j+1], articles[j]
				}
			}
		}
		
		// Limit results
		if limit > 0 && limit < len(articles) {
			articles = articles[:limit]
		}
		
		return articles, nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	findOptions := options.Find().SetSort(bson.D{{"published_at", -1}})
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	
	cursor, err := s.articles.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var articles []ArticleWithContent
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	
	return articles, nil
}

// Cleanup removes articles older than the specified duration
func (s *Storage) Cleanup(olderThan time.Duration) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	cutoff := time.Now().Add(-olderThan)
	
	// If using in-memory storage
	if s.useInMemory {
		count := int64(0)
		for id, article := range s.inMemoryArticles {
			if article.PublishedAt.Before(cutoff) {
				delete(s.inMemoryArticles, id)
				count++
			}
		}
		return count, nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	result, err := s.articles.DeleteMany(ctx, bson.M{"published_at": bson.M{"$lt": cutoff}})
	if err != nil {
		return 0, err
	}
	
	return result.DeletedCount, nil
}

// CreateUser creates a new user
func (s *Storage) CreateUser(username, passwordHash string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// If using in-memory storage
	if s.useInMemory {
		// Simple ID generation for in-memory storage
		id := int64(len(s.inMemoryUsers) + 1)
		user := User{
			ID:           id,
			Username:     username,
			PasswordHash: passwordHash,
			CreatedAt:    time.Now(),
		}
		s.inMemoryUsers[username] = user
		return id, nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	// Find the max ID
	var maxUser User
	err := s.users.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.D{{"id", -1}})).Decode(&maxUser)
	var id int64 = 1
	if err == nil {
		id = maxUser.ID + 1
	}
	
	user := User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
	
	_, err = s.users.InsertOne(ctx, user)
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

// GetUserByUsername gets user by username
func (s *Storage) GetUserByUsername(username string) (int64, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// If using in-memory storage
	if s.useInMemory {
		user, exists := s.inMemoryUsers[username]
		if !exists {
			return 0, "", mongo.ErrNoDocuments // This is just for consistency, we should return a different error
		}
		return user.ID, user.PasswordHash, nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	var user User
	err := s.users.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return 0, "", err
	}
	
	return user.ID, user.PasswordHash, nil
}

// AddBookmark adds a bookmark
func (s *Storage) AddBookmark(userID int64, articleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// If using in-memory storage
	if s.useInMemory {
		// Check if bookmark already exists
		exists := false
		for _, bookmarkedArticleID := range s.inMemoryBookmarks[userID] {
			if bookmarkedArticleID == articleID {
				exists = true
				break
			}
		}
		
		if !exists {
			s.inMemoryBookmarks[userID] = append(s.inMemoryBookmarks[userID], articleID)
		}
		return nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	bookmark := Bookmark{
		UserID:    userID,
		ArticleID: articleID,
		CreatedAt: time.Now(),
	}
	
	// Use upsert to avoid duplicates
	_, err := s.bookmarks.UpdateOne(
		ctx,
		bson.M{"user_id": userID, "article_id": articleID},
		bson.M{"$set": bookmark},
		options.Update().SetUpsert(true),
	)
	
	return err
}

// RemoveBookmark removes a bookmark
func (s *Storage) RemoveBookmark(userID int64, articleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// If using in-memory storage
	if s.useInMemory {
		bookmarks := s.inMemoryBookmarks[userID]
		for i, bookmarkedArticleID := range bookmarks {
			if bookmarkedArticleID == articleID {
				// Remove the bookmark
				s.inMemoryBookmarks[userID] = append(bookmarks[:i], bookmarks[i+1:]...)
				break
			}
		}
		return nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	_, err := s.bookmarks.DeleteOne(ctx, bson.M{"user_id": userID, "article_id": articleID})
	return err
}

// GetBookmarks gets user bookmarks
func (s *Storage) GetBookmarks(userID int64) ([]ArticleWithContent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// If using in-memory storage
	if s.useInMemory {
		articleIDs := s.inMemoryBookmarks[userID]
		articles := make([]ArticleWithContent, 0, len(articleIDs))
		
		for _, articleID := range articleIDs {
			if article, exists := s.inMemoryArticles[articleID]; exists {
				articles = append(articles, article)
			}
		}
		
		// Sort by bookmark creation date (newest first)
		// Since we don't store bookmark creation date in in-memory storage,
		// we'll just return them in the order they were added
		return articles, nil
	}
	
	// Use MongoDB
	ctx := context.Background()
	
	// First, get the bookmarked article IDs
	cursor, err := s.bookmarks.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var bookmarks []Bookmark
	if err = cursor.All(ctx, &bookmarks); err != nil {
		return nil, err
	}
	
	// Extract article IDs
	articleIDs := make([]string, len(bookmarks))
	for i, bookmark := range bookmarks {
		articleIDs[i] = bookmark.ArticleID
	}
	
	// Then get the articles
	filter := bson.M{"id": bson.M{"$in": articleIDs}}
	cursor, err = s.articles.Find(ctx, filter, options.Find().SetSort(bson.D{{"published_at", -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var articles []ArticleWithContent
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	
	return articles, nil
}