# Book Library API (Gorilla Mux)

A simple REST API for managing a book collection using the Gorilla Mux router.

## Features

- **Create** books
- **Read** all books or single book
- **Update** book details
- **Delete** books
- **Logging middleware** for request tracking

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/books` | Get all books |
| GET | `/books/{id}` | Get a specific book |
| POST | `/books` | Create a new book |
| PUT | `/books/{id}` | Update a book |
| DELETE | `/books/{id}` | Delete a book |

## Data Model

```json
{
  "id": "1",
  "title": "The Go Programming Language",
  "author": "Alan A. A. Donovan"
}
```

## Getting Started

1. Navigate to the book-library directory:
```bash
cd book-library
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8081`

## API Usage Examples

### Create a book
```bash
curl -X POST http://localhost:8081/books \
  -H "Content-Type: application/json" \
  -d '{"title": "Clean Architecture", "author": "Robert C. Martin"}'
```

### Get all books
```bash
curl http://localhost:8081/books
```

### Get a specific book
```bash
curl http://localhost:8081/books/1
```

### Update a book
```bash
curl -X PUT http://localhost:8081/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "The Go Programming Language (Updated)", "author": "Alan A. A. Donovan"}'
```

### Delete a book
```bash
curl -X DELETE http://localhost:8081/books/1
```

## Sample Data

The API starts with two sample books:
- Book 1: "The Go Programming Language" by Alan A. A. Donovan
- Book 2: "Clean Code" by Robert C. Martin

## Logging

The API includes a logging middleware that logs all HTTP requests with method, path, and client IP address.
