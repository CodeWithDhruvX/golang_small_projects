# API Testing Guide

## Quick Start - Test the Complete Flow

All services are now running. Here are the commands to test the event-driven flow:

### 1. Create a User

**PowerShell:**
```powershell
$userResponse = Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body '{"name": "John Doe", "email": "john@example.com"}'
Write-Host "Created user with ID: $($userResponse.id)"
$userId = $userResponse.id
```

**curl:**
```bash
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

**Expected Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2026-03-07T08:45:30.123456789Z"
}
```

---

### 2. Create an Order

**PowerShell:**
```powershell
$orderBody = @{
    user_id = $userId
    product_name = "Gaming Laptop"
    price = 1299
} | ConvertTo-Json

$orderResponse = Invoke-RestMethod -Uri "http://localhost:8082/orders" -Method POST -ContentType "application/json" -Body $orderBody
Write-Host "Created order with ID: $($orderResponse.id)"
$orderId = $orderResponse.id
```

**curl:**
```bash
curl -X POST http://localhost:8082/orders \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": \"$USER_ID\", \"product_name\": \"Gaming Laptop\", \"price\": 1299}"
```

**Expected Response:**
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

---

### 3. Verify Payment Was Created Automatically

The Payment Service automatically consumes the `order.created` event and creates a payment.

**PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8083/payments" | Format-List
```

**curl:**
```bash
curl http://localhost:8083/payments
```

**Expected Response:**
```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "order_id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "amount": 1299,
    "payment_status": "success",
    "created_at": "2026-03-07T08:46:15.345678901Z"
  }
]
```

---

## Complete Test Script

### PowerShell Full Test Script

Save this as `test-flow.ps1`:

```powershell
# Test the complete event-driven flow

Write-Host "=== Testing Event-Driven Microservices ===" -ForegroundColor Green
Write-Host ""

# 1. Create User
Write-Host "1. Creating User..." -ForegroundColor Yellow
$userResponse = Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body '{"name": "Alice Smith", "email": "alice@example.com"}'
Write-Host "   User ID: $($userResponse.id)"
Write-Host "   Name: $($userResponse.name)"
$userId = $userResponse.id
Write-Host ""

# 2. Get User (verify storage)
Write-Host "2. Verifying User in Database..." -ForegroundColor Yellow
$getUser = Invoke-RestMethod -Uri "http://localhost:8081/users/$userId"
Write-Host "   Retrieved: $($getUser.name) ($($getUser.email))"
Write-Host ""

# 3. Create Order
Write-Host "3. Creating Order..." -ForegroundColor Yellow
$orderBody = @{
    user_id = $userId
    product_name = "Wireless Headphones"
    price = 199
} | ConvertTo-Json

$orderResponse = Invoke-RestMethod -Uri "http://localhost:8082/orders" -Method POST -ContentType "application/json" -Body $orderBody
Write-Host "   Order ID: $($orderResponse.id)"
Write-Host "   Product: $($orderResponse.product_name)"
Write-Host "   Price: $($orderResponse.price)"
$orderId = $orderResponse.id
Write-Host ""

# 4. Wait for payment processing
Write-Host "4. Waiting for Payment Processing (2 seconds)..." -ForegroundColor Yellow
Start-Sleep -Seconds 2

# 5. Check Payment
Write-Host "5. Checking Automatic Payment Creation..." -ForegroundColor Yellow
$payments = Invoke-RestMethod -Uri "http://localhost:8083/payments"
Write-Host "   Total Payments: $($payments.Count)"
if ($payments.Count -gt 0) {
    $latestPayment = $payments | Select-Object -Last 1
    Write-Host "   Latest Payment ID: $($latestPayment.id)"
    Write-Host "   Order ID: $($latestPayment.order_id)"
    Write-Host "   Amount: $($latestPayment.amount)"
    Write-Host "   Status: $($latestPayment.payment_status)"
}
Write-Host ""

# 6. Check Order Service User Cache (Kafka event consumption)
Write-Host "6. Checking Order Service User Cache (from Kafka)..." -ForegroundColor Yellow
$cache = Invoke-RestMethod -Uri "http://localhost:8082/internal/cache"
Write-Host "   Cached Users: $($cache.Count)"
Write-Host ""

# 7. Health Checks
Write-Host "7. Health Checks..." -ForegroundColor Yellow
$userHealth = Invoke-RestMethod -Uri "http://localhost:8081/health"
$orderHealth = Invoke-RestMethod -Uri "http://localhost:8082/health"
$paymentHealth = Invoke-RestMethod -Uri "http://localhost:8083/health"

