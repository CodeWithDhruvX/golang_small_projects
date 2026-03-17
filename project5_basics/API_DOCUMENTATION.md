# API Endpoint Documentation

This document provides detailed API endpoint documentation for both projects in the repository.

---

## 🚀 Task Manager API (Gin Framework)
**Base URL**: `http://localhost:8080`

### Data Model
```json
{
  "id": 1,
  "title": "Learn Go",
  "completed": false
}
```

### Endpoints

#### 1. Get All Tasks
- **Method**: `GET`
- **Endpoint**: `/tasks`
- **Description**: Retrieves all tasks in the system
- **Response**: `200 OK`
```json
[
  {
    "id": 1,
    "title": "Learn Go",
    "completed": false
  },
  {
    "id": 2,
    "title": "Build REST API",
    "completed": true
  }
]
```

#### 2. Get Single Task
- **Method**: `GET`
- **Endpoint**: `/tasks/:id`
- **Description**: Retrieves a specific task by ID
- **Parameters**:
  - `id` (path): Task ID (integer)
- **Success Response**: `200 OK`
```json
{
  "id": 1,
  "title": "Learn Go",
  "completed": false
}
```
- **Error Responses**:
  - `400 Bad Request`: Invalid task ID
  - `404 Not Found`: Task not found

#### 3. Create Task
- **Method**: `POST`
- **Endpoint**: `/tasks`
- **Description**: Creates a new task
- **Request Body**:
```json
{
  "title": "New Task",
  "completed": false
}
```
- **Success Response**: `201 Created`
```json
{
  "id": 3,
  "title": "New Task",
  "completed": false
}
```
- **Error Response**: `400 Bad Request`: Invalid JSON format

#### 4. Update Task
- **Method**: `PUT`
- **Endpoint**: `/tasks/:id`
- **Description**: Updates an existing task
- **Parameters**:
  - `id` (path): Task ID (integer)
- **Request Body**:
```json
{
  "title": "Updated Task",
  "completed": true
}
```
- **Success Response**: `200 OK`
```json
{
  "id": 1,
  "title": "Updated Task",
  "completed": true
}
```
- **Error Responses**:
  - `400 Bad Request`: Invalid task ID or JSON
  - `404 Not Found`: Task not found

#### 5. Delete Task
- **Method**: `DELETE`
- **Endpoint**: `/tasks/:id`
- **Description**: Deletes a task by ID
- **Parameters**:
  - `id` (path): Task ID (integer)
- **Success Response**: `200 OK`
```json
{
  "message": "Task deleted successfully"
}
```
- **Error Responses**:
  - `400 Bad Request`: Invalid task ID
  - `404 Not Found`: Task not found

---

## 🌐 Book Library API (Gorilla Mux)
**Base URL**: `http://localhost:8081`

### Data Model
```json
{
  "id": "1",
  "title": "The Go Programming Language",
  "author": "Alan A. A. Donovan"
}
```

### Endpoints

#### 1. Get All Books
- **Method**: `GET`
- **Endpoint**: `/books`
- **Description**: Retrieves all books in the library
- **Response**: `200 OK`
```json
[
  {
    "id": "1",
    "title": "The Go Programming Language",
    "author": "Alan A. A. Donovan"
  },
  {
    "id": "2",
    "title": "Clean Code",
    "author": "Robert C. Martin"
  }
]
```

#### 2. Get Single Book
- **Method**: `GET`
- **Endpoint**: `/books/{id}`
- **Description**: Retrieves a specific book by ID
- **Parameters**:
  - `id` (path): Book ID (string)
- **Success Response**: `200 OK`
```json
{
  "id": "1",
  "title": "The Go Programming Language",
  "author": "Alan A. A. Donovan"
}
```
- **Error Response**: `404 Not Found`: Book not found

#### 3. Create Book
- **Method**: `POST`
- **Endpoint**: `/books`
- **Description**: Adds a new book to the library
- **Request Body**:
```json
{
  "title": "Clean Architecture",
  "author": "Robert C. Martin"
}
```
- **Success Response**: `201 Created`
```json
{
  "id": "3",
  "title": "Clean Architecture",
  "author": "Robert C. Martin"
}
```
- **Error Response**: `400 Bad Request`: Invalid JSON format

#### 4. Update Book
- **Method**: `PUT`
- **Endpoint**: `/books/{id}`
- **Description**: Updates an existing book's information
- **Parameters**:
  - `id` (path): Book ID (string)
- **Request Body**:
```json
{
  "title": "Updated Book Title",
  "author": "Updated Author Name"
}
```
- **Success Response**: `200 OK`
```json
{
  "id": "1",
  "title": "Updated Book Title",
  "author": "Updated Author Name"
}
```
- **Error Responses**:
  - `400 Bad Request`: Invalid JSON format
  - `404 Not Found`: Book not found

#### 5. Delete Book
- **Method**: `DELETE`
- **Endpoint**: `/books/{id}`
- **Description**: Removes a book from the library
- **Parameters**:
  - `id` (path): Book ID (string)
- **Success Response**: `200 OK`
```json
{
  "message": "Book deleted successfully"
}
```
- **Error Response**: `404 Not Found`: Book not found

---

## 🧪 Testing Examples

### Task Manager API (Gin)

```bash
# Create a task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Gin Framework", "completed": false}'

# Get all tasks
curl http://localhost:8080/tasks

# Get specific task
curl http://localhost:8080/tasks/1

# Update task
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Gin Framework", "completed": true}'

# Delete task
curl -X DELETE http://localhost:8080/tasks/1
```

### Book Library API (Gorilla Mux)

```bash
# Create a book
curl -X POST http://localhost:8081/books \
  -H "Content-Type: application/json" \
  -d '{"title": "Design Patterns", "author": "Gang of Four"}'

# Get all books
curl http://localhost:8081/books

# Get specific book
curl http://localhost:8081/books/1

# Update book
curl -X PUT http://localhost:8081/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Design Patterns", "author": "Erich Gamma et al."}'

# Delete book
curl -X DELETE http://localhost:8081/books/1
```

---

## 📊 Response Format

### Success Responses
- **200 OK**: Request successful
- **201 Created**: Resource created successfully

### Error Responses
- **400 Bad Request**: Invalid request data or parameters
- **404 Not Found**: Resource not found

All error responses include a descriptive error message in JSON format.

---

## 🔧 Technical Details

### Thread Safety
Both APIs use `sync.RWMutex` for thread-safe access to in-memory data stores.

### Data Persistence
Both APIs use in-memory storage (maps) for simplicity. Data is reset when the server restarts.

### Logging
The Book Library API includes request logging middleware that logs:
- HTTP method
- Request path
- Client IP address

### Content-Type
All endpoints expect and return JSON data with `Content-Type: application/json`.
