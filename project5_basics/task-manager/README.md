# Task Manager API (Gin)

A simple REST API for managing tasks using the Gin framework.

## Features

- **Create** tasks
- **Read** all tasks or single task
- **Update** task details
- **Delete** tasks

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks` | Get all tasks |
| GET | `/tasks/:id` | Get a specific task |
| POST | `/tasks` | Create a new task |
| PUT | `/tasks/:id` | Update a task |
| DELETE | `/tasks/:id` | Delete a task |

## Data Model

```json
{
  "id": 1,
  "title": "Learn Go",
  "completed": false
}
```

## Getting Started

1. Navigate to the task-manager directory:
```bash
cd task-manager
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Usage Examples

### Create a task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Gin", "completed": false}'
```

### Get all tasks
```bash
curl http://localhost:8080/tasks
```

### Get a specific task
```bash
curl http://localhost:8080/tasks/1
```

### Update a task
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Gin Framework", "completed": true}'
```

### Delete a task
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

## Sample Data

The API starts with two sample tasks:
- Task 1: "Learn Go" (not completed)
- Task 2: "Build REST API" (completed)