Write-Host "   User Service: $($userHealth.status) (DB: $($userHealth.database))"
Write-Host "   Order Service: $($orderHealth.status) (DB: $($orderHealth.database))"
Write-Host "   Payment Service: $($paymentHealth.status) (DB: $($paymentHealth.database))"
Write-Host ""

# 8. Internal Statistics
Write-Host "8. Internal Statistics..." -ForegroundColor Yellow
$userStats = Invoke-RestMethod -Uri "http://localhost:8081/internal/stats"
$orderStats = Invoke-RestMethod -Uri "http://localhost:8082/internal/stats"
$paymentStats = Invoke-RestMethod -Uri "http://localhost:8083/internal/stats"

Write-Host "   User Service - Total Users: $($userStats.totalUsers)"
Write-Host "   Order Service - Total Orders: $($orderStats.totalOrders), Cached Users: $($orderStats.cachedUsers)"
Write-Host "   Payment Service - Total Payments: $($paymentStats.totalPayments)"
Write-Host ""

Write-Host "=== Test Complete ===" -ForegroundColor Green
```

**Run it:**
```powershell
.\test-flow.ps1
```

---

## API Reference

### User Service (Port 8081)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users` | Create a new user |
| GET | `/users/{id}` | Get user by ID |
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |
| GET | `/internal/stats` | Internal statistics |

**Create User Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

---

### Order Service (Port 8082)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/orders` | Create a new order |
| GET | `/orders/{id}` | Get order by ID |
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |
| GET | `/internal/stats` | Internal statistics |
| GET | `/internal/cache` | View user cache from Kafka |

**Create Order Request:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "product_name": "Gaming Laptop",
  "price": 1299
}
```

---

### Payment Service (Port 8083)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/payments` | Get all payments |
| GET | `/payments/{id}` | Get payment by ID |
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |
| GET | `/internal/stats` | Internal statistics |

---

## Testing with curl (Linux/Mac/Git Bash)

### Full Test Script

```bash
#!/bin/bash

echo "=== Testing Event-Driven Microservices ==="
echo ""

# 1. Create User
echo "1. Creating User..."
USER_RESPONSE=$(curl -s -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}')
echo "   Response: $USER_RESPONSE"
USER_ID=$(echo $USER_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "   User ID: $USER_ID"
echo ""

# 2. Create Order
echo "2. Creating Order..."
ORDER_RESPONSE=$(curl -s -X POST http://localhost:8082/orders \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": \"$USER_ID\", \"product_name\": \"Gaming Laptop\", \"price\": 1299}")
echo "   Response: $ORDER_RESPONSE"
ORDER_ID=$(echo $ORDER_RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "   Order ID: $ORDER_ID"
echo ""

# 3. Wait for payment processing
echo "3. Waiting for payment processing (2s)..."
sleep 2

# 4. Check Payments
echo "4. Checking Payments..."
curl -s http://localhost:8083/payments | jq .
echo ""

# 5. Health Checks
echo "5. Health Checks..."
echo "   User Service: $(curl -s http://localhost:8081/health | jq -r '.status')"
echo "   Order Service: $(curl -s http://localhost:8082/health | jq -r '.status')"
echo "   Payment Service: $(curl -s http://localhost:8083/health | jq -r '.status')"
echo ""

echo "=== Test Complete ==="
```

---

## Testing with Postman

### Collection Setup

1. **Create a Collection** called "Event-Driven Microservices"

2. **Add Requests:**

**Create User (POST):**
- URL: `http://localhost:8081/users`
- Headers: `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Get User (GET):**
- URL: `http://localhost:8081/users/{{userId}}`

**Create Order (POST):**
- URL: `http://localhost:8082/orders`
- Headers: `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "user_id": "{{userId}}",
  "product_name": "Gaming Laptop",
  "price": 1299
}
```

**Get Payments (GET):**
- URL: `http://localhost:8083/payments`

3. **Set Environment Variables:**
```json
{
  "userId": "",
  "orderId": ""
}
```

4. **Add Tests to Extract IDs:**
In the "Create User" request Tests tab:
```javascript
var jsonData = pm.response.json();
pm.environment.set("userId", jsonData.id);
```

