# User Service API

**Port:** 8081  
**Base URL:** `http://localhost:8081`

---

## Endpoints

### 1. Create User

**POST** `/users`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "created_at": "2026-03-07T08:45:30.123456789Z"
}
```

**cURL:**
```bash
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john.doe@example.com"}'
```

**PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body '{"name": "John Doe", "email": "john.doe@example.com"}'
```

---

### 2. Get User by ID

**GET** `/users/{id}`

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "created_at": "2026-03-07T08:45:30.123456789Z"
}
```

**cURL:**
```bash
curl http://localhost:8081/users/550e8400-e29b-41d4-a716-446655440000
```

**PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8081/users/550e8400-e29b-41d4-a716-446655440000"
```

---

### 3. Health Check

**GET** `/health`

**Response (200 OK):**
```json
{
  "service": "user-service",
  "status": "healthy",
  "timestamp": "2026-03-07T08:45:30Z",
  "database": "connected"
}
```

**cURL:**
```bash
curl http://localhost:8081/health
```

---

### 4. Prometheus Metrics

**GET** `/metrics`

**Response (200 OK):**
```
# HELP user_service_http_requests_total Total HTTP requests
# TYPE user_service_http_requests_total counter
user_service_http_requests_total{method="POST",status="201"} 1

# HELP user_service_db_queries_total Total database queries
# TYPE user_service_db_queries_total counter
user_service_db_queries_total 5

# HELP user_service_kafka_messages_produced_total Kafka messages produced
# TYPE user_service_kafka_messages_produced_total counter
user_service_kafka_messages_produced_total{topic="user.created"} 1
```

**cURL:**
```bash
curl http://localhost:8081/metrics
```

---

### 5. Internal Statistics

**GET** `/internal/stats`

**Response (200 OK):**
```json
{
  "service": "user-service",
  "timestamp": "2026-03-07T08:45:30Z",
  "uptime": "5m30s",
  "totalUsers": 3,
  "database": {
    "host": "postgres:5432",
    "dbname": "userdb",
    "connected": true
  },
  "kafka": {
    "brokers": "kafka:29092",
    "producerReady": true
  }
}
```

**cURL:**
```bash
curl http://localhost:8081/internal/stats
```
