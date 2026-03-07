Below is the **final technical requirements specification** for the **3-microservice Golang Kafka project** with:

* **Kafka event streaming**
* **PostgreSQL + MongoDB**
* **Docker**
* **Kubernetes**
* **Event-driven microservices**

This is written like a **project requirements document** you could submit for architecture/design.

---

# Project: Event-Driven E-Commerce Microservices System

## Architecture Overview

```text
Client
   |
   v
API Gateway / REST
   |
   v
User Service (PostgreSQL)
   |
   | user.created
   v
Kafka
   |
   v
Order Service (PostgreSQL)
   |
   | order.created
   v
Kafka
   |
   v
Payment Service (MongoDB)
   |
   | payment.completed
   v
Kafka
```

Architecture pattern:

* **Event-driven microservices**
* **Database per service**
* **Asynchronous communication via Kafka**

---

# Microservices in the System

## 1️⃣ User Service

### Responsibilities

* User registration
* User profile management
* Produce Kafka event when user is created

### API Endpoints

```
POST /users
GET /users/{id}
```

### Database

**PostgreSQL**

Table:

```sql
users
```

Fields:

| Field      | Type      |
| ---------- | --------- |
| id         | UUID      |
| name       | String    |
| email      | String    |
| created_at | Timestamp |

### Kafka

Produces:

```
user.created
```

Consumes:

```
None
```

---

# 2️⃣ Order Service

### Responsibilities

* Create orders
* Store order information
* Consume user events
* Produce order events

### API Endpoints

```
POST /orders
GET /orders/{id}
```

### Database

**PostgreSQL**

Table:

```
orders
```

Fields:

| Field        | Type      |
| ------------ | --------- |
| id           | UUID      |
| user_id      | UUID      |
| product_name | String    |
| price        | Integer   |
| status       | String    |
| created_at   | Timestamp |

### Kafka

Consumes:

```
user.created
```

Produces:

```
order.created
```

---

# 3️⃣ Payment Service

### Responsibilities

* Process payments
* Store payment transactions
* Consume order events
* Produce payment events

### Database

**MongoDB (NoSQL)**

Collection:

```
payments
```

Example document:

```json
{
  "_id": "uuid",
  "order_id": "uuid",
  "amount": 1200,
  "payment_status": "success",
  "created_at": "timestamp"
}
```

### Kafka

Consumes:

```
order.created
```

Produces:

```
payment.completed
```

---

# Database Architecture

| Service         | Database   | Type  |
| --------------- | ---------- | ----- |
| User Service    | PostgreSQL | SQL   |
| Order Service   | PostgreSQL | SQL   |
| Payment Service | MongoDB    | NoSQL |

Pattern used:

**Database per microservice**

---

# Kafka Topics

| Topic             | Producer        | Consumer               |
| ----------------- | --------------- | ---------------------- |
| user.created      | User Service    | Order Service          |
| order.created     | Order Service   | Payment Service        |
| payment.completed | Payment Service | Notification/Analytics |

---

# Technology Stack

## Programming Language

| Technology  | Version |
| ----------- | ------- |
| Go (Golang) | 1.20+   |

---

## Messaging System

| Component    | Version |
| ------------ | ------- |
| Apache Kafka | 3.x     |
| Zookeeper    | 3.x     |

---

## Databases

| Database   | Version |
| ---------- | ------- |
| PostgreSQL | 14+     |
| MongoDB    | 6+      |

---

## Containerization

| Tool           | Purpose                       |
| -------------- | ----------------------------- |
| Docker         | Containerize microservices    |
| Docker Compose | Local development environment |

Docker images required:

```
user-service
order-service
payment-service
kafka
zookeeper
postgres
mongodb
```

---

# Kubernetes Deployment

Kubernetes will be used for **production orchestration**.

## Kubernetes Components

| Component   | Purpose                  |
| ----------- | ------------------------ |
| Deployment  | Run microservices        |
| Service     | Expose internal services |
| ConfigMap   | Store configuration      |
| Secret      | Store credentials        |
| StatefulSet | Kafka / databases        |
| Ingress     | External API access      |

---

# Kubernetes Workloads

