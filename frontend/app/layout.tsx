import './globals.css';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Game News',
  description: 'Latest news and updates from the gaming world',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}