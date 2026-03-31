# Project 7 - GORM Multi-Database Integration

A Go application that demonstrates integration between GORM (PostgreSQL) and MongoDB with Docker setup and database UIs.

## Features

- **PostgreSQL Database**: Uses GORM ORM for structured data (Users, Posts, Categories, Products)
- **MongoDB Database**: NoSQL database for flexible data (Logs, User Profiles, Analytics, Sessions)
- **Docker Integration**: Both databases run in Docker containers
- **Database UIs**: Adminer for PostgreSQL, Mongo Express for MongoDB
- **RESTful API**: Complete CRUD operations for both databases
- **Graceful Shutdown**: Proper cleanup and shutdown handling

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.25+ 
- Git

### Setup

1. **Clone and navigate to the project**:
   ```bash
   cd project7
   ```

2. **Start the databases with Docker**:
   ```bash
   docker-compose up -d
   ```

3. **Install Go dependencies**:
   ```bash
   go mod download
   ```

4. **Run the application**:
   ```bash
   go run cmd/main.go
   ```

The server will start on `http://localhost:8082`

## Database UIs

After starting Docker containers, you can access the database UIs:

- **Adminer (PostgreSQL)**: http://localhost:8080
  - Server: postgres
  - Username: admin
  - Password: password123
  - Database: project7_db

- **Mongo Express (MongoDB)**: http://localhost:8081
  - Username: admin
  - Password: admin123

## API Endpoints

### Health Check
- `GET /health` - Application health status

### PostgreSQL Endpoints (Users & Posts)

#### Users
- `POST /api/v1/users` - Create a user
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get a specific user
- `PUT /api/v1/users/:id` - Update a user
- `DELETE /api/v1/users/:id` - Delete a user

#### Posts
- `POST /api/v1/posts` - Create a post
- `GET /api/v1/posts` - Get all posts
- `GET /api/v1/posts/:id` - Get a specific post

### MongoDB Endpoints (Profiles, Logs, Analytics)

#### User Profiles
- `POST /api/v1/profiles` - Create a user profile
- `GET /api/v1/profiles/:user_id` - Get user profile

#### Logs
- `POST /api/v1/logs` - Create a log entry
- `GET /api/v1/logs` - Get log entries

#### Analytics
- `POST /api/v1/analytics/events` - Create an analytics event
- `GET /api/v1/analytics/events` - Get analytics events

## Example API Usage

### Create a User (PostgreSQL)
```bash
curl -X POST http://localhost:8082/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "age": 30
  }'
```

### Create a User Profile (MongoDB)
```bash
curl -X POST http://localhost:8082/api/v1/profiles \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "bio": "Software Developer",
    "location": "New York",
    "website": "https://johndoe.dev",
    "social_links": {
      "github": "johndoe",
      "twitter": "@johndoe"
    }
  }'
```

### Create a Log Entry (MongoDB)
```bash
curl -X POST http://localhost:8082/api/v1/logs \
  -H "Content-Type: application/json" \
  -d '{
    "level": "info",
    "message": "User logged in",
    "service": "auth-service",
    "user_id": 1,
    "request_id": "req-123456"
  }'
```

## Project Structure

```
project7/
├── cmd/
│   └── main.go              # Main application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── database/
│   │   ├── postgres.go      # PostgreSQL/GORM operations
│   │   └── mongodb.go       # MongoDB operations
│   ├── handlers/
│   │   └── handlers.go      # HTTP request handlers
│   └── models/
│       ├── postgres.go      # PostgreSQL/GORM models
│       └── mongodb.go       # MongoDB models
├── docker-compose.yml       # Docker configuration
├── .env.example            # Environment variables template
├── go.mod                  # Go module file
└── README.md               # This file
```

## Configuration

Copy `.env.example` to `.env` and modify as needed:

```bash
cp .env.example .env
```

## Database Schemas

### PostgreSQL (GORM)

#### Users Table
- id (PK)
- username (unique)
- email (unique)
- first_name, last_name
- age
- active (boolean)
- timestamps

#### Posts Table
- id (PK)
- title, content
- published (boolean)
- user_id (FK)
- timestamps
- Many-to-many relationship with Tags

#### Categories & Products
- Standard e-commerce structure with relationships

### MongoDB Collections

#### log_entries
- Application logs with timestamps and metadata
- Indexed by timestamp and level

#### user_profiles
- Extended user information with flexible schema
- Indexed by user_id

#### analytics_events
- Event tracking with properties and metadata
- Indexed by timestamp and event_type

#### sessions
- User session management
- Indexed by session_id

## Development

### Adding New Models

1. **PostgreSQL**: Add to `internal/models/postgres.go`
2. **MongoDB**: Add to `internal/models/mongodb.go`
3. **Update migrations** in the database connection files
4. **Add handlers** in `internal/handlers/handlers.go`

### Environment Variables

The application uses the following environment variables:
- `SERVER_PORT`: HTTP server port (default: 8082)
- `POSTGRES_*`: PostgreSQL connection settings
- `MONGO_*`: MongoDB connection settings

## Cleanup

To stop and remove all containers:
```bash
docker-compose down -v
```

## License

This project is for educational purposes to demonstrate multi-database integration in Go.
