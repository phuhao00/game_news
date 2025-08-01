import { generateMetadata } from '@/lib/metadata';

export const metadata = generateMetadata({
  title: 'Contact Us',
  description: 'Get in touch with the Game News team. Send us your questions, feedback, or story tips.',
});

export default function ContactLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}