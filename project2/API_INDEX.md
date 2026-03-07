# API Documentation Index

Complete API reference for all microservices.

---

## Quick Reference

| Service | Port | Base URL | File |
|---------|------|----------|------|
| **User Service** | 8081 | `http://localhost:8081` | [`API_USER_SERVICE.md`](API_USER_SERVICE.md) |
| **Order Service** | 8082 | `http://localhost:8082` | [`API_ORDER_SERVICE.md`](API_ORDER_SERVICE.md) |
| **Payment Service** | 8083 | `http://localhost:8083` | [`API_PAYMENT_SERVICE.md`](API_PAYMENT_SERVICE.md) |

---

## All Endpoints Summary

### User Service (Port 8081)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users` | Create new user |
| GET | `/users/{id}` | Get user by ID |
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |
| GET | `/internal/stats` | Internal statistics |

### Order Service (Port 8082)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/orders` | Create new order |
| GET | `/orders/{id}` | Get order by ID |
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |
| GET | `/internal/stats` | Internal statistics |
| GET | `/internal/cache` | Kafka user cache |

### Payment Service (Port 8083)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/payments` | Get all payments |
| GET | `/payments/{id}` | Get payment by ID |
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |
| GET | `/internal/stats` | Internal statistics |

---

## Quick Start Commands

### Create a Complete Flow

```powershell
# 1. Create User
$user = Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body '{"name": "John Doe", "email": "john@example.com"}'
Write-Host "User ID: $($user.id)"

# 2. Create Order
$order = Invoke-RestMethod -Uri "http://localhost:8082/orders" -Method POST -ContentType "application/json" -Body "{`"user_id`": `"$($user.id)`", `"product_name`": `"Laptop`", `"price`": 999}"
Write-Host "Order ID: $($order.id)"

# 3. Wait and check payment
Start-Sleep -Seconds 3
$payments = Invoke-RestMethod -Uri "http://localhost:8083/payments"
Write-Host "Payments: $($payments.Count)"
```

---

## Health Check All Services

```powershell
# Check all services
Write-Host "User Service:"
Invoke-RestMethod -Uri "http://localhost:8081/health"

Write-Host "`nOrder Service:"
Invoke-RestMethod -Uri "http://localhost:8082/health"

Write-Host "`nPayment Service:"
Invoke-RestMethod -Uri "http://localhost:8083/health"
```

---

## Event-Driven Architecture

```
┌─────────────────┐     POST /users      ┌─────────────────┐
│   User Service  │ ─────────────────────>│   PostgreSQL    │
│    :8081        │                        │   (userdb)      │
└─────────────────┘                        └─────────────────┘
         │
         │ Kafka: user.created
         ▼
┌─────────────────┐     POST /orders     ┌─────────────────┐
│  Order Service  │ ─────────────────────>│   PostgreSQL    │
│    :8082        │                        │   (orderdb)     │
│                 │<─────────────────────│                 │
│  [User Cache]   │   Consume user events  │                 │
└─────────────────┘                        └─────────────────┘
         │
         │ Kafka: order.created
         ▼
┌─────────────────┐                        ┌─────────────────┐
│ Payment Service │ ─────────────────────> │    MongoDB      │
│    :8083        │    Store payment       │  (paymentdb)    │
│                 │                        └─────────────────┘
│ [Auto-create    │
│  on order event]│
└─────────────────┘
         │
         │ Kafka: payment.completed
         ▼
        (...)
```

---

## Monitoring Endpoints

### Health Checks
```bash
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

### Prometheus Metrics
```bash
curl http://localhost:8081/metrics
curl http://localhost:8082/metrics
curl http://localhost:8083/metrics
```

### Internal Statistics
```bash
curl http://localhost:8081/internal/stats
curl http://localhost:8082/internal/stats
curl http://localhost:8083/internal/stats
```

### Kafka Cache Verification
```bash
# Verify Order Service received user events
curl http://localhost:8082/internal/cache
```

---

## Database Connections

| Service | Database | Connection String |
|---------|----------|-------------------|
| User | PostgreSQL | `postgres://postgres:password@postgres:5432/userdb` |
| Order | PostgreSQL | `postgres://postgres:password@postgres:5432/orderdb` |
| Payment | MongoDB | `mongodb://mongo:27017/paymentdb` |

---

## Kafka Topics

| Topic | Producer | Consumer | Purpose |
|-------|----------|----------|---------|
| `user.created` | User Service | Order Service | Sync user data |
| `order.created` | Order Service | Payment Service | Trigger payment |
| `payment.completed` | Payment Service | (Future use) | Payment confirmation |

---

## Status Codes

| Code | Meaning | When |
|------|---------|------|
| 200 | OK | GET requests successful |
| 201 | Created | POST requests successful |
| 400 | Bad Request | Invalid request body |
| 404 | Not Found | Resource doesn't exist |
| 405 | Method Not Allowed | Wrong HTTP method |
| 500 | Internal Server Error | Database/Kafka error |

---

## Related Documentation

- **[START_GUIDE.md](START_GUIDE.md)** - How to start the project
- **[POSTMAN_TEST_CASES.md](POSTMAN_TEST_CASES.md)** - Postman test documentation
- **[API_TESTING.md](API_TESTING.md)** - Complete testing guide
- **[TEST_CASES.md](TEST_CASES.md)** - Unit test documentation
- **[README.md](README.md)** - Project overview

---

## API Files

- **[API_USER_SERVICE.md](API_USER_SERVICE.md)** - User Service endpoints
- **[API_ORDER_SERVICE.md](API_ORDER_SERVICE.md)** - Order Service endpoints  
- **[API_PAYMENT_SERVICE.md](API_PAYMENT_SERVICE.md)** - Payment Service endpoints
