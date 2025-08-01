import MainLayout from '@/components/layout/MainLayout';
import NewsGrid from '@/components/ui/NewsGrid';
import { newsItems } from '@/lib/data';

export default function Home() {
  return (
    <MainLayout>
      <h2 className="text-3xl font-bold mb-6">Latest Gaming News</h2>
      <NewsGrid items={newsItems} />
    </MainLayout>
  );
}
