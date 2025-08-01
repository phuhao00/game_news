import { Metadata } from 'next';

export const siteConfig = {
  name: 'Game News',
  description: 'Latest news and updates from the gaming world',
  url: 'https://gamenews.com',
  ogImage: '/og-image.jpg',
  links: {
    twitter: 'https://twitter.com/gamenews',
    github: 'https://github.com/gamenews',
  },
};

export function generateMetadata({
  title,
  description,
  image,
  noIndex = false,
}: {
  title?: string;
  description?: string;
  image?: string;
  noIndex?: boolean;
} = {}): Metadata {
  return {
    title: title ? `${title} | ${siteConfig.name}` : siteConfig.name,
    description: description || siteConfig.description,
    keywords: [
      'gaming news',
      'video games',
      'game reviews',
      'esports',
      'gaming industry',
      'game releases',
    ],
    authors: [
      {
        name: siteConfig.name,
      },
    ],
    creator: siteConfig.name,
    openGraph: {
      type: 'website',
      locale: 'en_US',
      url: siteConfig.url,
      title: title ? `${title} | ${siteConfig.name}` : siteConfig.name,
      description: description || siteConfig.description,
      siteName: siteConfig.name,
      images: [
        {
          url: image || siteConfig.ogImage,
          width: 1200,
          height: 630,
          alt: siteConfig.name,
        },
      ],
    },
    twitter: {
      card: 'summary_large_image',
      title: title ? `${title} | ${siteConfig.name}` : siteConfig.name,
      description: description || siteConfig.description,
      images: [image || siteConfig.ogImage],
      creator: '@gamenews',
    },
    robots: {
      index: !noIndex,
      follow: !noIndex,
      googleBot: {
        index: !noIndex,
        follow: !noIndex,
        'max-video-preview': -1,
        'max-image-preview': 'large',
        'max-snippet': -1,
      },
    },
  };
}