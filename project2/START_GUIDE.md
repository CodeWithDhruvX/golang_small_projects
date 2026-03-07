# How to Start the Project

## Prerequisites

- **Docker Desktop** installed and running
- **Docker Compose** (included with Docker Desktop)
- **Postman** (optional, for API testing)
- **PowerShell** or **Git Bash** (for Windows)

---

## Quick Start (5 minutes)

### Step 1: Start Docker Desktop

Make sure Docker Desktop is running:
```powershell
# Check Docker status
docker info
```

---

### Step 2: Navigate to Project Directory

```powershell
cd c:\Users\dhruv\Downloads\personal_projects\golang_small_projects\project2
```

---

### Step 3: Build and Start All Services

```powershell
# Build all Docker images (first time only)
docker-compose build

# Start all services in detached mode
docker-compose up -d

# Wait 30 seconds for services to initialize
Start-Sleep -Seconds 30
```

---

### Step 4: Verify Services are Running

```powershell
# Check all containers are up
docker-compose ps
```

You should see:
- ✓ zookeeper - Up
- ✓ kafka - Up
- ✓ postgres - Up
- ✓ mongodb - Up
- ✓ user-service - Up (port 8081)
- ✓ order-service - Up (port 8082)
- ✓ payment-service - Up (port 8083)

---

### Step 5: Test the API

```powershell
# 1. Create a user
$user = Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body '{"name": "John Doe", "email": "john@example.com"}'
Write-Host "User created: $($user.id)"

# 2. Create an order
$order = Invoke-RestMethod -Uri "http://localhost:8082/orders" -Method POST -ContentType "application/json" -Body "{`"user_id`": `"$($user.id)`", `"product_name`": `"Gaming Laptop`", `"price`": 1299}"
Write-Host "Order created: $($order.id)"

# 3. Wait and check payment
Start-Sleep -Seconds 3
$payments = Invoke-RestMethod -Uri "http://localhost:8083/payments"
Write-Host "Payments found: $($payments.Count)"

# 4. View the flow
Write-Host ""
Write-Host "Event Flow Complete!"
Write-Host "User → Order → Payment (auto-created via Kafka)"
```

---

## Verification Commands

### Check Health
```powershell
# All services should return "healthy"
Invoke-RestMethod -Uri "http://localhost:8081/health"
Invoke-RestMethod -Uri "http://localhost:8082/health"
Invoke-RestMethod -Uri "http://localhost:8083/health"
```

### View Metrics
```powershell
# Prometheus metrics
Invoke-RestMethod -Uri "http://localhost:8081/metrics"
Invoke-RestMethod -Uri "http://localhost:8082/metrics"
Invoke-RestMethod -Uri "http://localhost:8083/metrics"
```

### View Internal Stats
```powershell
# Internal statistics
Invoke-RestMethod -Uri "http://localhost:8081/internal/stats"
Invoke-RestMethod -Uri "http://localhost:8082/internal/stats"
Invoke-RestMethod -Uri "http://localhost:8083/internal/stats"
```

---

## View Logs

```powershell
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f user-service
docker-compose logs -f order-service
docker-compose logs -f payment-service

# Kafka topics
docker-compose logs -f kafka
```

---

## Stop the Project

```powershell
# Stop and remove containers (keeps data)
docker-compose down

# Stop and remove containers + volumes (deletes all data)
docker-compose down -v
```

---

## Troubleshooting

### Port Already in Use
```powershell
# Find and kill processes using the ports
netstat -ano | findstr :8081
netstat -ano | findstr :8082
netstat -ano | findstr :8083

# Kill process by PID
Stop-Process -ID <PID> -Force
```

### Docker Not Running
```powershell
# Start Docker Desktop
& "C:\Program Files\Docker\Docker\Docker Desktop.exe"

# Or via command line
Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
```

### Services Not Starting
```powershell
# Full reset
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
```

### Check Kafka Topics
```powershell
# List topics
docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 --list

# Should show: user.created, order.created, payment.completed
```

### Check Databases
```powershell
# PostgreSQL - Users
docker-compose exec postgres psql -U postgres -d userdb -c "SELECT * FROM users;"

# PostgreSQL - Orders
docker-compose exec postgres psql -U postgres -d orderdb -c "SELECT * FROM orders;"

