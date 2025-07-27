import { Link } from 'react-router-dom'

const Header = () => {
  return (
    <header className="header">
      <div className="container">
        <Link to="/" className="logo">
          GameNews
        </Link>
        <nav className="nav">
          <Link to="/" className="nav-link">Home</Link>
          <Link to="/about" className="nav-link">About</Link>
        </nav>
      </div>
    </header>
  )
}

export default Header