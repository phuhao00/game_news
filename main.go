package main

import (
	"net/http"
	"time"
	"game-news/scraper"
	"game-news/storage"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

// News 结构体定义新闻数据结构，用于API响应
type News struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Content string `json:"content"`
	Image   string `json:"image"`
	Source  string `json:"source"`
	Date    string `json:"date"`
	URL     string `json:"url"`
}

func main() {
	// 创建存储实例
	store := storage.NewStorage()
	
	// 创建爬虫实例
	scraper := scraper.NewScraper()
	
	// 初始抓取新闻
	if articles, err := scraper.ScrapeGames(); err == nil {
		store.AddArticles(articles, scraper)
	}
	
	// 定期抓取新闻 (每小时一次)
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		
		for range ticker.C {
			if articles, err := scraper.ScrapeGames(); err == nil {
				store.AddArticles(articles, scraper)
				
				// 清理超过7天的旧新闻
				store.Cleanup(7 * 24 * time.Hour)
			}
		}
	}()
	
	// 设置Gin运行模式
	gin.SetMode(gin.ReleaseMode)
	
	// 创建Gin引擎实例
	router := gin.Default()
	
	// 配置CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080"}, // Vite默认开发端口和生产端口
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	// 设置静态文件目录（React构建后的文件）
	router.Static("/static", "./dist/static")
	router.StaticFile("/", "./dist/index.html")
	router.StaticFile("/about", "./dist/index.html")
	router.StaticFile("/news/*id", "./dist/index.html")
	
	// 设置API路由
	api := router.Group("/api")
	{
		api.GET("/news", getNews(store))
		api.GET("/news/:id", getNewsByID(store))
	}
	
	// 启动服务器
	router.Run(":8080")
}

// getNews 返回所有新闻
func getNews(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		articles := store.GetRecentArticles(0) // 0 means no limit
		
		// 转换为API响应格式
		newsList := make([]News, len(articles))
		for i, article := range articles {
			newsList[i] = News{
				ID:      article.ID,
				Title:   article.Title,
				Summary: article.Summary,
				Content: "", // 在列表中不包含完整内容以减少数据传输
				Image:   article.ImageURL,
				Source:  article.Source,
				Date:    article.PublishedAt.Format("2006-01-02"),
				URL:     article.URL,
			}
		}
		
		c.JSON(http.StatusOK, newsList)
	}
}

// getNewsByID 根据ID返回特定新闻
func getNewsByID(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		article, found := store.GetArticleByID(id)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
			return
		}
		
		// 转换为API响应格式
		news := News{
			ID:      article.ID,
			Title:   article.Title,
			Summary: article.Summary,
			Content: article.Content,
			Image:   article.ImageURL,
			Source:  article.Source,
			Date:    article.PublishedAt.Format("2006-01-02"),
			URL:     article.URL,
		}
		
		c.JSON(http.StatusOK, news)
	}
}