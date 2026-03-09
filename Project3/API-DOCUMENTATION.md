# Private Knowledge Base - API Documentation

## 📡 **API Overview**

The Private Knowledge Base provides a RESTful API for document management, AI-powered chat, and user authentication. All API endpoints are protected with JWT authentication except for health checks and login.

**Base URL**: `http://localhost:8080/api/v1`  
**Authentication**: Bearer Token (JWT)  
**Content-Type**: `application/json` (except file uploads)

---

## 🔐 **Authentication**

### **Login**
Authenticate user and receive JWT token for subsequent requests.

```http
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "username": "testuser",
  "password": "testpass"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "user123",
  "username": "testuser",
  "expires_in": 3600
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "Invalid credentials"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Username and password are required"
}
```

---

### **Refresh Token**
Refresh an existing JWT token.

```http
POST /api/v1/auth/refresh
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "token": "new.jwt.token.here",
  "expires_in": 3600
}
```

---

### **Logout**
Invalidate user session and clear token.

```http
POST /api/v1/auth/logout
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "message": "Logged out successfully"
}
```

---

## 📄 **Document Management**

### **Upload Document**
Upload a document for processing and indexing.

```http
POST /api/v1/documents/upload
```

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data
```

**Request Body (multipart/form-data):**
```
file: <binary_file_data>
```

**Response (200 OK):**
```json
{
  "document_id": "550e8400-e29b-41d4-a716-446655440000",
  "filename": "document.pdf",
  "size": 1024000,
  "content_type": "application/pdf",
  "status": "uploaded",
  "processing_status": "pending"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "No file uploaded"
}
```

**Response (413 Payload Too Large):**
```json
{
  "error": "File size exceeds maximum limit (50MB)"
}
```

---

### **List Documents**
Get paginated list of uploaded documents.

```http
GET /api/v1/documents
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)
- `search` (optional): Search term for filename
- `processed` (optional): Filter by processing status (true/false)

**Response (200 OK):**
```json
{
  "documents": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "filename": "document.pdf",
      "content_type": "application/pdf",
      "file_size": 1024000,
      "upload_time": "2026-03-09T12:00:00Z",
      "processed": true,
      "processing_status": "completed",
      "chunk_count": 25,
      "metadata": {
        "pages": 10,
        "title": "Sample Document"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### **Get Document Details**
Get detailed information about a specific document.

```http
GET /api/v1/documents/{document_id}
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `document_id`: UUID of the document

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "filename": "document.pdf",
  "content_type": "application/pdf",
  "file_size": 1024000,
  "upload_time": "2026-03-09T12:00:00Z",
  "processed": true,
  "processing_status": "completed",
  "chunk_count": 25,
  "metadata": {
    "pages": 10,
    "title": "Sample Document",
    "author": "John Doe",
    "created_date": "2026-03-01T10:00:00Z"
  },
  "chunks": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "chunk_index": 0,
      "content": "This is the first chunk of the document...",
      "page_number": 1,
      "metadata": {
        "type": "paragraph",
        "section": "introduction"
      }
    }
  ]
}
```

**Response (404 Not Found):**
```json
{
  "error": "Document not found"
}
```

---

### **Delete Document**
Delete a document and all associated chunks.

```http
DELETE /api/v1/documents/{document_id}
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `document_id`: UUID of the document

**Response (200 OK):**
```json
{
  "message": "Document deleted successfully"
}
```

**Response (404 Not Found):**
```json
{
  "error": "Document not found"
}
```

---

### **Reindex Document**
Reprocess and re-index a document.

```http
POST /api/v1/documents/{document_id}/reindex
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `document_id`: UUID of the document

**Response (200 OK):**
```json
{
  "message": "Document reindexing started",
  "document_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "processing"
}
```

---

### **Get Document Chunks**
Get all chunks for a specific document.

```http
GET /api/v1/documents/{document_id}/chunks
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `document_id`: UUID of the document

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 50)

**Response (200 OK):**
```json
{
  "chunks": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "document_id": "550e8400-e29b-41d4-a716-446655440000",
      "chunk_index": 0,
      "content": "This is the first chunk of the document...",
      "page_number": 1,
      "metadata": {
        "type": "paragraph",
        "section": "introduction"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 25,
    "total_pages": 1
  }
}
```

