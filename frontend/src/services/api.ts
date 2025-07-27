const API_BASE_URL = '/api'

export interface News {
  id: string
  title: string
  summary: string
  content: string
  image: string
  source: string
  date: string
  url: string
}

export interface User {
  id: number
  username: string
}

export interface SearchParams {
  q?: string
  source?: string
}

class ApiService {
  async getNews(params?: SearchParams): Promise<News[]> {
    let url = `${API_BASE_URL}/news`
    
    if (params) {
      const queryParams = new URLSearchParams()
      if (params.q) queryParams.append('q', params.q)
      if (params.source) queryParams.append('source', params.source)
      const queryString = queryParams.toString()
      if (queryString) url += `?${queryString}`
    }
    
    const response = await fetch(url)
    if (!response.ok) {
      throw new Error('Failed to fetch news')
    }
    return response.json()
  }

  async getNewsById(id: string): Promise<News> {
    const response = await fetch(`${API_BASE_URL}/news/${id}`)
    if (!response.ok) {
      throw new Error('Failed to fetch news')
    }
    return response.json()
  }

  async searchNews(query: string): Promise<News[]> {
    const response = await fetch(`${API_BASE_URL}/search?q=${encodeURIComponent(query)}`)
    if (!response.ok) {
      throw new Error('Failed to search news')
    }
    return response.json()
  }

  async getSources(): Promise<string[]> {
    const response = await fetch(`${API_BASE_URL}/sources`)
    if (!response.ok) {
      throw new Error('Failed to fetch sources')
    }
    return response.json()
  }

  async registerUser(username: string, password: string): Promise<User> {
    const response = await fetch(`${API_BASE_URL}/users/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to register user')
    }
    
    return response.json()
  }

  async loginUser(username: string, password: string): Promise<User> {
    const response = await fetch(`${API_BASE_URL}/users/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to login')
    }
    
    return response.json()
  }

  async addBookmark(articleId: string): Promise<void> {
    // 在实际实现中，需要提供认证信息
    const response = await fetch(`${API_BASE_URL}/protected/bookmarks`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        // 'Authorization': `Bearer ${token}` // 在实际实现中需要添加认证
      },
      body: JSON.stringify({ article_id: articleId }),
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to add bookmark')
    }
  }

  async removeBookmark(articleId: string): Promise<void> {
    // 在实际实现中，需要提供认证信息
    const response = await fetch(`${API_BASE_URL}/protected/bookmarks`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        // 'Authorization': `Bearer ${token}` // 在实际实现中需要添加认证
      },
      body: JSON.stringify({ article_id: articleId }),
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to remove bookmark')
    }
  }

  async getUserBookmarks(): Promise<News[]> {
    // 在实际实现中，需要提供认证信息
    const response = await fetch(`${API_BASE_URL}/protected/bookmarks`, {
      headers: {
        // 'Authorization': `Bearer ${token}` // 在实际实现中需要添加认证
      },
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to fetch bookmarks')
    }
    
    return response.json()
  }
}

export default new ApiService()