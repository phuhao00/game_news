import Link from 'next/link';
import Button from './Button';
import { formatDate, truncateText } from '@/lib/utils';

interface NewsCardProps {
  id: number;
  title: string;
  description: string;
  date: string;
}

export default function NewsCard({ id, title, description, date }: NewsCardProps) {
  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300">
      <div className="h-48 bg-gradient-to-br from-gray-200 to-gray-300 relative">
        <div className="absolute inset-0 bg-black bg-opacity-20 flex items-center justify-center">
          <span className="text-white text-sm font-medium">Game Image</span>
        </div>
      </div>
      <div className="p-4">
        <h3 className="text-xl font-semibold mb-2 line-clamp-2">
          {truncateText(title, 60)}
        </h3>
        <p className="text-gray-600 mb-4 line-clamp-3">
          {truncateText(description, 120)}
        </p>
        <div className="flex justify-between items-center">
          <span className="text-sm text-gray-500">{formatDate(date)}</span>
          <Link href={`/news/${id}`}>
            <Button variant="outline" size="sm">
              Read More
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
