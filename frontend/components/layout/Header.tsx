'use client';

import { Menu } from 'lucide-react';
import Link from 'next/link';
import Search from '@/components/ui/Search';
import { useState } from 'react';

export default function Header() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const handleSearch = (query: string) => {
    console.log('Searching for:', query);
    // TODO: Implement search functionality
  };

  return (
    <header className="bg-blue-600 text-white shadow-md">
      <div className="container mx-auto px-4 py-4">
        <div className="flex justify-between items-center">
          <Link href="/" className="text-2xl font-bold hover:text-blue-200 transition-colors">
            Game News
          </Link>
          
          <div className="hidden md:flex items-center space-x-4 flex-1 max-w-md mx-8">
            <Search 
              onSearch={handleSearch}
              placeholder="Search gaming news..."
              className="w-full"
            />
          </div>
          
          <nav className="hidden md:flex items-center space-x-6">
            <Link href="/" className="hover:text-blue-200 transition-colors">
              Home
            </Link>
            <Link href="/about" className="hover:text-blue-200 transition-colors">
              About
            </Link>
            <Link href="/contact" className="hover:text-blue-200 transition-colors">
              Contact
            </Link>
          </nav>
          
          <button 
            className="md:hidden p-2 rounded hover:bg-blue-700 transition-colors"
            onClick={() => setIsMenuOpen(!isMenuOpen)}
          >
            <Menu size={24} />
          </button>
        </div>
        
        {/* Mobile menu */}
        {isMenuOpen && (
          <div className="md:hidden mt-4 pb-4 border-t border-blue-500">
            <div className="mt-4 mb-4">
              <Search 
                onSearch={handleSearch}
                placeholder="Search gaming news..."
              />
            </div>
            <nav className="flex flex-col space-y-2">
              <Link href="/" className="py-2 hover:text-blue-200 transition-colors">
                Home
              </Link>
              <Link href="/about" className="py-2 hover:text-blue-200 transition-colors">
                About
              </Link>
              <Link href="/contact" className="py-2 hover:text-blue-200 transition-colors">
                Contact
              </Link>
            </nav>
          </div>
        )}
      </div>
    </header>
  );
}
