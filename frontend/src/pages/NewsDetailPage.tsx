import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import api, { News } from '../services/api'

const NewsDetailPage = () => {
  const { id } = useParams<{ id: string }>()
  const [news, setNews] = useState<News | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isBookmarked, setIsBookmarked] = useState(false)

  useEffect(() => {
    const fetchNews = async () => {
      try {
        if (!id) return
        
        setLoading(true)
        const data = await api.getNewsById(id)
        setNews(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An unknown error occurred')
      } finally {
        setLoading(false)
      }
    }

    fetchNews()
  }, [id])

  const toggleBookmark = async () => {
    if (!news) return
    
    try {
      if (isBookmarked) {
        await api.removeBookmark(news.id)
      } else {
        await api.addBookmark(news.id)
      }
      setIsBookmarked(!isBookmarked)
    } catch (err) {
      console.error('Failed to toggle bookmark:', err)
    }
  }

  if (loading) {
    return <div className="loading">Loading...</div>
  }

  if (error) {
    return <div className="error">Error: {error}</div>
  }

  if (!news) {
    return <div className="not-found">News not found</div>
  }

  return (
    <div className="news-detail-page">
      <article className="news-detail">
        <h1>{news.title}</h1>
        <div className="news-meta">
          <span className="news-source">{news.source}</span>
          <span className="news-date">{news.date}</span>
        </div>
        <img src={news.image} alt={news.title} className="news-image" />
        <div className="news-content">
          {news.content}
        </div>
        <div className="news-actions">
          {news.url && (
            <a 
              href={news.url} 
              target="_blank" 
              rel="noopener noreferrer"
              className="read-original"
            >
              Read original article
            </a>
          )}
          <button 
            onClick={toggleBookmark}
            className={`bookmark-button ${isBookmarked ? 'bookmarked' : ''}`}
          >
            {isBookmarked ? 'Remove Bookmark' : 'Bookmark Article'}
          </button>
        </div>
      </article>
    </div>
  )
}

export default NewsDetailPage