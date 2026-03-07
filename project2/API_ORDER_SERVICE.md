# Order Service API

**Port:** 8082  
**Base URL:** `http://localhost:8082`

---

## Endpoints

### 1. Create Order

**POST** `/orders`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "product_name": "Gaming Laptop",
  "price": 1299
}
```

**Response (201 Created):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "product_name": "Gaming Laptop",
  "price": 1299,
  "status": "pending",
  "created_at": "2026-03-07T08:46:15.234567890Z"
}
```

**cURL:**
```bash
curl -X POST http://localhost:8082/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id": "550e8400-e29b-41d4-a716-446655440000", "product_name": "Gaming Laptop", "price": 1299}'
```

**PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8082/orders" -Method POST -ContentType "application/json" -Body '{"user_id": "550e8400-e29b-41d4-a716-446655440000", "product_name": "Gaming Laptop", "price": 1299}'
```

---

### 2. Get Order by ID

**GET** `/orders/{id}`

**Response (200 OK):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "product_name": "Gaming Laptop",
  "price": 1299,
  "status": "pending",
  "created_at": "2026-03-07T08:46:15.234567890Z"
}
```

**cURL:**
```bash
curl http://localhost:8082/orders/660e8400-e29b-41d4-a716-446655440001
```

---

### 3. Health Check

**GET** `/health`

**Response (200 OK):**
```json
{
  "service": "order-service",
  "status": "healthy",
  "timestamp": "2026-03-07T08:46:15Z",
  "database": "connected",
  "cached_users": 3
}
```

**cURL:**
```bash
curl http://localhost:8082/health
```

---

### 4. Prometheus Metrics

**GET** `/metrics`

**Response (200 OK):**
```
# HELP order_service_http_requests_total Total HTTP requests
# TYPE order_service_http_requests_total counter
order_service_http_requests_total{method="POST",status="201"} 1

# HELP order_service_db_queries_total Total database queries
# TYPE order_service_db_queries_total counter
order_service_db_queries_total 5

# HELP order_service_kafka_messages_produced_total Kafka messages produced
# TYPE order_service_kafka_messages_produced_total counter
order_service_kafka_messages_produced_total{topic="order.created"} 1

# HELP order_service_kafka_messages_consumed_total Kafka messages consumed
# TYPE order_service_kafka_messages_consumed_total counter
order_service_kafka_messages_consumed_total{topic="user.created"} 3

# HELP order_service_users_cached Users in cache
# TYPE order_service_users_cached gauge
order_service_users_cached 3
```

**cURL:**
```bash
curl http://localhost:8082/metrics
```

---

### 5. Internal Statistics

**GET** `/internal/stats`

**Response (200 OK):**
```json
{
  "service": "order-service",
  "timestamp": "2026-03-07T08:46:15Z",
  "uptime": "5m30s",
  "totalOrders": 2,
  "cachedUsers": 3,
  "database": {
    "host": "postgres:5432",
    "dbname": "orderdb",
    "connected": true
  },
  "kafka": {
    "brokers": "kafka:29092",
    "producerReady": true,
    "consumerReady": true,
    "topics": ["user.created", "order.created"]
  }
}
```

**cURL:**
```bash
curl http://localhost:8082/internal/stats
```

---

### 6. User Cache (Kafka Event Viewer)

**GET** `/internal/cache`

**Response (200 OK):**
```json
{
  "550e8400-e29b-41d4-a716-446655440000": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "created_at": "2026-03-07T08:45:30.123456789Z"
  }
}
```

**cURL:**
```bash
curl http://localhost:8082/internal/cache
```

**Note:** This endpoint shows users received via Kafka `user.created` events. It verifies the event-driven architecture is working.