In the "Create Order" request Tests tab:
```javascript
var jsonData = pm.response.json();
pm.environment.set("orderId", jsonData.id);
```

---

## Monitoring Endpoints

### Prometheus Metrics

**View all metrics:**
```powershell
# User Service
Invoke-RestMethod -Uri "http://localhost:8081/metrics"

# Order Service
Invoke-RestMethod -Uri "http://localhost:8082/metrics"

# Payment Service
Invoke-RestMethod -Uri "http://localhost:8083/metrics"
```

### Key Metrics to Watch

| Metric | Description |
|--------|-------------|
| `*_http_requests_total` | Total HTTP requests by method/endpoint/status |
| `*_http_request_duration_seconds` | Request duration histogram |
| `*_db_queries_total` | Database queries by operation |
| `*_kafka_messages_produced_total` | Kafka messages produced by topic |
| `*_kafka_messages_consumed_total` | Kafka messages consumed by topic |
| `*_active_connections` | Active HTTP connections |

### Internal Statistics

**User Service Stats:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8081/internal/stats"
```
Output:
```json
{
  "service": "user-service",
  "timestamp": "2026-03-07T08:50:00Z",
  "totalUsers": 5
}
```

**Order Service Stats:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8082/internal/stats"
```
Output:
```json
{
  "service": "order-service",
  "timestamp": "2026-03-07T08:50:00Z",
  "totalOrders": 3,
  "cachedUsers": 5
}
```

**Payment Service Stats:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8083/internal/stats"
```
Output:
```json
{
  "service": "payment-service",
  "timestamp": "2026-03-07T08:50:00Z",
  "totalPayments": 3,
  "paymentsByStatus": {
    "success": 2,
    "failed": 1,
    "pending": 0
  }
}
```

---

## Load Testing

### Using PowerShell (Simple)

```powershell
# Create 10 users and orders
1..10 | ForEach-Object {
    $user = Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body "{`"name`": `"User$_`", `"email`": `"user$_@example.com`"}"
    $order = Invoke-RestMethod -Uri "http://localhost:8082/orders" -Method POST -ContentType "application/json" -Body "{`"user_id`": `"$($user.id)`", `"product_name`": `"Product$_`", `"price`": $($_ * 100)}"
    Write-Host "Created user $($user.id) and order $($order.id)"
}
```

### Using curl (Parallel)

```bash
# Create 5 users in parallel
for i in {1..5}; do
  curl -s -X POST http://localhost:8081/users \
    -H "Content-Type: application/json" \
    -d "{\"name\": \"User$i\", \"email\": \"user$i@example.com\"}" &
done
wait
```

---

## Troubleshooting

### Check if services are running:
```powershell
docker-compose ps
```

### View logs:
```powershell
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f user-service
docker-compose logs -f order-service
docker-compose logs -f payment-service
```

### Test Kafka topics:
```powershell
# List topics
docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 --list

# View messages in topic
docker exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic user.created --from-beginning
docker exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic order.created --from-beginning
docker exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic payment.completed --from-beginning
```

### Check databases:
```powershell
# PostgreSQL - Users
docker-compose exec postgres psql -U postgres -d userdb -c "SELECT * FROM users;"

# PostgreSQL - Orders
docker-compose exec postgres psql -U postgres -d orderdb -c "SELECT * FROM orders;"

# MongoDB
docker-compose exec mongodb mongosh paymentdb --eval "db.payments.find()"
```

---

## Expected Event Flow

```
1. POST /users (User Service)
   ↓
   Stores in PostgreSQL
   ↓
   Publishes user.created → Kafka
   ↓
2. Order Service consumes user.created
   ↓
   Caches user in memory
   ↓
3. POST /orders (Order Service)
   ↓
   Stores in PostgreSQL
   ↓
   Publishes order.created → Kafka
   ↓
4. Payment Service consumes order.created
   ↓
   Processes payment
   ↓
   Stores in MongoDB
   ↓
   Publishes payment.completed → Kafka
```

---

## Summary

| Action | Command |
|--------|---------|
| Create User | `POST http://localhost:8081/users` |
| Create Order | `POST http://localhost:8082/orders` |
| View Payments | `GET http://localhost:8083/payments` |
| Health Check | `GET http://localhost:8081/health` |
| Metrics | `GET http://localhost:8081/metrics` |
| Stats | `GET http://localhost:8081/internal/stats` |

All monitoring endpoints are available on all three services (ports 8081, 8082, 8083).
