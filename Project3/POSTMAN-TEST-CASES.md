# Private Knowledge Base - Postman Test Cases Guide

## 🧪 **Complete API Testing Collection**

This guide provides comprehensive test cases for the Private Knowledge Base API using the Postman collection. The collection covers all endpoints with proper authentication, data validation, and error handling tests.

---

## 📋 **Setup Instructions**

### **1. Import Collection**
1. Open Postman
2. Click "Import" → "Select Files"
3. Choose `postman-collection.json`
4. Select the collection and click "Import"

### **2. Configure Environment Variables**
The collection includes the following variables:

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `baseUrl` | `http://localhost:8080` | API base URL |
| `jwt_token` | *auto-populated* | Authentication token |
| `user_id` | *auto-populated* | Current user ID |
| `username` | *auto-populated* | Current username |
| `document_id` | *auto-populated* | Uploaded document ID |
| `session_id` | *auto-populated* | Chat session ID |
| `chat_model` | `llama3.1:8b` | Default chat model |

### **3. Test Data Preparation**
Create a test document (PDF, TXT, or MD) for upload tests:
- Sample PDF: `test-document.pdf`
- Sample text: `test-document.txt`
- Sample markdown: `test-document.md`

---

## 🔐 **Authentication Test Cases**

### **Test Case 1: Successful Login**
**Request**: `POST /api/v1/auth/login`
```json
{
  "username": "testuser",
  "password": "testpass"
}
```

