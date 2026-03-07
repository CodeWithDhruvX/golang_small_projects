# Postman API Test Cases

Complete Postman collection documentation for testing the Event-Driven E-Commerce Microservices.

## How to Use This Guide

1. Open **Postman**
2. Create a new **Collection** named "E-Commerce Microservices"
3. Add requests as documented below
4. Create **Environment** with variables: `userId`, `orderId`, `paymentId`
5. Run requests in sequence

---

## Collection Structure

```
E-Commerce Microservices
├── 1. User Service (Port 8081)
├── 2. Order Service (Port 8082)
├── 3. Payment Service (Port 8083)
└── 4. Complete Flow Test
```

---

## 1. User Service Tests

### 1.1 Create User

**Method:** POST  
**URL:** `http://localhost:8081/users`  
**Headers:**
```
Content-Type: application/json
```
**Body:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

**Tests (Postman Test Script):**
```javascript
// Test 1: Status code is 201 Created
pm.test('Status code is 201', function () {
    pm.response.to.have.status(201);
});

// Test 2: Response has Content-Type header
pm.test('Content-Type is application/json', function () {
    pm.response.to.have.header('Content-Type');
    pm.expect(pm.response.headers.get('Content-Type')).to.include('application/json');
});

// Test 3: Response has required fields
pm.test('Response has required fields', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('id');
    pm.expect(jsonData).to.have.property('name');
    pm.expect(jsonData).to.have.property('email');
    pm.expect(jsonData).to.have.property('created_at');
});

// Test 4: ID is a valid UUID
pm.test('ID is a valid UUID', function () {
    var jsonData = pm.response.json();
    var uuidPattern = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    pm.expect(jsonData.id).to.match(uuidPattern);
});

// Test 5: Name and email match request
pm.test('Name and email match request', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.name).to.eql('John Doe');
    pm.expect(jsonData.email).to.eql('john.doe@example.com');
});

// Save user ID for next requests
var jsonData = pm.response.json();
pm.collectionVariables.set('userId', jsonData.id);
console.log('Created user with ID:', jsonData.id);
```

**Expected Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "created_at": "2026-03-07T08:45:30.123456789Z"
}
```

---

### 1.2 Get User by ID

**Method:** GET  
**URL:** `http://localhost:8081/users/{{userId}}`

**Tests:**
```javascript
// Test 1: Status code is 200
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

// Test 2: Response contains user data
pm.test('Response contains user data', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('id');
    pm.expect(jsonData).to.have.property('name');
    pm.expect(jsonData).to.have.property('email');
    pm.expect(jsonData.id).to.eql(pm.collectionVariables.get('userId'));
});
```

---

### 1.3 Health Check

**Method:** GET  
**URL:** `http://localhost:8081/health`

**Tests:**
```javascript
// Test 1: Status code is 200
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

// Test 2: Service is healthy
pm.test('Service is healthy', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('healthy');
    pm.expect(jsonData.service).to.eql('user-service');
});

// Test 3: Database is connected
pm.test('Database is connected', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.database).to.eql('connected');
});
```

**Expected Response:**
```json
{
  "service": "user-service",
  "status": "healthy",
  "timestamp": "2026-03-07T08:45:30Z",
  "database": "connected"
}
```

---

### 1.4 Get Metrics

**Method:** GET  
**URL:** `http://localhost:8081/metrics`

**Tests:**
```javascript
// Test 1: Status code is 200
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

// Test 2: Response contains Prometheus metrics
pm.test('Response contains Prometheus metrics', function () {
    var responseText = pm.response.text();
    pm.expect(responseText).to.include('user_service_http_requests_total');
    pm.expect(responseText).to.include('user_service_db_queries_total');
    pm.expect(responseText).to.include('user_service_kafka_messages_produced_total');
});
```

---

### 1.5 Get Internal Stats

**Method:** GET  
**URL:** `http://localhost:8081/internal/stats`

**Tests:**
```javascript
// Test 1: Status code is 200
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

// Test 2: Response contains stats
pm.test('Response contains internal stats', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('service');
    pm.expect(jsonData).to.have.property('timestamp');
    pm.expect(jsonData).to.have.property('totalUsers');
});
```

