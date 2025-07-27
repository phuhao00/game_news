import { useState, useEffect } from 'react'
import NewsList from '../components/NewsList'
import { News } from '../types/news'

const HomePage = () => {
  const [news, setNews] = useState<News[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchNews = async () => {
      try {
        // 模拟从API获取数据
        const response = await fetch('/api/news')
        if (!response.ok) {
          throw new Error('Failed to fetch news')
        }
        const data = await response.json()
        setNews(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An unknown error occurred')
        // 使用模拟数据
        setNews([
          {
            id: '1',
            title: 'New Game Update Coming Soon',
            summary: 'Developers announce major update with new features and improvements.',
            content: 'Developers today officially announced that the highly anticipated game update will be released next month. This update will include brand new maps, characters, and gameplay mechanics, promising to deliver a completely new gaming experience. The development team said they spent over a year perfecting these new features and conducted multiple rounds of testing to ensure game balance.',
            image: 'https://picsum.photos/600/400?random=1',
            source: 'GameNews Network',
            date: '2023-07-15',
            url: 'https://example.com/news/1'
          },
          {
            id: '2',
            title: 'Esports Tournament Results Are Out',
            summary: "The year's biggest esports tournament has ended, with the champion team winning a million-dollar prize.",
            content: 'After a week of intense competition, the annual esports tournament has finally come to a close. The Korean team won the championship with a 3:2 score against their opponents, taking home the title and a million-dollar prize. This competition attracted over 10 million viewers worldwide, making it one of the highest-rated esports events in history.',
            image: 'https://picsum.photos/600/400?random=2',
            source: 'eSports Daily',
            date: '2023-07-10',
            url: 'https://example.com/news/2'
          },
          {
            id: '3',
            title: 'Indie Game Sensation Gains Popularity',
            summary: 'An indie game developed by a small team goes viral, selling over 500,000 copies.',
            content: 'A small team of only five developers has seen their indie game "Journey of Dreams" sell over 500,000 copies within a month of its release, far exceeding expectations. The game has received unanimous praise from both players and critics for its unique art style and innovative gameplay mechanics. The development team said they initially just wanted to create a game that could express their creativity, never expecting such tremendous success.',
            image: 'https://picsum.photos/600/400?random=3',
            source: 'Indie Game Watch',
            date: '2023-07-05',
            url: 'https://example.com/news/3'
          }
        ])
      } finally {
        setLoading(false)
      }
    }

    fetchNews()
  }, [])

  if (loading) {
    return <div className="loading">Loading...</div>
  }

  if (error) {
    return <div className="error">Error: {error}</div>
  }

  return (
    <div className="home-page">
      <h1>Latest Game News</h1>
      <NewsList news={news} />
    </div>
  )
}

export default HomePage