**Expected Response** (200 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "user123",
  "username": "testuser",
  "expires_in": 3600
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains token
- ✅ Response contains user_id
- ✅ Response contains username
- ✅ Token is not empty
- ✅ Token stored in collection variables

---

### **Test Case 2: Invalid Login**
**Request**: `POST /api/v1/auth/login`
```json
{
  "username": "wronguser",
  "password": "wrongpass"
}
```

**Expected Response** (401 Unauthorized):
```json
{
  "error": "Invalid credentials"
}
```

**Test Assertions**:
- ✅ Status code is 401
- ✅ Response contains error message

---

### **Test Case 3: Missing Credentials**
**Request**: `POST /api/v1/auth/login`
```json
{
  "username": "",
  "password": ""
}
```

**Expected Response** (400 Bad Request):
```json
{
  "error": "Username and password are required"
}
```

**Test Assertions**:
- ✅ Status code is 400
- ✅ Response contains validation error

---

### **Test Case 4: Token Refresh**
**Request**: `POST /api/v1/auth/refresh`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
```json
{
  "token": "new.jwt.token.here",
  "expires_in": 3600
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains new token
- ✅ Token stored in collection variables

---

### **Test Case 5: Logout**
**Request**: `POST /api/v1/auth/logout`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
```json
{
  "message": "Logged out successfully"
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Token cleared from collection variables

---

## 📄 **Document Management Test Cases**

### **Test Case 6: Upload Document**
**Request**: `POST /api/v1/documents/upload`
**Headers**: `Authorization: Bearer <jwt_token>`
**Body**: multipart/form-data with file

**Expected Response** (200 OK):
```json
{
  "document_id": "550e8400-e29b-41d4-a716-446655440000",
  "filename": "test-document.pdf",
  "size": 1024000,
  "content_type": "application/pdf",
  "status": "uploaded",
  "processing_status": "pending"
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains document_id
- ✅ Response contains filename
- ✅ Response contains file size
- ✅ Status is "uploaded"
- ✅ Document ID stored in collection variables

---

### **Test Case 7: Upload Without File**
**Request**: `POST /api/v1/documents/upload`
**Headers**: `Authorization: Bearer <jwt_token>`
**Body**: Empty multipart form

**Expected Response** (400 Bad Request):
```json
{
  "error": "No file uploaded"
}
```

**Test Assertions**:
- ✅ Status code is 400
- ✅ Response contains error message

---

### **Test Case 8: Upload Oversized File**
**Request**: `POST /api/v1/documents/upload`
**Headers**: `Authorization: Bearer <jwt_token>`
**Body**: File larger than 50MB

**Expected Response** (413 Payload Too Large):
```json
{
  "error": "File size exceeds maximum limit (50MB)"
}
```

**Test Assertions**:
- ✅ Status code is 413
- ✅ Response contains size limit error

---

### **Test Case 9: List Documents**
**Request**: `GET /api/v1/documents?page=1&limit=20`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
```json
{
  "documents": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "filename": "test-document.pdf",
      "content_type": "application/pdf",
      "file_size": 1024000,
      "upload_time": "2026-03-09T12:00:00Z",
      "processed": true,
      "processing_status": "completed",
      "chunk_count": 25
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

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains documents array
- ✅ Response contains pagination object
- ✅ Pagination structure is correct

---

### **Test Case 10: Get Document Details**
**Request**: `GET /api/v1/documents/{{document_id}}`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "filename": "test-document.pdf",
  "content_type": "application/pdf",
  "file_size": 1024000,
  "upload_time": "2026-03-09T12:00:00Z",
  "processed": true,
  "processing_status": "completed",
  "chunk_count": 25,
  "chunks": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "chunk_index": 0,
      "content": "This is the first chunk...",
      "page_number": 1
    }
  ]
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains document details
- ✅ Response contains chunks array
- ✅ Chunk structure is correct

---

### **Test Case 11: Get Non-Existent Document**
**Request**: `GET /api/v1/documents/invalid-uuid`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (404 Not Found):
```json
{
  "error": "Document not found"
}
```

**Test Assertions**:
- ✅ Status code is 404
- ✅ Response contains not found error

---

### **Test Case 12: Get Document Chunks**
**Request**: `GET /api/v1/documents/{{document_id}}/chunks?page=1&limit=50`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
```json
{
  "chunks": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "document_id": "550e8400-e29b-41d4-a716-446655440000",
      "chunk_index": 0,
      "content": "This is the first chunk...",
      "page_number": 1,
      "metadata": {
        "type": "paragraph"
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

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains chunks array
- ✅ Response contains pagination
- ✅ Chunk structure is correct

---

### **Test Case 13: Reindex Document**
**Request**: `POST /api/v1/documents/{{document_id}}/reindex`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
```json
{
  "message": "Document reindexing started",
  "document_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "processing"
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains reindexing message
- ✅ Status is "processing"

---

### **Test Case 14: Delete Document**
**Request**: `DELETE /api/v1/documents/{{document_id}}`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
```json
{
  "message": "Document deleted successfully"
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains success message
- ✅ Document ID cleared from collection variables

---

## 💬 **Chat & AI Test Cases**

### **Test Case 15: Create Chat Session**
**Request**: `POST /api/v1/chat/sessions`
**Headers**: `Authorization: Bearer <jwt_token>`
```json
{
  "session_name": "Test Session",
  "model": "llama3.1:8b"
}
```

**Expected Response** (200 OK):
```json
{
  "session_id": "770e8400-e29b-41d4-a716-446655440000",
  "session_name": "Test Session",
  "model": "llama3.1:8b",
  "created_at": "2026-03-09T12:00:00Z",
  "last_activity": "2026-03-09T12:00:00Z"
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains session_id
- ✅ Response contains session details
- ✅ Session ID stored in collection variables

---

### **Test Case 16: Get Chat Sessions**
**Request**: `GET /api/v1/chat/sessions`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
```json
{
  "sessions": [
    {
      "session_id": "770e8400-e29b-41d4-a716-446655440000",
      "session_name": "Test Session",
      "model": "llama3.1:8b",
      "created_at": "2026-03-09T12:00:00Z",
      "last_activity": "2026-03-09T12:30:00Z",
      "message_count": 15
    }
  ]
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains sessions array
- ✅ Session structure is correct

---

### **Test Case 17: Send Chat Message**
**Request**: `POST /api/v1/chat`
**Headers**: `Authorization: Bearer <jwt_token>`
```json
{
  "session_id": "{{session_id}}",
  "message": "What are the main points in the uploaded document?",
  "model": "{{chat_model}}",
  "stream": false
}
```

**Expected Response** (200 OK):
```json
{
  "session_id": "770e8400-e29b-41d4-a716-446655440000",
  "message_id": "880e8400-e29b-41d4-a716-446655440000",
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
  "created_at": "2026-03-09T12:00:00Z"
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains chat response
- ✅ Response contains citations
- ✅ Response contains metadata
- ✅ Content is not empty
- ✅ Citations is an array

---

### **Test Case 18: Send Message to Invalid Session**
**Request**: `POST /api/v1/chat`
**Headers**: `Authorization: Bearer <jwt_token>`
```json
{
  "session_id": "invalid-session-id",
  "message": "Test message",
  "model": "llama3.1:8b"
}
```

**Expected Response** (404 Not Found):
```json
{
  "error": "Chat session not found"
}
```

**Test Assertions**:
- ✅ Status code is 404
- ✅ Response contains session not found error

---

### **Test Case 19: Get Chat History**
**Request**: `GET /api/v1/chat/sessions/{{session_id}}/messages?page=1&limit=50`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
```json
{
  "messages": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440000",
      "session_id": "770e8400-e29b-41d4-a716-446655440000",
      "message_type": "user",
      "content": "What are the main points?",
      "created_at": "2026-03-09T12:00:00Z"
    },
    {
      "id": "990e8400-e29b-41d4-a716-446655440000",
      "session_id": "770e8400-e29b-41d4-a716-446655440000",
      "message_type": "assistant",
      "content": "Based on the document...",
      "citations": [...],
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

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains messages array
- ✅ Message structure is correct
- ✅ Message types are valid ('user' or 'assistant')

---

### **Test Case 20: Delete Chat Session**
**Request**: `DELETE /api/v1/chat/sessions/{{session_id}}`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
```json
{
  "message": "Chat session deleted successfully"
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains success message
- ✅ Session variables cleared

---

## 🔍 **Search Test Cases**

### **Test Case 21: Semantic Search**
**Request**: `POST /api/v1/search`
**Headers**: `Authorization: Bearer <jwt_token>`
```json
{
  "query": "machine learning algorithms",
  "limit": 10,
  "threshold": 0.7
}
```

**Expected Response** (200 OK):
```json
{
  "results": [
    {
      "chunk": {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "document_id": "550e8400-e29b-41d4-a716-446655440000",
        "chunk_index": 5,
        "content": "Machine learning algorithms are computational methods...",
        "page_number": 3
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

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains results array
- ✅ Results contain similarity scores
- ✅ Similarity scores are between 0 and 1
- ✅ Processing time is recorded

---

### **Test Case 22: Empty Search Query**
**Request**: `POST /api/v1/search`
**Headers**: `Authorization: Bearer <jwt_token>`
```json
{
  "query": "",
  "limit": 10
}
```

**Expected Response** (400 Bad Request):
```json
{
  "error": "Query cannot be empty"
}
```

**Test Assertions**:
- ✅ Status code is 400
- ✅ Response contains validation error

---

### **Test Case 23: Hybrid Search**
**Request**: `POST /api/v1/search/hybrid`
**Headers**: `Authorization: Bearer <jwt_token>`
```json
{
  "query": "deep learning neural networks",
  "semantic_weight": 0.7,
  "keyword_weight": 0.3,
  "limit": 10
}
```

**Expected Response** (200 OK):
```json
{
  "results": [
    {
      "chunk": {
        "id": "660e8400-e29b-41d4-a716-446655440002",
        "content": "Deep learning neural networks are a subset...",
        "page_number": 7
      },
      "similarity": 0.92,
      "keyword_score": 0.85,
      "combined_score": 0.89,
      "document": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "filename": "ml-guide.pdf"
      }
    }
  ],
  "total": 1,
  "query": "deep learning neural networks",
  "processing_time_ms": 62
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Results contain similarity scores
- ✅ Results contain keyword scores
- ✅ Results contain combined scores
- ✅ All scores are valid numbers

---

## 🤖 **AI Models Test Cases**

### **Test Case 24: List Available Models**
**Request**: `GET /api/v1/models`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
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

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains models array
- ✅ Model structure is correct
- ✅ Model types are valid ('chat' or 'embedding')

---

### **Test Case 25: Generate Embeddings**
**Request**: `POST /api/v1/models/embeddings`
**Headers**: `Authorization: Bearer <jwt_token>`
```json
{
  "model": "nomic-embed-text",
  "text": "This is a sample text for embedding generation"
}
```

**Expected Response** (200 OK):
```json
{
  "model": "nomic-embed-text",
  "embeddings": [0.1234, -0.5678, 0.9012, ...],
  "dimensions": 768,
  "processing_time_ms": 120
}
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains embeddings array
- ✅ Embeddings array length matches dimensions
- ✅ All embeddings are numbers
- ✅ Processing time is recorded

---

### **Test Case 26: Invalid Model for Embeddings**
**Request**: `POST /api/v1/models/embeddings`
**Headers**: `Authorization: Bearer <jwt_token>`
```json
{
  "model": "llama3.1:8b",
  "text": "Test text"
}
```

**Expected Response** (400 Bad Request):
```json
{
  "error": "Model 'llama3.1:8b' is not an embedding model"
}
```

**Test Assertions**:
- ✅ Status code is 400
- ✅ Response contains model type error

---

## 📊 **User & Analytics Test Cases**

### **Test Case 27: Get User Profile**
**Request**: `GET /api/v1/user/profile`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
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

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains user profile
- ✅ Response contains statistics
- ✅ Statistics values are numbers

---

### **Test Case 28: Get Usage Statistics**
**Request**: `GET /api/v1/user/statistics?period=month`
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (200 OK):
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

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains statistics
- ✅ Response contains daily breakdown
- ✅ All values are valid numbers

---

## 🏥 **System & Health Test Cases**

### **Test Case 29: Health Check**
**Request**: `GET /health`

**Expected Response** (200 OK):
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

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains health status
- ✅ All system checks are "ok"
- ✅ Status is "healthy"

---

### **Test Case 30: Readiness Check**
**Request**: `GET /ready`

**Expected Response** (200 OK):
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

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains readiness status
- ✅ All checks are "ok"
- ✅ Status is "ready"

---

### **Test Case 31: Metrics Endpoint**
**Request**: `GET /metrics`

**Expected Response** (200 OK):
```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",status="200"} 1250
http_requests_total{method="POST",status="200"} 340

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.1"} 1200
http_request_duration_seconds_bucket{le="0.5"} 1450
```

**Test Assertions**:
- ✅ Status code is 200
- ✅ Response contains Prometheus metrics
- ✅ Metrics format is correct
- ✅ Key metrics are present

---

## 🚨 **Error Handling Test Cases**

### **Test Case 32: Unauthorized Access**
**Request**: Any protected endpoint without token
**Headers**: No Authorization header

**Expected Response** (401 Unauthorized):
```json
{
  "error": "Authorization token required"
}
```

**Test Assertions**:
- ✅ Status code is 401
- ✅ Response contains auth error

---

### **Test Case 33: Invalid Token**
**Request**: Any protected endpoint with invalid token
**Headers**: `Authorization: Bearer invalid.token.here`

**Expected Response** (401 Unauthorized):
```json
{
  "error": "Invalid or expired token"
}
```

**Test Assertions**:
- ✅ Status code is 401
- ✅ Response contains token error

---

### **Test Case 34: Rate Limiting**
**Request**: Rapid requests to any endpoint
**Headers**: `Authorization: Bearer <jwt_token>`

**Expected Response** (429 Too Many Requests):
```json
{
  "error": "Rate limit exceeded",
  "code": "RATE_LIMIT_EXCEEDED",
  "retry_after": 60
}
```

**Test Assertions**:
- ✅ Status code is 429
- ✅ Response contains rate limit error
- ✅ Response contains retry_after header

---

### **Test Case 35: Invalid JSON**
**Request**: POST endpoint with malformed JSON
**Headers**: `Authorization: Bearer <jwt_token>`
**Body**: `{"invalid": json}`

**Expected Response** (400 Bad Request):
```json
{
  "error": "Invalid JSON format"
}
```

**Test Assertions**:
- ✅ Status code is 400
- ✅ Response contains JSON error

---

## 🔄 **Workflow Test Scenarios**

### **Scenario 1: Complete Document Workflow**
1. **Login** → Get JWT token
2. **Upload Document** → Upload test file
3. **List Documents** → Verify document appears
4. **Get Document Details** → Verify processing
5. **Search Document** → Find content
6. **Delete Document** → Clean up

### **Scenario 2: Complete Chat Workflow**
1. **Login** → Get JWT token
2. **Create Chat Session** → Start conversation
3. **Send Message** → Get AI response
4. **Get Chat History** → Verify conversation
5. **Delete Chat Session** → Clean up

### **Scenario 3: Document + Chat Integration**
1. **Login** → Get JWT token
2. **Upload Document** → Add knowledge
3. **Create Chat Session** → Start conversation
4. **Ask Question** → Get response with citations
5. **Verify Citations** → Check document references

---

## 🧪 **Performance Test Cases**

### **Test Case 36: Concurrent Requests**
- Send 10 simultaneous requests to `/health`
- Verify all responses are successful
- Check response times are reasonable

### **Test Case 37: Large Document Upload**
- Upload a 10MB document
- Verify upload completes successfully
- Check processing completes within reasonable time

### **Test Case 38: Complex Search Query**
- Search with long query and many results
- Verify search completes within 5 seconds
- Check result relevance and accuracy

---

## 📊 **Test Execution Guide**

### **Running All Tests**
1. Open the imported collection in Postman
2. Click on the collection name
3. Click "Run" button
4. Select "Run entire collection"
5. Configure iterations (recommended: 1)
6. Click "Run Private Knowledge Base API"

### **Running Specific Test Suites**
- **Authentication**: Run only "🔐 Authentication" folder
- **Document Management**: Run only "📄 Document Management" folder
- **Chat & AI**: Run only "💬 Chat & AI" folder
- **Search**: Run only "🔍 Search" folder

### **Test Results Analysis**
- **Pass Rate**: Should be 100% for all tests
- **Response Times**: Should be under 5 seconds for most endpoints
- **Error Handling**: Verify proper error responses for edge cases
- **Data Validation**: Check response structures match expected formats

---

## 🔧 **Troubleshooting**

### **Common Issues**
1. **401 Unauthorized**: Check if JWT token is valid and not expired
2. **Connection Refused**: Verify backend is running on correct port
3. **File Upload Fails**: Check file size and format constraints
4. **Search No Results**: Ensure documents are processed and indexed

### **Debug Tips**
- Use Postman Console for detailed request/response logs
- Check collection variables after each test run
- Verify environment variables are correctly set
- Use Postman's "Test Results" tab for assertion failures

---

This comprehensive test suite provides complete coverage of the Private Knowledge Base API, ensuring all functionality works correctly and handles edge cases appropriately.