---

## 2. Order Service Tests

### 2.1 Create Order

**Method:** POST  
**URL:** `http://localhost:8082/orders`  
**Headers:**
```
Content-Type: application/json
```
**Body:**
```json
{
  "user_id": "{{userId}}",
  "product_name": "Gaming Laptop",
  "price": 1299
}
```

**Tests:**
```javascript
// Test 1: Status code is 201 Created
pm.test('Status code is 201', function () {
    pm.response.to.have.status(201);
});

// Test 2: Response has required fields
pm.test('Response has required fields', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('id');
    pm.expect(jsonData).to.have.property('user_id');
    pm.expect(jsonData).to.have.property('product_name');
    pm.expect(jsonData).to.have.property('price');
    pm.expect(jsonData).to.have.property('status');
    pm.expect(jsonData).to.have.property('created_at');
});

// Test 3: ID is valid UUID
pm.test('Order ID is valid UUID', function () {
    var jsonData = pm.response.json();
    var uuidPattern = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    pm.expect(jsonData.id).to.match(uuidPattern);
});

// Test 4: User ID matches request
pm.test('User ID matches request', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.user_id).to.eql(pm.collectionVariables.get('userId'));
});

// Test 5: Order status is pending
pm.test('Order status is pending', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('pending');
});

// Test 6: Price is correct
pm.test('Price is correct', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.price).to.eql(1299);
});

// Save order ID for next requests
var jsonData = pm.response.json();
pm.collectionVariables.set('orderId', jsonData.id);
console.log('Created order with ID:', jsonData.id);
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

### 2.2 Get Order by ID

**Method:** GET  
**URL:** `http://localhost:8082/orders/{{orderId}}`

**Tests:**
```javascript
// Test 1: Status code is 200
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

// Test 2: Order data is correct
pm.test('Order data is correct', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.id).to.eql(pm.collectionVariables.get('orderId'));
    pm.expect(jsonData.user_id).to.eql(pm.collectionVariables.get('userId'));
    pm.expect(jsonData.product_name).to.eql('Gaming Laptop');
});
```

---

### 2.3 Health Check

**Method:** GET  
**URL:** `http://localhost:8082/health`

**Tests:**
```javascript
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

pm.test('Service is healthy', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('healthy');
    pm.expect(jsonData.service).to.eql('order-service');
    pm.expect(jsonData).to.have.property('cached_users');
});
```

---

### 2.4 Get Metrics

**Method:** GET  
**URL:** `http://localhost:8082/metrics`

**Tests:**
```javascript
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

pm.test('Contains Kafka consumer metrics', function () {
    var responseText = pm.response.text();
    pm.expect(responseText).to.include('order_service_kafka_messages_consumed_total');
    pm.expect(responseText).to.include('order_service_kafka_messages_produced_total');
    pm.expect(responseText).to.include('order_service_users_cached');
});
```

---

### 2.5 Get User Cache (Kafka Event Verification)

**Method:** GET  
**URL:** `http://localhost:8082/internal/cache`

**Tests:**
```javascript
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

pm.test('Cache contains user from Kafka', function () {
    var jsonData = pm.response.json();
    var userId = pm.collectionVariables.get('userId');
    pm.expect(jsonData).to.have.property(userId);
});
```

**What this verifies:** The Order Service consumed the `user.created` event from Kafka and cached the user.

---

### 2.6 Get Internal Stats

**Method:** GET  
**URL:** `http://localhost:8082/internal/stats`

**Tests:**
```javascript
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

pm.test('Contains order stats', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('totalOrders');
    pm.expect(jsonData).to.have.property('cachedUsers');
});
```

---

## 3. Payment Service Tests

### 3.1 Wait for Payment Processing

**Method:** GET  
**URL:** `http://localhost:8082/internal/stats`  
**Pre-request Script:**
```javascript
// Wait 3 seconds for payment to be processed
console.log('Waiting 3 seconds for payment processing...');
setTimeout(function() {}, 3000);
```

---

### 3.2 Get All Payments