| Component       | Kubernetes Resource |
| --------------- | ------------------- |
| User Service    | Deployment          |
| Order Service   | Deployment          |
| Payment Service | Deployment          |
| Kafka           | StatefulSet         |
| Zookeeper       | StatefulSet         |
| PostgreSQL      | StatefulSet         |
| MongoDB         | StatefulSet         |

---

# Ports Configuration

| Component       | Port  |
| --------------- | ----- |
| User Service    | 8081  |
| Order Service   | 8082  |
| Payment Service | 8083  |
| Kafka           | 9092  |
| Zookeeper       | 2181  |
| PostgreSQL      | 5432  |
| MongoDB         | 27017 |

---

# Required Go Libraries

Kafka client:

```
github.com/IBM/sarama
```

PostgreSQL driver:

```
github.com/lib/pq
```

MongoDB driver:

```
go.mongodb.org/mongo-driver/mongo
```

UUID generation:

```
github.com/google/uuid
```

Optional ORM:

```
gorm.io/gorm
gorm.io/driver/postgres
```

---

# Configuration Management

Environment variables:

```
KAFKA_BROKERS=localhost:9092
USER_DB_URL=postgres://user:pass@postgres:5432/userdb
ORDER_DB_URL=postgres://user:pass@postgres:5432/orderdb
MONGO_URI=mongodb://mongo:27017
```

---

# Docker Compose (Local Development)

Docker Compose will run:

```
zookeeper
kafka
postgres
mongodb
user-service
order-service
payment-service
```

Purpose:

* Local testing
* Integration development

---

# Kubernetes Deployment (Production)

Kubernetes cluster runs:

```
Kafka cluster
Zookeeper
PostgreSQL
MongoDB
User Service
Order Service
Payment Service
API Gateway
```

Benefits:

* Auto scaling
* High availability
* Service discovery
* Rolling deployments

---

# Event Flow

### Step 1 — User Registration

Client:

```
POST /users
```

User Service:

```
Store in PostgreSQL
Publish event → user.created
```

---

### Step 2 — Order Creation

Client:

```
POST /orders
```

Order Service:

```
Store order
Publish event → order.created
```

---

### Step 3 — Payment Processing

Payment Service:

```
Consume order.created
Process payment
Store in MongoDB
Publish payment.completed
```

---

# Security Requirements

Recommended security features:

* TLS for Kafka
* JWT authentication
* Secure database credentials
* Kubernetes secrets

---

# Monitoring & Observability

Recommended tools:

| Tool          | Purpose       |
| ------------- | ------------- |
| Prometheus    | Metrics       |
| Grafana       | Visualization |
| Loki          | Logging       |
| OpenTelemetry | Tracing       |

---

# Testing Requirements

| Test Type         | Tool           |
| ----------------- | -------------- |
| Unit tests        | Go testing     |
| API testing       | Postman        |
| Kafka testing     | Testcontainers |
| Integration tests | Docker Compose |

---

# CI/CD Pipeline (Optional)

Typical pipeline:

```
GitHub
   ↓
Build Docker Images
   ↓
Push to Container Registry
   ↓
Deploy to Kubernetes
```

Tools:

* GitHub Actions
* GitLab CI
* ArgoCD

---

# Final Architecture Summary

| Layer            | Technology           |
| ---------------- | -------------------- |
| Backend          | Go microservices     |
| Messaging        | Kafka                |
| Databases        | PostgreSQL + MongoDB |
| Containerization | Docker               |
| Orchestration    | Kubernetes           |
| Monitoring       | Prometheus + Grafana |

---

✅ This project demonstrates:

* **Event-driven architecture**
* **Microservices communication using Kafka**
* **Polyglot persistence (SQL + NoSQL)**
* **Containerization with Docker**
* **Orchestration with Kubernetes**

---

If you'd like, I can also give you a **much stronger “industry-level” version** of this project that includes:

* **API Gateway**
* **Redis caching**
* **Kafka retry topics + DLQ**
* **Saga pattern for distributed transactions**
* **Complete Kubernetes manifests**

That version is **closer to real systems used by companies like Netflix, Uber, and LinkedIn**.



# Start all services
docker-compose up -d

# Wait 30 seconds
Start-Sleep -Seconds 30

# Check health manually
Invoke-RestMethod -Uri "http://localhost:8081/health"
Invoke-RestMethod -Uri "http://localhost:8082/health"
Invoke-RestMethod -Uri "http://localhost:8083/health"