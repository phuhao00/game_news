# Game News Collector

A modern web application for collecting and browsing the latest game news, built with cutting-edge technologies.

## Technology Stack

### Frontend
- **React 18** with TypeScript
- **React Router** for navigation
- **Vite** for fast development and building
- Modern CSS with responsive design

### Backend
- **Go** programming language
- **Gin Web Framework** for REST API
- **Colly** for web scraping
- **CORS** support for frontend integration

## Features

- Real-time game news aggregation
- Modern and responsive UI
- Fast and efficient backend
- Type-safe codebase with TypeScript
- Single Page Application (SPA) architecture
- Automated news scraping with Colly
- In-memory data storage with cleanup
- Full article content retrieval
- Original source linking

## Project Structure

```
game_news/
├── frontend/           # React frontend application
│   ├── src/            # Source code
│   │   ├── components/ # Reusable UI components
│   │   ├── pages/      # Page components
│   │   ├── hooks/      # Custom React hooks
│   │   ├── types/      # TypeScript types
│   │   └── ...
│   ├── public/         # Static assets
│   └── ...
├── scraper/            # Web scraping functionality
├── storage/            # Data storage management
├── main.go             # Go backend application
├── go.mod              # Go module dependencies
└── README.md           # This file
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Node.js 16 or higher
- npm or yarn

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd game_news
   ```

2. Install frontend dependencies:
   ```bash
   cd frontend
   npm install
   ```

3. Install backend dependencies:
   ```bash
   cd ..
   go mod tidy
   ```

### Development

1. Start the Go backend server:
   ```bash
   go run main.go
   ```

2. Start the React frontend development server:
   ```bash
   cd frontend
   npm run dev
   ```

The application will be available at:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

### Building for Production

1. Build the React frontend:
   ```bash
   cd frontend
   npm run build
   ```

2. The built files will be placed in the `dist/` directory, which is served by the Go backend.

## API Endpoints

- `GET /api/news` - Get all news
- `GET /api/news/:id` - Get a specific news by ID with full content

## Web Scraping

The application uses the Colly web scraping framework to collect game news from various sources. The scraper runs periodically to fetch the latest news and update the storage.

In the current implementation, we're using mock data to demonstrate the functionality. In a production environment, you would replace this with actual scraping logic for real game news websites.

## Data Storage

Articles are stored in-memory with the following features:
- Automatic cleanup of articles older than 7 days
- Efficient lookup by ID
- Sorting by publication date
- Content caching

## Future Improvements

- [ ] Implement real web scraping from actual game news websites
- [ ] Add search and filtering functionality
- [ ] Implement user authentication and bookmarks
- [ ] Add news comments and ratings
- [ ] Create a mobile app version
- [ ] Add database storage instead of in-memory storage
- [ ] Implement rate limiting for API endpoints
- [ ] Add RSS feed generation
- [ ] Implement push notifications for breaking news