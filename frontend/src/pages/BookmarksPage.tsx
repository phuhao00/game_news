import { useState, useEffect } from 'react'
import api, { News } from '../services/api'
import NewsList from '../components/NewsList'

const BookmarksPage = () => {
  const [news, setNews] = useState<News[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchBookmarks = async () => {
      try {
        setLoading(true)
        // 在实际实现中，需要先登录并获取认证令牌
        const data = await api.getUserBookmarks()
        setNews(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An unknown error occurred')
      } finally {
        setLoading(false)
      }
    }

    fetchBookmarks()
  }, [])

  if (loading) {
    return <div className="loading">Loading bookmarks...</div>
  }

  if (error) {
    return <div className="error">Error: {error}</div>
  }

  return (
    <div className="bookmarks-page">
      <h1>Your Bookmarks</h1>
      {news.length > 0 ? (
        <NewsList news={news} />
      ) : (
        <p>You haven't bookmarked any articles yet.</p>
      )}
    </div>
  )
}

export default BookmarksPage