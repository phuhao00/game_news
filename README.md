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
- **PostgreSQL** or **SQLite** for data storage
- **CORS** support for frontend integration
- **BCrypt** for password hashing

## Features

- Real-time game news aggregation from multiple sources (GameSpot, IGN)
- Modern and responsive UI
- Fast and efficient backend
- Type-safe codebase with TypeScript
- Single Page Application (SPA) architecture
- Automated news scraping with Colly
- Persistent data storage with PostgreSQL or SQLite
- Full article content retrieval
- Original source linking
- Search functionality
- Source filtering
- User authentication (registration/login)
- Bookmarking system

## Project Structure

```
game_news/
├── frontend/           # React frontend application
│   ├── src/            # Source code
│   │   ├── components/ # Reusable UI components
│   │   ├── pages/      # Page components
│   │   ├── hooks/      # Custom React hooks
│   │   ├── services/   # API service layer
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
- Docker and Docker Compose (for containerized deployment)

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

3. Run the production server:
   ```bash
   cd ..
   go run main.go
   ```

The application will be available at http://localhost:8080

## Docker & Docker Compose Deployment

This application supports deployment using Docker and Docker Compose for easy setup and scalability.

### Using Docker Compose (Recommended)

1. Build and start all services:
   ```bash
   docker-compose up -d
   ```

This will start the following services:
- Backend application on port 8080
- PostgreSQL database on port 5432
- Adminer database management tool on port 8081

2. Access the application:
   - Main application: http://localhost:8080
   - Database management: http://localhost:8081

### Using Individual Docker Commands

1. Build the Docker image:
   ```bash
   docker build -t game-news .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 game-news
   ```

### Development with Hot Reloading

For development with hot reloading, use the development Docker setup:

1. Start the development environment:
   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d
   ```

This uses [Air](https://github.com/cosmtrek/air) for Go hot reloading.

## Troubleshooting

### Network Issues with Docker Compose

If you encounter network issues when pulling images with Docker Compose:

1. Try using a different Docker registry mirror
2. Manually pull the images:
   ```bash
   docker pull postgres:15-alpine
   docker pull adminer:4.8.1
   ```
3. Then run:
   ```bash
   docker-compose up -d
   ```

### Database Connection Issues

If the application cannot connect to the database:

1. Ensure all services are running:
   ```bash
   docker-compose ps
   ```

2. Check the logs:
   ```bash
   docker-compose logs db
   docker-compose logs backend
   ```

3. Verify the database connection parameters in the docker-compose.yml file

### Frontend Build Issues

If you encounter issues building the frontend:

1. Clean the node_modules:
   ```bash
   cd frontend
   rm -rf node_modules
   npm install
   ```

2. Try building again:
   ```bash
   npm run build
   ```

## API Endpoints

### Public Endpoints
- `GET /api/news` - Get all news (with optional `source` query parameter)
- `GET /api/news/:id` - Get a specific news by ID with full content
- `GET /api/search` - Search news by query string (`q` parameter)
- `GET /api/sources` - Get all news sources
- `POST /api/users/register` - Register a new user
- `POST /api/users/login` - Login as a user

### Protected Endpoints
- `POST /api/protected/bookmarks` - Add a bookmark
- `DELETE /api/protected/bookmarks` - Remove a bookmark
- `GET /api/protected/bookmarks` - Get user's bookmarks

## Web Scraping

The application uses the Colly web scraping framework to collect game news from various sources:
- GameSpot (https://www.gamespot.com/news/)
- IGN (https://www.ign.com/news)

The scraper runs periodically to fetch the latest news and update the storage. It respects website rate limits to avoid being blocked.

## Data Storage

Articles are stored persistently in either PostgreSQL or SQLite with the following features:
- Automatic cleanup of articles older than 7 days
- Efficient lookup by ID
- Sorting by publication date
- Content caching
- User management with password hashing
- Bookmark system

The application automatically detects if a PostgreSQL database is available and connects to it. If not, it falls back to SQLite for simpler setups.

## Mobile App Version

For information about creating a mobile app version of this application, please see [MobileREADME.md](MobileREADME.md).

## Deployment

To deploy this application to a production environment:

1. Build the frontend:
   ```bash
   cd frontend
   npm run build
   ```

2. The Go backend will automatically serve the built frontend from the `dist/` directory.

3. Run the Go application:
   ```bash
   cd ..
   go run main.go
   ```

4. Set up a reverse proxy (like Nginx) in front of the application for production use.

## Future Improvements

- [ ] Add RSS feed generation
- [ ] Implement push notifications for breaking news
- [ ] Add pagination for large result sets
- [ ] Improve search algorithm with full-text search
- [ ] Add social sharing features
- [ ] Implement caching for better performance
- [ ] Add admin panel for content management