**Method:** GET  
**URL:** `http://localhost:8083/payments`

**Tests:**
```javascript
// Test 1: Status code is 200
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

// Test 2: Response is an array
pm.test('Response is an array', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.be.an('array');
});

// Test 3: At least one payment exists
pm.test('At least one payment exists', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.length).to.be.at.least(1);
});

// Test 4: Payment has required fields
pm.test('Payment has required fields', function () {
    var jsonData = pm.response.json();
    var payment = jsonData[jsonData.length - 1]; // Get latest
    pm.expect(payment).to.have.property('_id');
    pm.expect(payment).to.have.property('order_id');
    pm.expect(payment).to.have.property('user_id');
    pm.expect(payment).to.have.property('amount');
    pm.expect(payment).to.have.property('payment_status');
});

// Test 5: Payment links to our order
pm.test('Payment links to our order', function () {
    var jsonData = pm.response.json();
    var orderId = pm.collectionVariables.get('orderId');
    var matchingPayment = jsonData.find(function(p) {
        return p.order_id === orderId;
    });
    pm.expect(matchingPayment).to.not.be.undefined;
    pm.collectionVariables.set('paymentId', matchingPayment._id);
    console.log('Found payment:', matchingPayment._id);
});

// Test 6: Payment amount matches order
pm.test('Payment amount matches order price', function () {
    var jsonData = pm.response.json();
    var orderId = pm.collectionVariables.get('orderId');
    var payment = jsonData.find(function(p) {
        return p.order_id === orderId;
    });
    pm.expect(payment.amount).to.eql(1299);
});

// Test 7: Payment status is valid
pm.test('Payment status is valid (success or failed)', function () {
    var jsonData = pm.response.json();
    var orderId = pm.collectionVariables.get('orderId');
    var payment = jsonData.find(function(p) {
        return p.order_id === orderId;
    });
    pm.expect(['success', 'failed']).to.include(payment.payment_status);
    console.log('Payment status:', payment.payment_status);
});
```

**What this verifies:** The Payment Service automatically consumed the `order.created` event and created a payment.

**Expected Response:**
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

---

### 3.3 Get Payment by ID

**Method:** GET  
**URL:** `http://localhost:8083/payments/{{paymentId}}`

**Tests:**
```javascript
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

pm.test('Payment ID matches', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData._id).to.eql(pm.collectionVariables.get('paymentId'));
});
```

---

### 3.4 Health Check

**Method:** GET  
**URL:** `http://localhost:8083/health`

**Tests:**
```javascript
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

pm.test('Service is healthy', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('healthy');
    pm.expect(jsonData.service).to.eql('payment-service');
    pm.expect(jsonData.database).to.eql('connected');
});
```

---

### 3.5 Get Metrics

**Method:** GET  
**URL:** `http://localhost:8083/metrics`

**Tests:**
```javascript
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

pm.test('Contains payment metrics', function () {
    var responseText = pm.response.text();
    pm.expect(responseText).to.include('payment_service_payments_processed_total');
    pm.expect(responseText).to.include('payment_service_kafka_messages_consumed_total');
    pm.expect(responseText).to.include('payment_service_db_queries_total');
});
```

---

### 3.6 Get Internal Stats

**Method:** GET  
**URL:** `http://localhost:8083/internal/stats`

**Tests:**
```javascript
pm.test('Status code is 200', function () {
    pm.response.to.have.status(200);
});

pm.test('Contains payment stats', function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('totalPayments');
    pm.expect(jsonData).to.have.property('paymentsByStatus');
    var statusCounts = jsonData.paymentsByStatus;
    pm.expect(statusCounts).to.have.property('success');
    pm.expect(statusCounts).to.have.property('failed');
});
```

---

## 4. Complete Flow Verification

### 4.1 Full Flow Summary

**Method:** GET  
**URL:** `http://localhost:8081/health`

