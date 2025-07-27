import { News } from '../types/news'
import NewsItem from './NewsItem'

interface NewsListProps {
  news: News[]
}

const NewsList = ({ news }: NewsListProps) => {
  return (
    <div className="news-list">
      {news.map((item) => (
        <NewsItem key={item.id} news={item} />
      ))}
    </div>
  )
}

export default NewsList