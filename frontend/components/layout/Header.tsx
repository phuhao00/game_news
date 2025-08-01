import { Menu } from 'lucide-react';

export default function Header() {
  return (
    <header className="bg-blue-600 text-white p-4 shadow-md">
      <div className="container mx-auto flex justify-between items-center">
        <h1 className="text-2xl font-bold">Game News</h1>
        <button className="p-2 rounded hover:bg-blue-700">
          <Menu size={24} />
        </button>
      </div>
    </header>
  );
}