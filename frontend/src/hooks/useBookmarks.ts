import { useState, useEffect } from 'react'
import api from '../services/api'

interface Bookmark {
  id: string
  title: string
}

const useBookmarks = () => {
  const [bookmarks, setBookmarks] = useState<Bookmark[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchBookmarks = async () => {
      try {
        setLoading(true)
        // In a real implementation, this would fetch actual bookmarks
        // For now, we'll use mock data
        const mockBookmarks = [
          { id: '1', title: 'New Game Update Coming Soon' },
          { id: '2', title: 'Esports Tournament Results Are Out' }
        ]
        setBookmarks(mockBookmarks)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An unknown error occurred')
      } finally {
        setLoading(false)
      }
    }

    fetchBookmarks()
  }, [])

  const addBookmark = async (id: string, title: string) => {
    try {
      // In a real implementation, this would call the API
      await api.addBookmark(id)
      setBookmarks(prev => [...prev, { id, title }])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add bookmark')
      throw err
    }
  }

  const removeBookmark = async (id: string) => {
    try {
      // In a real implementation, this would call the API
      await api.removeBookmark(id)
      setBookmarks(prev => prev.filter(bookmark => bookmark.id !== id))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to remove bookmark')
      throw err
    }
  }

  const isBookmarked = (id: string) => {
    return bookmarks.some(bookmark => bookmark.id === id)
  }

  return {
    bookmarks,
    loading,
    error,
    addBookmark,
    removeBookmark,
    isBookmarked
  }
}

export default useBookmarks