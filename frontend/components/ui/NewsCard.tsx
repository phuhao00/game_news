import Link from 'next/link';

interface NewsCardProps {
  id: number;
  title: string;
  description: string;
  date: string;
}

export default function NewsCard({ id, title, description, date }: NewsCardProps) {
  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden">
      <div className="h-48 bg-gray-300"></div>
      <div className="p-4">
        <h3 className="text-xl font-semibold mb-2">{title}</h3>
        <p className="text-gray-600 mb-4">{description}</p>
        <div className="flex justify-between items-center">
          <span className="text-sm text-gray-500">{date}</span>
          <Link href={`/news/${id}`} className="text-blue-600 hover:text-blue-800">
            Read More
          </Link>
        </div>
      </div>
    </div>
  );
}
