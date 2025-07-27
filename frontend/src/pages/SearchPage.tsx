import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import api, { News } from '../services/api'
import NewsList from '../components/NewsList'

const SearchPage = () => {
  const [searchParams] = useSearchParams()
  const [news, setNews] = useState<News[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const query = searchParams.get('q') || ''

  useEffect(() => {
    const fetchNews = async () => {
      if (!query) return
      
      try {
        setLoading(true)
        const data = await api.searchNews(query)
        setNews(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An unknown error occurred')
      } finally {
        setLoading(false)
      }
    }

    fetchNews()
  }, [query])

  if (loading) {
    return <div className="loading">Searching...</div>
  }

  if (error) {
    return <div className="error">Error: {error}</div>
  }

  return (
    <div className="search-page">
      <h1>Search Results for "{query}"</h1>
      {news.length > 0 ? (
        <NewsList news={news} />
      ) : (
        <p>No results found.</p>
      )}
    </div>
  )
}

export default SearchPage