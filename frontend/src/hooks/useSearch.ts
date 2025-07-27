import { useState, useEffect, useCallback } from 'react'
import api, { News } from '../services/api'

const useSearch = () => {
  const [results, setResults] = useState<News[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [query, setQuery] = useState('')

  const search = useCallback(async (searchQuery: string) => {
    if (!searchQuery.trim()) {
      setResults([])
      setQuery('')
      return
    }

    try {
      setLoading(true)
      setError(null)
      setQuery(searchQuery)
      const data = await api.searchNews(searchQuery)
      setResults(data)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Search failed'
      setError(errorMessage)
      setResults([])
    } finally {
      setLoading(false)
    }
  }, [])

  const clearResults = useCallback(() => {
    setResults([])
    setQuery('')
    setError(null)
  }, [])

  return {
    results,
    loading,
    error,
    query,
    search,
    clearResults
  }
}

export default useSearch