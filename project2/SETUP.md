# Event-Driven E-Commerce Microservices System - Implementation Guide

A complete event-driven microservices architecture using Go, Kafka, PostgreSQL, MongoDB, Docker, and Kubernetes.

## Architecture Overview

```
Client
   |
   v
API Gateway / REST
   |
   v
User Service (PostgreSQL) ──user.created──> Kafka
   |                                              |
   |                                              v
   |                                        Order Service (PostgreSQL)
   |                                              |
   |                                              | order.created
   |                                              v
   |                                        Payment Service (MongoDB)
   |                                              |
   |                                              | payment.completed
   |                                              v
   |                                        Kafka
```

## Services

### 1. User Service (Port 8081)
- **Database**: PostgreSQL
- **Responsibilities**: User registration, profile management
- **API Endpoints**:
  - `POST /users` - Create a new user
  - `GET /users/{id}` - Get user by ID
  - `GET /health` - Health check
  - `GET /metrics` - Prometheus metrics
  - `GET /internal/stats` - Internal statistics
- **Kafka**: Produces `user.created` events

### 2. Order Service (Port 8082)
- **Database**: PostgreSQL
- **Responsibilities**: Order creation, consumes user events
- **API Endpoints**:
  - `POST /orders` - Create a new order
  - `GET /orders/{id}` - Get order by ID
  - `GET /health` - Health check
  - `GET /metrics` - Prometheus metrics
  - `GET /internal/stats` - Internal statistics
  - `GET /internal/cache` - View cached users from Kafka
- **Kafka**: Consumes `user.created`, Produces `order.created`

### 3. Payment Service (Port 8083)
- **Database**: MongoDB
- **Responsibilities**: Payment processing
- **API Endpoints**:
  - `GET /payments` - Get all payments
  - `GET /payments/{id}` - Get payment by ID
  - `GET /health` - Health check
  - `GET /metrics` - Prometheus metrics
  - `GET /internal/stats` - Internal statistics
- **Kafka**: Consumes `order.created`, Produces `payment.completed`

## Event Flow

1. **User Registration**: Client → `POST /users` → User Service publishes `user.created`
2. **Order Creation**: Client → `POST /orders` → Order Service publishes `order.created`
3. **Payment Processing**: Payment Service consumes `order.created` → processes payment → publishes `payment.completed`

## Quick Start (Docker Compose)

### Prerequisites
- Docker Desktop
- Docker Compose

### Run the System

```powershell
# Build all services
docker-compose build

# Start all services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f

# Stop all services
docker-compose down -v
```

### Test the Flow

```powershell
# 1. Create a user
$userResponse = Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body '{"name": "John Doe", "email": "john@example.com"}'
Write-Host "Created user: $($userResponse.id)"

# 2. Create an order
$orderBody = @{
    user_id = $userResponse.id
    product_name = "Laptop"
    price = 1200
} | ConvertTo-Json

$orderResponse = Invoke-RestMethod -Uri "http://localhost:8082/orders" -Method POST -ContentType "application/json" -Body $orderBody
Write-Host "Created order: $($orderResponse.id)"

# 3. Check payments (automatically created)
Invoke-RestMethod -Uri "http://localhost:8083/payments" | Format-List

# 4. Check health of services
Invoke-RestMethod -Uri "http://localhost:8081/health"
Invoke-RestMethod -Uri "http://localhost:8082/health"
Invoke-RestMethod -Uri "http://localhost:8083/health"

# 5. View internal stats
Invoke-RestMethod -Uri "http://localhost:8081/internal/stats"
Invoke-RestMethod -Uri "http://localhost:8082/internal/stats"
Invoke-RestMethod -Uri "http://localhost:8083/internal/stats"

# 6. View Prometheus metrics
Invoke-RestMethod -Uri "http://localhost:8081/metrics"
```

## Kubernetes Deployment

### Prerequisites
- kubectl
- minikube or any Kubernetes cluster

### Deploy to Kubernetes

```powershell
# Apply all manifests
kubectl apply -f k8s/

# Wait for pods to be ready
kubectl wait --for=condition=ready pod -l app=user-service -n ecommerce --timeout=120s
kubectl wait --for=condition=ready pod -l app=order-service -n ecommerce --timeout=120s
kubectl wait --for=condition=ready pod -l app=payment-service -n ecommerce --timeout=120s

# Check status
kubectl get all -n ecommerce

# Port forward to access services
kubectl port-forward svc/user-service 8081:8081 -n ecommerce
kubectl port-forward svc/order-service 8082:8082 -n ecommerce
kubectl port-forward svc/payment-service 8083:8083 -n ecommerce

# Delete all resources
kubectl delete -f k8s/
```

## Monitoring & Observability

Each service exposes Prometheus metrics at `/metrics`:

- `*_http_requests_total` - HTTP request count by method, endpoint, status
- `*_http_request_duration_seconds` - HTTP request duration histogram
- `*_db_queries_total` - Database query count by operation, status
- `*_kafka_messages_produced_total` - Kafka messages produced by topic
- `*_kafka_messages_consumed_total` - Kafka messages consumed by topic
- `*_active_connections` - Active HTTP connections gauge

### Internal Statistics Endpoints

- `/internal/stats` - Service statistics including DB counts
- `/internal/cache` - Order Service user cache view (shows users consumed from Kafka)

## Project Structure

