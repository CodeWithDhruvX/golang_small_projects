# Kafka Microservices (In-House)

This project contains an implementation of the Golang Kafka Microservice interview project. It operates completely in-house without any cloud dependencies by running local Docker containers for the messaging system.

## Setup Instructions

### 1. Start Infrastructure
Run Kafka, Zookeeper, and Kafka UI using Docker Compose. If you have Docker running, you can use the provided script to start everything at once:

```bash
.\start.ps1
```

Or manually:
```bash
docker compose up -d
```
You can access Kafka UI at `http://localhost:8081`.

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Run Microservices
If you didn't use `start.ps1`, open two different terminals.

**Terminal 1: Notification Service (Consumer)**
```bash
go run cmd/notification-service/main.go
```

**Terminal 2: Order Service (Producer API)**
```bash
go run cmd/order-service/main.go
```

### 4. Test the API
You can test the Order Service API using `curl` or Postman:
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"item": "Laptop", "quantity": 1, "price": 1200.50}'
```

You should see an HTTP 201 response and verify the event logged by the Notification Service in its terminal.

---

## Comprehensive Testing Guide

This project is built with testability in mind, decoupling the business logic and Kafka integrations to enable comprehensive Mock testing.

### Running All Tests

To run the entire test suite including the HTTP API and Kafka mock tests, run:
```bash
go test ./... -v
```

### Test Coverage Breakdown

#### 1. API HTTP Handler Tests (`internal/api/handler_test.go`)
The `order-service` HTTP handler is tested using Go's `httptest` package. It injects a `MockPublisher` to prevent real messages from being sent to Kafka during testing.
- **Success Case**: Ensures `HTTP 201 Created` is returned, a valid UUID is assigned to the `Order`, and the message is published.
- **Method Not Allowed**: Validates the handler rejects `GET` and other non-`POST` requests with `HTTP 405`.
- **Bad Request**: Validates the handler gracefully handles invalid JSON payloads with `HTTP 400`.
- **Kafka Down Scenario**: Simulates a producer failure and verifies the API returns `HTTP 500 Internal Server Error`.

#### 2. Kafka Producer Tests (`internal/kafka/producer_test.go`)
Testing Kafka publishing logic without an actual Kafka broker is achieved using `github.com/IBM/sarama/mocks`.
- A `sarama.MockSyncProducer` is instantiated.
- It expects a message to be published securely.
- The wrapper `PublishEvent` method serializes the structs correctly and bridges to the Sarama publisher.

#### 3. Kafka Consumer Tests (`internal/kafka/consumer_test.go`)
Tests the consumer processing logic isolated from network bindings.
- Mocks both `sarama.ConsumerGroupSession` and `sarama.ConsumerGroupClaim`.
- Pushes a simulated Kafka event into the mock channel.
- Validates the `ProcessFunc` correctly executes against the delivered message within a bounded timeframe context.

#### 4. Model Tests (`internal/models/order_test.go`)
Validates structural definitions.
- Ensures native JSON serialization and deserialization behaves predictably according to established JSON struct tags.