# MongoDB
docker-compose exec mongodb mongosh paymentdb --eval "db.payments.find()"
```

---

## Project URLs

Once running, access these URLs:

| Service | Endpoint | URL |
|---------|----------|-----|
| User Service | Health | http://localhost:8081/health |
| User Service | Create User | http://localhost:8081/users |
| User Service | Metrics | http://localhost:8081/metrics |
| Order Service | Health | http://localhost:8082/health |
| Order Service | Create Order | http://localhost:8082/orders |
| Order Service | User Cache | http://localhost:8082/internal/cache |
| Payment Service | Health | http://localhost:8083/health |
| Payment Service | Get Payments | http://localhost:8083/payments |
| Payment Service | Metrics | http://localhost:8083/metrics |

---

## Development Mode (Without Docker)

If you want to run services locally for development:

### Prerequisites
- Go 1.21+
- PostgreSQL running locally
- MongoDB running locally
- Kafka running locally

### Start Infrastructure
```powershell
# Start PostgreSQL, MongoDB, Kafka locally
# Then set environment variables

$env:KAFKA_BROKERS="localhost:9092"
$env:USER_DB_URL="postgres://postgres:password@localhost:5432/userdb?sslmode=disable"
$env:ORDER_DB_URL="postgres://postgres:password@localhost:5432/orderdb?sslmode=disable"
$env:MONGO_URI="mongodb://localhost:27017"
```

### Run Services
```powershell
# Terminal 1 - User Service
cd user-service
go run main.go

# Terminal 2 - Order Service
cd order-service
go run main.go

# Terminal 3 - Payment Service
cd payment-service
go run main.go
```

---

## Quick Test Script

Save as `test-all.ps1` and run:

```powershell
# test-all.ps1
Write-Host "=== Testing Event-Driven Microservices ===" -ForegroundColor Green

# Test User Service
Write-Host "`n1. Creating User..." -ForegroundColor Yellow
$user = Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body '{"name": "Alice", "email": "alice@example.com"}'
Write-Host "   User ID: $($user.id)"

# Test Order Service
Write-Host "`n2. Creating Order..." -ForegroundColor Yellow
$order = Invoke-RestMethod -Uri "http://localhost:8082/orders" -Method POST -ContentType "application/json" -Body "{`"user_id`": `"$($user.id)`", `"product_name`": `"Laptop`", `"price`": 999}"
Write-Host "   Order ID: $($order.id)"

# Wait for payment
Write-Host "`n3. Waiting for payment processing..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

# Test Payment Service
Write-Host "`n4. Checking Payments..." -ForegroundColor Yellow
$payments = Invoke-RestMethod -Uri "http://localhost:8083/payments"
$payment = $payments | Where-Object { $_.order_id -eq $order.id }
if ($payment) {
    Write-Host "   Payment ID: $($payment._id)"
    Write-Host "   Status: $($payment.payment_status)"
} else {
    Write-Host "   Payment not found yet, try again in a few seconds" -ForegroundColor Red
}

# Health checks
Write-Host "`n5. Health Checks..." -ForegroundColor Yellow
$userHealth = Invoke-RestMethod -Uri "http://localhost:8081/health"
$orderHealth = Invoke-RestMethod -Uri "http://localhost:8082/health"
$paymentHealth = Invoke-RestMethod -Uri "http://localhost:8083/health"
Write-Host "   User Service: $($userHealth.status)"
Write-Host "   Order Service: $($orderHealth.status)"
Write-Host "   Payment Service: $($paymentHealth.status)"

Write-Host "`n=== Test Complete ===" -ForegroundColor Green
```

Run it:
```powershell
.\test-all.ps1
```

---

## Summary

| Step | Command |
|------|---------|
| **Start** | `docker-compose up -d` |
| **Verify** | `docker-compose ps` |
| **Test** | `Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST ...` |
| **Logs** | `docker-compose logs -f` |
| **Stop** | `docker-compose down` |
| **Reset** | `docker-compose down -v` |

---

## Next Steps

1. **Import Postman Collection** - Use `E-Commerce-Microservices-Postman-Collection.json`
2. **Read API Testing Guide** - See `API_TESTING.md`
3. **Read Test Cases** - See `POSTMAN_TEST_CASES.md`
4. **Deploy to Kubernetes** - See `k8s/` directory
