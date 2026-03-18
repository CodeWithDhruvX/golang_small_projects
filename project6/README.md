# Go CRUD API Projects

This repository contains two complete REST API implementations in Go, demonstrating different routing frameworks and approaches to building CRUD applications.

## Projects

### 1. Task Manager API (Gin Framework)
- **Location**: `task-manager/`
- **Framework**: Gin
- **Port**: 8080
- **Description**: A simple to-do/task management system

### 2. Book Library API (Gorilla Mux)
- **Location**: `book-library/`
- **Framework**: Gorilla Mux
- **Port**: 8081
- **Description**: A book collection management system

## Quick Start

### Task Manager (Gin)
```bash
cd task-manager
go mod tidy
go run main.go
# Server runs on http://localhost:8080
```

### Book Library (Gorilla Mux)
```bash
cd book-library
go mod tidy
go run main.go
# Server runs on http://localhost:8081
```

## Comparison

| Feature | Gin | Gorilla Mux |
|---------|-----|-------------|
| Speed | ⚡ Very fast | Fast |
| Boilerplate | Low | Moderate |
| Routing | Simple | More flexible |
| Learning Curve | Easy | Medium |
| URL Parameters | `:id` | `{id}` |

## API Features

Both APIs implement full CRUD operations:

- **Create**: Add new items
- **Read**: Get all items or single item
- **Update**: Modify existing items
- **Delete**: Remove items

## Testing

You can use `curl` or any API client (Postman, Insomnia) to test the endpoints. See individual README files for detailed API examples.

**📖 Complete API Documentation**: [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) - Comprehensive endpoint documentation with examples

## Project Structure

```
project5_basics/
├── task-manager/          # Gin-based Task Manager API
│   ├── main.go           # Main application file
│   ├── go.mod            # Go module file
│   └── README.md         # API documentation
├── book-library/         # Gorilla Mux-based Book Library API
│   ├── main.go           # Main application file
│   ├── go.mod            # Go module file
│   └── README.md         # API documentation
├── requirements.md       # Original requirements
├── API_DOCUMENTATION.md  # Complete API endpoint documentation
└── README.md            # This file
```

## Technologies Used

- **Go**: Programming language
- **Gin**: HTTP web framework (Task Manager)
- **Gorilla Mux**: HTTP router (Book Library)
- **In-memory storage**: Maps for data persistence
- **Concurrency**: Mutex for thread safety

## Future Enhancements

- Database integration (PostgreSQL, MongoDB)
- Authentication (JWT)
- Pagination and filtering
- Rate limiting
- Docker support
- Unit tests
- Swagger documentation
