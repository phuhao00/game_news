import { notFound } from 'next/navigation';
import Link from 'next/link';
import { ArrowLeft } from 'lucide-react';
import MainLayout from '@/components/layout/MainLayout';
import { newsItems } from '@/lib/data';

export default function NewsPage({ params }: { params: { id: string } }) {
  const newsId = parseInt(params.id);
  const newsItem = newsItems.find(item => item.id === newsId);
  
  if (!newsItem) {
    notFound();
  }
  
  return (
    <MainLayout>
      <Link href="/" className="inline-flex items-center text-blue-600 hover:text-blue-800 mb-6">
        <ArrowLeft size={20} className="mr-2" />
        Back to Home
      </Link>
      
      <article className="bg-white rounded-lg shadow-md overflow-hidden">
        <div className="h-64 bg-gray-300"></div>
        <div className="p-6">
          <h1 className="text-3xl font-bold mb-2">{newsItem.title}</h1>
          <p className="text-gray-500 mb-6">{newsItem.date}</p>
          <div className="prose max-w-none">
            <p className="text-gray-700">{newsItem.content}</p>
          </div>
        </div>
      </article>
    </MainLayout>
  );
}
