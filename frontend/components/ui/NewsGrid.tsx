import NewsCard from './NewsCard';
import { NewsItem } from '@/lib/data';

interface NewsGridProps {
  items: NewsItem[];
}

export default function NewsGrid({ items }: NewsGridProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {items.map((item) => (
        <NewsCard
          key={item.id}
          id={item.id}
          title={item.title}
          description={item.description}
          date={item.date}
        />
      ))}
    </div>
  );
}