---

## 💬 **Chat & AI**

### **Create Chat Session**
Create a new chat session.

```http
POST /api/v1/chat/sessions
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "session_name": "Document Analysis",
  "model": "llama3.1:8b"
}
```

**Response (200 OK):**
```json
{
  "session_id": "770e8400-e29b-41d4-a716-446655440000",
  "session_name": "Document Analysis",
  "model": "llama3.1:8b",
  "created_at": "2026-03-09T12:00:00Z",
  "last_activity": "2026-03-09T12:00:00Z"
}
```

---

### **Get Chat Sessions**
Get all chat sessions for the current user.

```http
GET /api/v1/chat/sessions
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "sessions": [
    {
      "session_id": "770e8400-e29b-41d4-a716-446655440000",
      "session_name": "Document Analysis",
      "model": "llama3.1:8b",
      "created_at": "2026-03-09T12:00:00Z",
      "last_activity": "2026-03-09T12:30:00Z",
      "message_count": 15
    }
  ]
}
```

---

### **Send Chat Message**
Send a message and get AI response (Server-Sent Events streaming).

```http
POST /api/v1/chat
```

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "session_id": "770e8400-e29b-41d4-a716-446655440000",
  "message": "What are the main points in the uploaded document?",
  "model": "llama3.1:8b",
  "stream": true
}
```

**Response (Server-Sent Events):**
```
data: {"type": "start", "session_id": "770e8400-e29b-41d4-a716-446655440000"}

data: {"type": "chunk", "content": "Based on the document I analyzed"}

data: {"type": "chunk", "content": ", the main points are:"}

data: {"type": "chunk", "content": "1. Introduction to the topic"}

data: {"type": "chunk", "content": "2. Key findings and analysis"}

data: {"type": "chunk", "content": "3. Conclusion and recommendations"}

data: {"type": "end", "message_id": "880e8400-e29b-41d4-a716-446655440000", "citations": ["doc1_chunk5", "doc2_chunk12"], "tokens_used": 150, "response_time_ms": 3200}
```

**Non-Streaming Response:**
```json
{
  "session_id": "770e8400-e29b-41d4-a716-446655440000",
  "message_id": "880e8400-e29b-41d4-a716-446655440000",
  "content": "Based on the document I analyzed, the main points are: 1. Introduction to the topic 2. Key findings and analysis 3. Conclusion and recommendations",
  "citations": [
    {
      "document_id": "550e8400-e29b-41d4-a716-446655440000",
      "chunk_id": "660e8400-e29b-41d4-a716-446655440001",
      "content": "Introduction to the topic...",
      "page_number": 1
    }
  ],
  "model_used": "llama3.1:8b",
  "tokens_used": 150,
  "response_time_ms": 3200,
  "created_at": "2026-03-09T12:00:00Z"
}
```

---

### **Get Chat History**
Get message history for a specific session.

```http
GET /api/v1/chat/sessions/{session_id}/messages
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `session_id`: UUID of the chat session

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 50)

**Response (200 OK):**
```json
{
  "messages": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440000",
      "session_id": "770e8400-e29b-41d4-a716-446655440000",
      "message_type": "user",
      "content": "What are the main points in the uploaded document?",
      "created_at": "2026-03-09T12:00:00Z"
    },
    {
      "id": "990e8400-e29b-41d4-a716-446655440000",
      "session_id": "770e8400-e29b-41d4-a716-446655440000",
      "message_type": "assistant",
      "content": "Based on the document I analyzed, the main points are...",
      "citations": [
        {
          "document_id": "550e8400-e29b-41d4-a716-446655440000",
          "chunk_id": "660e8400-e29b-41d4-a716-446655440001",
          "content": "Introduction to the topic...",
          "page_number": 1
        }
      ],
      "model_used": "llama3.1:8b",
      "tokens_used": 150,
      "response_time_ms": 3200,
      "created_at": "2026-03-09T12:00:05Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 2,
    "total_pages": 1
  }
}
```

---

### **Delete Chat Session**
Delete a chat session and all messages.

