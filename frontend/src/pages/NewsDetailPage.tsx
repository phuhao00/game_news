import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { News } from '../types/news'

const NewsDetailPage = () => {
  const { id } = useParams<{ id: string }>()
  const [news, setNews] = useState<News | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchNews = async () => {
      try {
        // 模拟从API获取数据
        const response = await fetch(`/api/news/${id}`)
        if (!response.ok) {
          throw new Error('Failed to fetch news')
        }
        const data = await response.json()
        setNews(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An unknown error occurred')
      } finally {
        setLoading(false)
      }
    }

    if (id) {
      fetchNews()
    }
  }, [id])

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
        {news.url && (
          <div className="news-link">
            <a href={news.url} target="_blank" rel="noopener noreferrer">
              Read original article
            </a>
          </div>
        )}
      </article>
    </div>
  )
}

export default NewsDetailPage