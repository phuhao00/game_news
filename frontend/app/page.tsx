import { Menu } from 'lucide-react';

export default function Home() {
  return (
    <main className="min-h-screen bg-gray-100">
      <header className="bg-blue-600 text-white p-4 shadow-md">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-2xl font-bold">Game News</h1>
          <button className="p-2 rounded hover:bg-blue-700">
            <Menu size={24} />
          </button>
        </div>
      </header>
      
      <div className="container mx-auto p-4 mt-8">
        <h2 className="text-3xl font-bold mb-6">Latest Gaming News</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {/* Sample news cards */}
          {[1, 2, 3, 4, 5, 6].map((item) => (
            <div key={item} className="bg-white rounded-lg shadow-md overflow-hidden">
              <div className="h-48 bg-gray-300"></div>
              <div className="p-4">
                <h3 className="text-xl font-semibold mb-2">Game Title {item}</h3>
                <p className="text-gray-600 mb-4">
                  Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
                </p>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-500">June 15, 2023</span>
                  <button className="text-blue-600 hover:text-blue-800">Read More</button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </main>
  );
}