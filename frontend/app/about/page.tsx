import MainLayout from '@/components/layout/MainLayout';
import { generateMetadata } from '@/lib/metadata';

export const metadata = generateMetadata({
  title: 'About Us',
  description: 'Learn more about Game News, your premier destination for gaming updates, reviews, and industry insights.',
});

export default function About() {
  return (
    <MainLayout>
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold mb-8">About Game News</h1>
        
        <div className="bg-white rounded-lg shadow-md p-8 mb-8">
          <h2 className="text-2xl font-semibold mb-4">Our Mission</h2>
          <p className="text-gray-700 mb-6">
            Game News is your premier destination for the latest updates, reviews, and insights 
            from the gaming world. We're passionate about bringing you comprehensive coverage 
            of everything from indie gems to AAA blockbusters.
          </p>
          
          <h2 className="text-2xl font-semibold mb-4">What We Cover</h2>
          <ul className="list-disc list-inside text-gray-700 mb-6 space-y-2">
            <li>Latest game releases and updates</li>
            <li>E-sports tournaments and competitions</li>
            <li>Gaming hardware and technology</li>
            <li>Industry news and developer insights</li>
            <li>Game reviews and recommendations</li>
            <li>Mobile and indie gaming spotlights</li>
          </ul>
          
          <h2 className="text-2xl font-semibold mb-4">Our Team</h2>
          <p className="text-gray-700">
            Our dedicated team of gaming enthusiasts and journalists work around the clock 
            to bring you accurate, timely, and engaging content. We believe in the power 
            of gaming to bring people together and create unforgettable experiences.
          </p>
        </div>
      </div>
    </MainLayout>
  );
}