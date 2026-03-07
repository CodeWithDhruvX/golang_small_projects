# Payment Service API

**Port:** 8083  
**Base URL:** `http://localhost:8083`

---

## Endpoints

### 1. Get All Payments

**GET** `/payments`

**Response (200 OK):**
```json
[
  {
    "_id": "770e8400-e29b-41d4-a716-446655440002",
    "order_id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "amount": 1299,
    "payment_status": "success",
    "created_at": "2026-03-07T08:46:15.345678901Z"
  }
]
```

**cURL:**
```bash
curl http://localhost:8083/payments
```

**PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8083/payments"
```

---

### 2. Get Payment by ID

**GET** `/payments/{id}`

**Response (200 OK):**
```json
{
  "_id": "770e8400-e29b-41d4-a716-446655440002",
  "order_id": "660e8400-e29b-41d4-a716-446655440001",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 1299,
  "payment_status": "success",
  "created_at": "2026-03-07T08:46:15.345678901Z"
}
```

**cURL:**
```bash
curl http://localhost:8083/payments/770e8400-e29b-41d4-a716-446655440002
```

---

### 3. Health Check

**GET** `/health`

**Response (200 OK):**
```json
{
  "service": "payment-service",
  "status": "healthy",
  "timestamp": "2026-03-07T08:46:15Z",
  "database": "connected"
}
```

**cURL:**
```bash
curl http://localhost:8083/health
```

---

### 4. Prometheus Metrics

**GET** `/metrics`

**Response (200 OK):**
```
# HELP payment_service_http_requests_total Total HTTP requests
# TYPE payment_service_http_requests_total counter
payment_service_http_requests_total{method="GET",status="200"} 5

# HELP payment_service_db_queries_total Total database queries
# TYPE payment_service_db_queries_total counter
payment_service_db_queries_total 10

# HELP payment_service_kafka_messages_produced_total Kafka messages produced
# TYPE payment_service_kafka_messages_produced_total counter
payment_service_kafka_messages_produced_total{topic="payment.completed"} 2

# HELP payment_service_kafka_messages_consumed_total Kafka messages consumed
# TYPE payment_service_kafka_messages_consumed_total counter
payment_service_kafka_messages_consumed_total{topic="order.created"} 2

# HELP payment_service_payments_processed_total Payments processed
# TYPE payment_service_payments_processed_total counter
payment_service_payments_processed_total{status="success"} 1
payment_service_payments_processed_total{status="failed"} 1
```

**cURL:**
```bash
curl http://localhost:8083/metrics
```

---

### 5. Internal Statistics

**GET** `/internal/stats`

**Response (200 OK):**
```json
{
  "service": "payment-service",
  "timestamp": "2026-03-07T08:46:15Z",
  "uptime": "5m30s",
  "totalPayments": 2,
  "paymentsByStatus": {
    "success": 1,
    "failed": 1
  },
  "database": {
    "host": "mongo:27017",
    "dbname": "paymentdb",
    "connected": true
  },
  "kafka": {
    "brokers": "kafka:29092",
    "producerReady": true,
    "consumerReady": true,
    "topics": ["order.created", "payment.completed"]
  }
}
```

**cURL:**
```bash
curl http://localhost:8083/internal/stats
```

---

## Payment Status Logic

Payments are created **automatically** when an order is created (via Kafka `order.created` event).

### Status Rules:

| Price Ends With | Status | Reason |
|----------------|--------|---------|
| Any number except 99 | `success` | Normal payment |
| Ends with 99 (e.g., 199, 299, 1299) | `failed` | Simulates payment failure |

### Example:

```json
// Price: 1299 → Status: failed (ends with 99)
{
  "_id": "770e8400-e29b-41d4-a716-446655440002",
  "order_id": "660e8400-e29b-41d4-a716-446655440001",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 1299,
  "payment_status": "failed",
  "created_at": "2026-03-07T08:46:15.345678901Z"
}

// Price: 1000 → Status: success (does not end with 99)
{
  "_id": "880e8400-e29b-41d4-a716-446655440003",
  "order_id": "770e8400-e29b-41d4-a716-446655440004",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 1000,
  "payment_status": "success",
  "created_at": "2026-03-07T08:47:20.456789012Z"
}
```

---

## Event Flow

1. **Order Service** creates order
2. **Order Service** publishes `order.created` event to Kafka
3. **Payment Service** consumes `order.created` event
4. **Payment Service** processes payment (simulated)
5. **Payment Service** stores payment in MongoDB
6. **Payment Service** publishes `payment.completed` event to Kafka

**Note:** No manual POST to create payments - they are created automatically via the event-driven architecture!
