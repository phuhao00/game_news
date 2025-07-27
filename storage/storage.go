package storage

import (
	"sync"
	"time"
	"game-news/scraper"
)

// ArticleWithContent 扩展文章结构以包含详细内容
type ArticleWithContent struct {
	scraper.Article
	Content string
}

// Storage handles in-memory storage of news articles
type Storage struct {
	articles []ArticleWithContent
	mu       sync.RWMutex
}

// NewStorage creates a new Storage instance
func NewStorage() *Storage {
	return &Storage{
		articles: make([]ArticleWithContent, 0),
	}
}

// AddArticle adds a new article to storage
func (s *Storage) AddArticle(article scraper.Article, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Check if article already exists
	for i, existing := range s.articles {
		if existing.ID == article.ID {
			// Update existing article
			s.articles[i] = ArticleWithContent{
				Article: article,
				Content: content,
			}
			return
		}
	}
	
	// Add new article
	s.articles = append(s.articles, ArticleWithContent{
		Article: article,
		Content: content,
	})
}

// AddArticles adds multiple articles to storage
func (s *Storage) AddArticles(articles []scraper.Article, scraper *scraper.Scraper) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for _, article := range articles {
		// Check if article already exists
		found := false
		for i, existing := range s.articles {
			if existing.ID == article.ID {
				// Get content for existing article
				content, _ := scraper.ScrapeGameDetails(article.URL)
				
				// Update existing article
				s.articles[i] = ArticleWithContent{
					Article: article,
					Content: content,
				}
				found = true
				break
			}
		}
		
		// Add new article if not found
		if !found {
			content, _ := scraper.ScrapeGameDetails(article.URL)
			s.articles = append(s.articles, ArticleWithContent{
				Article: article,
				Content: content,
			})
		}
	}
}

// GetArticles returns all articles
func (s *Storage) GetArticles() []ArticleWithContent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Return a copy to prevent modifications
	articles := make([]ArticleWithContent, len(s.articles))
	copy(articles, s.articles)
	
	return articles
}

// GetArticleByID returns a specific article by ID
func (s *Storage) GetArticleByID(id string) (ArticleWithContent, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, article := range s.articles {
		if article.ID == id {
			return article, true
		}
	}
	
	return ArticleWithContent{}, false
}

// GetRecentArticles returns the most recent articles
func (s *Storage) GetRecentArticles(limit int) []ArticleWithContent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Sort articles by published date (newest first)
	articles := make([]ArticleWithContent, len(s.articles))
	copy(articles, s.articles)
	
	// Simple bubble sort by date (newest first)
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
	
	return articles
}

// Cleanup removes articles older than the specified duration
func (s *Storage) Cleanup(olderThan time.Duration) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	cutoff := time.Now().Add(-olderThan)
	
	removed := 0
	n := 0
	for _, article := range s.articles {
		if article.PublishedAt.After(cutoff) {
			s.articles[n] = article
			n++
		} else {
			removed++
		}
	}
	
	s.articles = s.articles[:n]
	return removed
}