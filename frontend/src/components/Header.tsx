import { useState } from 'react'
import { Link, useNavigate, useLocation } from 'react-router-dom'

const Header = () => {
  const [searchQuery, setSearchQuery] = useState('')
  const navigate = useNavigate()
  const location = useLocation()

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      navigate(`/search?q=${encodeURIComponent(searchQuery.trim())}`)
    }
  }

  // In a real implementation, you would check if the user is authenticated
  const isAuthenticated = false // Placeholder

  return (
    <header className="header">
      <div className="container">
        <Link to="/" className="logo">
          GameNews
        </Link>
        
        <form onSubmit={handleSearch} className="search-form">
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search news..."
            className="search-input"
          />
          <button type="submit" className="search-button">
            Search
          </button>
        </form>
        
        <nav className="nav">
          <Link 
            to="/" 
            className={`nav-link ${location.pathname === '/' ? 'active' : ''}`}
          >
            Home
          </Link>
          <Link 
            to="/bookmarks" 
            className={`nav-link ${location.pathname === '/bookmarks' ? 'active' : ''}`}
          >
            Bookmarks
          </Link>
          <Link 
            to="/about" 
            className={`nav-link ${location.pathname === '/about' ? 'active' : ''}`}
          >
            About
          </Link>
          {isAuthenticated ? (
            <Link 
              to="/profile" 
              className={`nav-link ${location.pathname === '/profile' ? 'active' : ''}`}
            >
              Profile
            </Link>
          ) : (
            <Link 
              to="/auth" 
              className={`nav-link ${location.pathname === '/auth' ? 'active' : ''}`}
            >
              Login
            </Link>
          )}
        </nav>
      </div>
    </header>
  )
}

export default Header