```
project2/
├── docker-compose.yml          # Local development orchestration
├── init-postgres.sh           # PostgreSQL initialization
├── readme.md                  # Requirements specification
├── SETUP.md                   # This file - Implementation guide
├── user-service/              # User microservice
│   ├── main.go               # Service implementation with Kafka producer, DB, metrics
│   ├── Dockerfile            # Multi-stage container image
│   ├── go.mod                # Go dependencies
│   └── go.sum                # Go checksums
├── order-service/             # Order microservice
│   ├── main.go               # Service with Kafka consumer/producer, DB, metrics
│   ├── Dockerfile            # Multi-stage container image
│   ├── go.mod                # Go dependencies
│   └── go.sum                # Go checksums
├── payment-service/           # Payment microservice
│   ├── main.go               # Service with Kafka consumer/producer, MongoDB, metrics
│   ├── Dockerfile            # Multi-stage container image
│   ├── go.mod                # Go dependencies
│   └── go.sum                # Go checksums
└── k8s/                       # Kubernetes manifests
    ├── 01-namespace.yaml      # Namespace definition
    ├── 02-config.yaml         # ConfigMaps and Secrets
    ├── 03-zookeeper.yaml      # Zookeeper StatefulSet
    ├── 04-kafka.yaml          # Kafka StatefulSet
    ├── 05-postgres.yaml       # PostgreSQL StatefulSet
    ├── 06-mongodb.yaml        # MongoDB StatefulSet
    ├── 07-user-service.yaml   # User Service Deployment
    ├── 08-order-service.yaml  # Order Service Deployment
    ├── 09-payment-service.yaml # Payment Service Deployment
    └── 10-kafka-init.yaml     # Kafka topic initialization Job
```

## Kafka Topics

| Topic              | Producer         | Consumer         | Partitions |
|-------------------|------------------|------------------|------------|
| user.created      | User Service     | Order Service    | 3          |
| order.created     | Order Service    | Payment Service  | 3          |
| payment.completed | Payment Service  | Analytics        | 3          |

## Environment Variables

| Variable       | Description                    | Default                           |
|----------------|--------------------------------|-----------------------------------|
| KAFKA_BROKERS  | Kafka broker addresses         | localhost:9092                    |
| USER_DB_URL    | User Service PostgreSQL URL    | postgres://.../userdb             |
| ORDER_DB_URL   | Order Service PostgreSQL URL   | postgres://.../orderdb            |
| MONGO_URI      | Payment Service MongoDB URI    | mongodb://localhost:27017         |
| PORT           | Service HTTP port              | 8081/8082/8083                    |

## Implementation Details

### User Service Features
- UUID generation for user IDs
- PostgreSQL table auto-creation
- Kafka synchronous producer with retry
- Prometheus metrics for HTTP requests, DB queries, Kafka messages
- Health check endpoint with DB connectivity check

### Order Service Features
- Kafka consumer group for `user.created` events
- In-memory user cache from Kafka events
- Dual role: consumes user events, produces order events
- User cache viewer endpoint for debugging
- PostgreSQL persistence

### Payment Service Features
- MongoDB with automatic index creation
- Kafka consumer group for `order.created` events
- Automatic payment processing (simulated 95% success rate)
- Payment status tracking (success/failed)
- Aggregation queries for payment statistics

### Docker Configuration
- Multi-stage builds for minimal image size
- Alpine Linux base images
- Automatic Go module download with retry
- Health checks in Kubernetes manifests

### Kubernetes Configuration
- Namespace isolation (`ecommerce`)
- ConfigMaps for environment variables
- Secrets for sensitive data
- StatefulSets for databases and Kafka
- Deployments for microservices (2 replicas each)
- Services for internal communication
- Liveness and readiness probes

## Troubleshooting

### Services fail to start
```powershell
# Check Docker status
docker info

# View logs for specific service
docker-compose logs user-service
docker-compose logs order-service
docker-compose logs payment-service

# Restart services
docker-compose restart

# Full reset
docker-compose down -v
docker-compose up -d
```

### Kafka connection issues
```powershell
# Check Kafka is running
docker-compose ps

# List Kafka topics
docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 --list

# Check topic contents
docker exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic user.created --from-beginning
```

### Database connection issues
```powershell
# Check PostgreSQL
docker-compose exec postgres psql -U postgres -l

# Check specific database
docker-compose exec postgres psql -U postgres -d userdb -c "SELECT * FROM users;"
docker-compose exec postgres psql -U postgres -d orderdb -c "SELECT * FROM orders;"

# Check MongoDB
docker-compose exec mongodb mongosh --eval "show dbs"
docker-compose exec mongodb mongosh paymentdb --eval "db.payments.find()"
```

### Port conflicts
```powershell
# Find processes using ports
netstat -ano | findstr :8081
netstat -ano | findstr :8082
netstat -ano | findstr :8083

# Kill processes
Stop-Process -ID <PID>
```

## Architecture Patterns Demonstrated

1. **Event-Driven Architecture**: Services communicate asynchronously via Kafka events
2. **Database per Service**: Each service owns its own data store
3. **Asynchronous Communication**: Non-blocking message passing
4. **Polyglot Persistence**: PostgreSQL (SQL) + MongoDB (NoSQL) based on use case
5. **Containerization**: Docker containers for consistency
6. **Orchestration**: Kubernetes for production deployment
7. **Observability**: Prometheus metrics for monitoring
8. **Health Checking**: Readiness and liveness probes

## Development Commands

```powershell
# Run individual service locally (from service directory)
cd user-service
go mod tidy
go run main.go

# Run with environment variables
$env:KAFKA_BROKERS="localhost:9092"
$env:USER_DB_URL="postgres://postgres:password@localhost:5432/userdb?sslmode=disable"
go run main.go
```

## Next Steps

1. Start Docker Desktop
2. Run `docker-compose up -d` to start all services
3. Test the API endpoints using the PowerShell commands above
4. View metrics at `http://localhost:8081/metrics`, etc.
5. Deploy to Kubernetes using `kubectl apply -f k8s/`

## License

MIT License - Feel free to use this project for learning and development.