```http
DELETE /api/v1/chat/sessions/{session_id}
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `session_id`: UUID of the chat session

**Response (200 OK):**
```json
{
  "message": "Chat session deleted successfully"
}
```

---

## 🔍 **Search**

### **Semantic Search**
Perform semantic search across all documents.

```http
POST /api/v1/search
```

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "query": "machine learning algorithms",
  "limit": 10,
  "threshold": 0.7,
  "document_ids": ["550e8400-e29b-41d4-a716-446655440000"]
}
```

**Response (200 OK):**
```json
{
  "results": [
    {
      "chunk": {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "document_id": "550e8400-e29b-41d4-a716-446655440000",
        "chunk_index": 5,
        "content": "Machine learning algorithms are computational methods that enable systems to learn patterns from data...",
        "page_number": 3,
        "metadata": {
          "type": "paragraph",
          "section": "algorithms"
        }
      },
      "similarity": 0.89,
      "document": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "filename": "ml-guide.pdf",
        "title": "Machine Learning Guide"
      }
    }
  ],
  "total": 1,
  "query": "machine learning algorithms",
  "processing_time_ms": 45
}
```

---

### **Hybrid Search**
Perform hybrid search combining semantic and keyword search.

```http
POST /api/v1/search/hybrid
```

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "query": "deep learning neural networks",
  "semantic_weight": 0.7,
  "keyword_weight": 0.3,
  "limit": 10
}
```

**Response (200 OK):**
```json
{
  "results": [
    {
      "chunk": {
        "id": "660e8400-e29b-41d4-a716-446655440002",
        "document_id": "550e8400-e29b-41d4-a716-446655440000",
        "chunk_index": 12,
        "content": "Deep learning neural networks are a subset of machine learning...",
        "page_number": 7,
        "metadata": {
          "type": "paragraph",
          "section": "deep-learning"
        }
      },
      "similarity": 0.92,
      "keyword_score": 0.85,
      "combined_score": 0.89,
      "document": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "filename": "ml-guide.pdf",
        "title": "Machine Learning Guide"
      }
    }
  ],
  "total": 1,
  "query": "deep learning neural networks",
  "processing_time_ms": 62
}
```

---

## 📊 **System & Health**

### **Health Check**
Check system health status.

```http
GET /health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2026-03-09T12:00:00Z",
  "version": "1.0.0",
  "uptime": "2h30m15s",
  "checks": {
    "database": "ok",
    "redis": "ok",
    "ollama": "ok"
  }
}
```

**Response (503 Service Unavailable):**
```json
{
  "status": "unhealthy",
  "timestamp": "2026-03-09T12:00:00Z",
  "checks": {
    "database": "error",
    "redis": "ok",
    "ollama": "ok"
  }
}
```

---

### **Readiness Check**
Check if system is ready to accept requests.

```http
GET /ready
```

**Response (200 OK):**
```json
{
  "status": "ready",
  "checks": {
    "database": "ok",
    "redis": "ok",
    "ollama": "ok",
    "disk_space": "ok"
  }
}
```

---

### **Metrics**
Get Prometheus metrics.

```http
GET /metrics
```

**Response (200 OK):**
```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",status="200"} 1250
http_requests_total{method="POST",status="200"} 340

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.1"} 1200
http_request_duration_seconds_bucket{le="0.5"} 1450
http_request_duration_seconds_bucket{le="1.0"} 1580

# HELP documents_processed_total Total number of documents processed
# TYPE documents_processed_total counter
documents_processed_total{status="completed"} 45
documents_processed_total{status="failed"} 2
```

---

## 🤖 **AI Models**

### **List Available Models**
Get list of available AI models.

```http
GET /api/v1/models
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "models": [
    {
      "name": "llama3.1:8b",
      "display_name": "Llama 3.1 8B",
      "type": "chat",
      "size": "4.9GB",
      "description": "General purpose chat model",
      "status": "available"
    },
    {
      "name": "nomic-embed-text",
      "display_name": "Nomic Embed Text",
      "type": "embedding",
      "size": "274MB",
      "description": "Text embedding model",
      "status": "available"
    }
  ]
}
```

---

### **Generate Embeddings**
Generate embeddings for text.

```http
POST /api/v1/models/embeddings
```

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "model": "nomic-embed-text",
  "text": "This is a sample text for embedding generation"
}
```

