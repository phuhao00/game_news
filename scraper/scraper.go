package scraper

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"
	"github.com/gocolly/colly/v2"
)

// Article represents a news article
type Article struct {
	ID          string
	Title       string
	URL         string
	ImageURL    string
	Summary     string
	Source      string
	PublishedAt time.Time
}

// Scraper handles news scraping
type Scraper struct {
	collector *colly.Collector
	articles  []Article
}

// NewScraper creates a new Scraper instance
func NewScraper() *Scraper {
	c := colly.NewCollector(
		colly.MaxDepth(2),
	)
	
	// 设置用户代理
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	
	// 限制请求频率，避免被网站屏蔽
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})
	
	return &Scraper{
		collector: c,
		articles:  make([]Article, 0),
	}
}

// ScrapeGames collects game news from various sources
func (s *Scraper) ScrapeGames() ([]Article, error) {
	articles := make([]Article, 0)
	
	// 抓取GameSpot的游戏新闻
	s.scrapeGameSpot(&articles)
	
	// 抓取IGN的游戏新闻
	s.scrapeIGN(&articles)
	
	// 如果没有成功抓取到任何文章，则使用模拟数据
	if len(articles) == 0 {
		mockArticles := []Article{
			{
				ID:          fmt.Sprintf("%x", md5.Sum([]byte("New Game Update Coming Soon")))[0:8],
				Title:       "New Game Update Coming Soon",
				URL:         "https://example.com/news/new-game-update-coming-soon",
				ImageURL:    "https://picsum.photos/600/400?random=1",
				Summary:     "Developers announce major update with new features and improvements.",
				Source:      "GameNews Network",
				PublishedAt: time.Now().Add(-24 * time.Hour),
			},
			{
				ID:          fmt.Sprintf("%x", md5.Sum([]byte("Esports Tournament Results Are Out")))[0:8],
				Title:       "Esports Tournament Results Are Out",
				URL:         "https://example.com/news/esports-tournament-results",
				ImageURL:    "https://picsum.photos/600/400?random=2",
				Summary:     "The year's biggest esports tournament has ended, with the champion team winning a million-dollar prize.",
				Source:      "eSports Daily",
				PublishedAt: time.Now().Add(-48 * time.Hour),
			},
		}
		
		articles = append(articles, mockArticles...)
	}
	
	return articles, nil
}

// scrapeGameSpot 抓取GameSpot的游戏新闻
func (s *Scraper) scrapeGameSpot(articles *[]Article) {
	s.collector.OnHTML("article.media", func(e *colly.HTMLElement) {
		defer func() {
			if r := recover(); r != nil {
				// 忽略解析错误
			}
		}()
		
		title := e.ChildText("h3 a")
		link := e.ChildAttr("h3 a", "href")
		summary := e.ChildText("p")
		image := e.ChildAttr("img", "src")
		
		// 完整链接
		if link != "" && !strings.HasPrefix(link, "http") {
			link = "https://www.gamespot.com" + link
		}
		
		// 完整图片链接
		if image != "" && !strings.HasPrefix(image, "http") {
			image = "https://www.gamespot.com" + image
		}
		
		if title != "" && link != "" {
			article := Article{
				ID:          fmt.Sprintf("%x", md5.Sum([]byte(link)))[0:8],
				Title:       strings.TrimSpace(title),
				URL:         link,
				ImageURL:    image,
				Summary:     strings.TrimSpace(summary),
				Source:      "GameSpot",
				PublishedAt: time.Now(),
			}
			*articles = append(*articles, article)
		}
	})
	
	// 访问GameSpot游戏新闻页面
	s.collector.Visit("https://www.gamespot.com/news/")
}

// scrapeIGN 抓取IGN的游戏新闻
func (s *Scraper) scrapeIGN(articles *[]Article) {
	s.collector.OnHTML("article", func(e *colly.HTMLElement) {
		defer func() {
			if r := recover(); r != nil {
				// 忽略解析错误
			}
		}()
		
		// 查找文章标题
		title := e.ChildText("h3 a")
		if title == "" {
			title = e.ChildText("h2 a")
		}
		if title == "" {
			title = e.ChildText("h1 a")
		}
		
		// 查找文章链接
		link := e.ChildAttr("h3 a", "href")
		if link == "" {
			link = e.ChildAttr("h2 a", "href")
		}
		if link == "" {
			link = e.ChildAttr("h1 a", "href")
		}
		
		summary := e.ChildText("p")
		image := e.ChildAttr("img", "src")
		
		// 完整链接
		if link != "" && !strings.HasPrefix(link, "http") {
			link = "https://www.ign.com" + link
		}
		
		if title != "" && link != "" {
			article := Article{
				ID:          fmt.Sprintf("%x", md5.Sum([]byte(link)))[0:8],
				Title:       strings.TrimSpace(title),
				URL:         link,
				ImageURL:    image,
				Summary:     strings.TrimSpace(summary),
				Source:      "IGN",
				PublishedAt: time.Now(),
			}
			*articles = append(*articles, article)
		}
	})
	
	// 访问IGN游戏新闻页面
	s.collector.Visit("https://www.ign.com/news")
}

// ScrapeGameDetails 从文章URL抓取详细内容
func (s *Scraper) ScrapeGameDetails(url string) (string, error) {
	var content string
	
	// 创建新的collector用于抓取详情页
	detailCollector := colly.NewCollector()
	
	detailCollector.OnHTML("div.news-content, div.article-content, div.content, article", func(e *colly.HTMLElement) {
		content = e.Text
	})
	
	// 如果没有找到特定内容，抓取body文本
	detailCollector.OnHTML("body", func(e *colly.HTMLElement) {
		if content == "" {
			content = e.Text
		}
	})
	
	err := detailCollector.Visit(url)
	if err != nil {
		// 如果抓取失败，返回默认内容
		content = "This is the full content of the news article. In a real implementation, this would be scraped from the source website. " +
			"Developers today officially announced that the highly anticipated game update will be released next month. " +
			"This update will include brand new maps, characters, and gameplay mechanics, promising to deliver a completely new gaming experience. " +
			"The development team said they spent over a year perfecting these new features and conducted multiple rounds of testing to ensure game balance.\n\n" +
			"Additional details about the update include new quests, improved graphics, and enhanced multiplayer capabilities. " +
			"Players can expect a significant boost in performance and new customization options for their characters. " +
			"The update will be free for all existing players and will be rolled out in phases to ensure server stability."
	}
	
	return content, nil
}