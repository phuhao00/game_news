package main

import (
	"net/http"
	"time"
	"game-news/scraper"
	"game-news/storage"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"golang.org/x/crypto/bcrypt"
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

// User 结构体定义用户数据结构
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// BookmarkRequest 书签请求结构体
type BookmarkRequest struct {
	ArticleID string `json:"article_id"`
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
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	// 设置静态文件目录（React构建后的文件）
	router.Static("/static", "./dist/static")
	router.StaticFile("/", "./dist/index.html")
	router.StaticFile("/about", "./dist/index.html")
	router.StaticFile("/news/*id", "./dist/index.html")
	router.StaticFile("/search", "./dist/index.html")
	router.StaticFile("/bookmarks", "./dist/index.html")
	
	// 设置API路由
	api := router.Group("/api")
	{
		// 公开API路由
		public := api.Group("/")
		{
			public.GET("/news", getNews(store))
			public.GET("/news/:id", getNewsByID(store))
			public.GET("/search", searchNews(store))
			public.GET("/sources", getSources(store))
			public.POST("/users/register", registerUser(store))
			public.POST("/users/login", loginUser(store))
		}
		
		// 需要认证的路由
		protected := api.Group("/protected")
		protected.Use(authMiddleware(store))
		{
			protected.POST("/bookmarks", addBookmark(store))
			protected.DELETE("/bookmarks", removeBookmark(store))
			protected.GET("/bookmarks", getUserBookmarks(store))
		}
	}
	
	// 启动服务器
	router.Run(":8080")
}

// getNews 返回所有新闻
func getNews(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取查询参数
		source := c.Query("source")
		
		var articles []storage.ArticleWithContent
		var err error
		
		if source != "" {
			// 按来源过滤
			articles, err = store.FilterArticlesBySource(source)
		} else {
			// 获取所有文章
			articles, err = store.GetRecentArticles(20) // 限制20篇文章
		}
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
			return
		}
		
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
		
		article, found, err := store.GetArticleByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
			return
		}
		
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

// searchNews 搜索新闻
func searchNews(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
			return
		}
		
		articles, err := store.SearchArticles(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search news"})
			return
		}
		
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

// getSources 获取所有新闻来源
func getSources(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里应该从数据库查询所有不同的来源
		sources := []string{"GameSpot", "IGN", "GameNews Network", "eSports Daily", "Indie Game Watch"}
		c.JSON(http.StatusOK, sources)
	}
}

// registerUser 用户注册
func registerUser(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// 检查用户名是否已存在
		_, _, err := store.GetUserByUsername(user.Username)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		
		// Hash密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		
		// 创建用户
		userID, err := store.CreateUser(user.Username, string(hashedPassword))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		
		c.JSON(http.StatusOK, User{ID: userID, Username: user.Username})
	}
}

// loginUser 用户登录
func loginUser(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// 获取用户
		userID, hashedPassword, err := store.GetUserByUsername(user.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		
		// 验证密码
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		
		// 在实际应用中，这里应该生成JWT令牌
		// 为简化起见，我们直接返回用户信息
		c.JSON(http.StatusOK, User{ID: userID, Username: user.Username})
	}
}

// authMiddleware 认证中间件
func authMiddleware(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在实际应用中，这里应该验证JWT令牌
		// 为简化起见，我们检查是否存在用户ID头部
		userID := c.GetHeader("User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			c.Abort()
			return
		}
		
		c.Set("user_id", userID)
		c.Next()
	}
}

// addBookmark 添加书签
func addBookmark(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)
		
		var req BookmarkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// 在实际应用中，需要将字符串ID转换为整数
		// 为简化起见，我们假设用户ID是有效的
		c.JSON(http.StatusOK, gin.H{"message": "Bookmark added"})
	}
}

// removeBookmark 删除书签
func removeBookmark(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)
		
		var req BookmarkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// 在实际应用中，需要将字符串ID转换为整数
		// 为简化起见，我们假设用户ID是有效的
		c.JSON(http.StatusOK, gin.H{"message": "Bookmark removed"})
	}
}

// getUserBookmarks 获取用户书签
func getUserBookmarks(store *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在实际应用中，需要从上下文中获取用户ID
		// 为简化起见，我们返回空列表
		
		articles := make([]storage.ArticleWithContent, 0)
		
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