**Test Script:**
```javascript
// Complete event-driven flow test
pm.test('Event-Driven Flow Complete', function () {
    console.log('========================================');
    console.log('EVENT-DRIVEN FLOW TEST SUMMARY');
    console.log('========================================');
    console.log('1. User created:', pm.collectionVariables.get('userId'));
    console.log('2. Order created:', pm.collectionVariables.get('orderId'));
    console.log('3. Payment created:', pm.collectionVariables.get('paymentId'));
    console.log('');
    console.log('Event Flow:');
    console.log('  User Service → user.created → Order Service (cached)');
    console.log('  Order Service → order.created → Payment Service');
    console.log('  Payment Service → payment.completed → Kafka');
    console.log('========================================');
});

// Verify all IDs were created
pm.test('All entities created', function () {
    pm.expect(pm.collectionVariables.get('userId')).to.not.be.undefined;
    pm.expect(pm.collectionVariables.get('orderId')).to.not.be.undefined;
    pm.expect(pm.collectionVariables.get('paymentId')).to.not.be.undefined;
});
```

---

## Environment Variables

Create these variables in your Postman Environment:

| Variable | Initial Value | Current Value |
|----------|---------------|---------------|
| `userId` | (empty) | (auto-filled) |
| `orderId` | (empty) | (auto-filled) |
| `paymentId` | (empty) | (auto-filled) |

---

## Execution Order

Run requests in this sequence:

1. **User Service**
   - Create User → sets `{{userId}}`
   - Get User by ID (verify storage)
   - Health Check
   - Get Metrics
   - Get Internal Stats

2. **Order Service**
   - Create Order → uses `{{userId}}`, sets `{{orderId}}`
   - Get Order by ID
   - Health Check
   - Get User Cache (verify Kafka consumption)
   - Get Metrics
   - Get Internal Stats

3. **Payment Service**
   - Wait for Payment Processing (delay)
   - Get All Payments → finds payment for `{{orderId}}`, sets `{{paymentId}}`
   - Get Payment by ID
   - Health Check
   - Get Metrics
   - Get Internal Stats

4. **Complete Flow**
   - Full Flow Summary

---

## Expected Event Flow

```
POST /users (User Service)
  ↓
  Stores in PostgreSQL
  ↓
  Publishes user.created → Kafka
  ↓
Order Service consumes user.created
  ↓
  Caches user in memory (verify via /internal/cache)
  ↓
POST /orders (Order Service)  
  ↓
  Stores in PostgreSQL
  ↓
  Publishes order.created → Kafka
  ↓
Payment Service consumes order.created
  ↓
  Processes payment (simulated)
  ↓
  Stores in MongoDB
  ↓
  Publishes payment.completed → Kafka
  ↓
GET /payments shows auto-created payment
```

---

## Quick Verification Commands

After running the collection, verify with PowerShell:

```powershell
# Check user exists
Invoke-RestMethod -Uri "http://localhost:8081/users/$env:userId"

# Check order exists
Invoke-RestMethod -Uri "http://localhost:8082/orders/$env:orderId"

# Check payment exists
Invoke-RestMethod -Uri "http://localhost:8083/payments"

# View all metrics
Invoke-RestMethod -Uri "http://localhost:8081/metrics"
Invoke-RestMethod -Uri "http://localhost:8082/metrics"
Invoke-RestMethod -Uri "http://localhost:8083/metrics"
```

---

## Troubleshooting

### Tests Failing?

1. **Check services are running:**
   ```powershell
   docker-compose ps
   ```

2. **Check logs:**
   ```powershell
   docker-compose logs user-service
   docker-compose logs order-service
   docker-compose logs payment-service
   ```

3. **Verify Kafka topics:**
   ```powershell
   docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 --list
   ```

### Variables Not Set?

- Ensure you run **Create User** before **Create Order**
- Ensure you run **Create Order** before **Get All Payments**
- Wait 3 seconds between order creation and payment check

---

## Test Summary

| Service | Requests | Key Tests |
|---------|----------|-----------|
| User | 5 | Create, Get, Health, Metrics, Stats |
| Order | 6 | Create, Get, Health, Cache, Metrics, Stats |
| Payment | 6 | Wait, List, Get, Health, Metrics, Stats |
| Flow | 1 | Summary verification |
| **Total** | **18** | **100+ assertions** |
