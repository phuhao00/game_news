import { Link } from 'react-router-dom'
import { News } from '../types/news'

interface NewsItemProps {
  news: News
}

const NewsItem = ({ news }: NewsItemProps) => {
  return (
    <div className="news-item">
      <img src={news.image} alt={news.title} className="news-image" />
      <div className="news-content">
        <h2 className="news-title">{news.title}</h2>
        <p className="news-summary">{news.summary}</p>
        <div className="news-meta">
          <span className="news-source">{news.source}</span>
          <span className="news-date">{news.date}</span>
        </div>
        <div className="news-actions">
          <Link to={`/news/${news.id}`} className="read-more">
            Read more
          </Link>
          {news.url && (
            <a 
              href={news.url} 
              target="_blank" 
              rel="noopener noreferrer"
              className="read-original"
            >
              Read original
            </a>
          )}
        </div>
      </div>
    </div>
  )
}

export default NewsItem