**Response (200 OK):**
```json
{
  "model": "nomic-embed-text",
  "embeddings": [0.1234, -0.5678, 0.9012, ...],
  "dimensions": 768,
  "processing_time_ms": 120
}
```

---

## 📈 **User & Analytics**

### **Get User Profile**
Get current user profile.

```http
GET /api/v1/user/profile
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "user_id": "user123",
  "username": "testuser",
  "created_at": "2026-03-01T10:00:00Z",
  "last_login": "2026-03-09T12:00:00Z",
  "statistics": {
    "documents_uploaded": 15,
    "chat_sessions": 8,
    "total_messages": 125,
    "tokens_used": 15420
  }
}
```

---

### **Get Usage Statistics**
Get usage statistics for the current user.

```http
GET /api/v1/user/statistics
```

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `period` (optional): Time period (day/week/month/year, default: month)

**Response (200 OK):**
```json
{
  "period": "month",
  "statistics": {
    "documents_uploaded": 5,
    "chat_sessions": 3,
    "total_messages": 45,
    "tokens_used": 2340,
    "average_response_time_ms": 2100,
    "most_used_model": "llama3.1:8b"
  },
  "daily_breakdown": [
    {
      "date": "2026-03-09",
      "messages": 15,
      "tokens_used": 780,
      "documents_uploaded": 2
    }
  ]
}
```

---

## ⚠️ **Error Handling**

### **Standard Error Response Format**

```json
{
  "error": "Error message description",
  "code": "ERROR_CODE",
  "timestamp": "2026-03-09T12:00:00Z",
  "request_id": "req_123456789"
}
```

### **Common HTTP Status Codes**

| Status | Code | Description |
|--------|------|-------------|
| 200 | OK | Request successful |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid request parameters |
| 401 | Unauthorized | Authentication required or invalid |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource conflict |
| 413 | Payload Too Large | File size exceeds limit |
| 422 | Unprocessable Entity | Validation failed |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |
| 503 | Service Unavailable | Service temporarily unavailable |

### **Rate Limiting**

- **Default Limit**: 100 requests per minute per user
- **File Upload Limit**: 10 uploads per minute per user
- **Chat Messages**: 50 messages per minute per user

**Rate Limit Response (429 Too Many Requests):**
```json
{
  "error": "Rate limit exceeded",
  "code": "RATE_LIMIT_EXCEEDED",
  "retry_after": 60
}
```

---

## 🔧 **Configuration**

### **Request Headers**

| Header | Required | Description |
|--------|----------|-------------|
| `Authorization` | Yes (except public endpoints) | JWT token in format `Bearer <token>` |
| `Content-Type` | Yes (POST/PUT requests) | `application/json` or `multipart/form-data` |
| `Accept` | No | Response content type preference |
| `X-Request-ID` | No | Unique request identifier for tracing |

### **Response Headers**

| Header | Description |
|--------|-------------|
| `Content-Type` | Response content type |
| `X-Request-ID` | Request identifier (if provided) |
| `X-Rate-Limit-Remaining` | Remaining requests in current window |
| `X-Rate-Limit-Reset` | Time when rate limit resets |

---

## 📝 **API Versioning**

- **Current Version**: v1
- **Version in URL**: `/api/v1/`
- **Backward Compatibility**: Maintained within major versions
- **Depreciation**: 6 months notice for breaking changes

---

## 🧪 **Testing Examples**

### **cURL Examples**

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "testpass"}'
```

**Upload Document:**
```bash
curl -X POST http://localhost:8080/api/v1/documents/upload \
  -H "Authorization: Bearer <jwt_token>" \
  -F "file=@document.pdf"
```

**Send Chat Message:**
```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "770e8400-e29b-41d4-a716-446655440000",
    "message": "What are the main points?",
    "model": "llama3.1:8b"
  }'
```

**Search Documents:**
```bash
curl -X POST http://localhost:8080/api/v1/search \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "machine learning",
    "limit": 5
  }'
```

---

This API documentation provides comprehensive information for integrating with the Private Knowledge Base system. All endpoints are designed with RESTful principles, proper error handling, and comprehensive testing capabilities.
