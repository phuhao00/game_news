import Header from './Header';

interface MainLayoutProps {
  children: React.ReactNode;
}

export default function MainLayout({ children }: MainLayoutProps) {
  return (
    <div className="min-h-screen bg-gray-100">
      <Header />
      <main className="container mx-auto p-4 mt-8">
        {children}
      </main>
      <footer className="bg-blue-600 text-white p-4 mt-12">
        <div className="container mx-auto text-center">
          <p>© {new Date().getFullYear()} Game News. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
}