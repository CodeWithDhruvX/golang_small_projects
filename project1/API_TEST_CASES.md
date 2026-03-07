# API Test Cases - Order Service

## Base URL
```
http://localhost:8080
```

## Endpoints
- `POST /orders` - Create a new order

---

## Test Case 1: Legacy Format (Simple Order)
**Request:**
```bash
POST http://localhost:8080/orders
Content-Type: application/json

{
  "item": "Laptop",
  "quantity": 1,
  "price": 1286.58
}
```

**Expected Response (201 Created):**
```json
{
  "id": "generated-uuid",
  "customer_id": "legacy-customer",
  "items": [
    {
      "product_id": "Laptop",
      "quantity": 1,
      "price": 1286.58
    }
  ],
  "total_amount": 1286.58,
  "item": "Laptop",
  "quantity": 1,
  "price": 1286.58,
  "status": "PENDING",
  "created_at": "2026-03-07T07:06:41Z"
}
```

---

## Test Case 2: New Format (Multi-item Order)
**Request:**
```bash
POST http://localhost:8080/orders
Content-Type: application/json

{
  "customer_id": "john_doe",
  "items": [
    {
      "product_id": "laptop",
      "quantity": 1,
      "price": 1286.58
    },
    {
      "product_id": "mouse",
      "quantity": 2,
      "price": 25.99
    }
  ],
  "total_amount": 1338.56
}
```

**Expected Response (201 Created):**
```json
{
  "id": "generated-uuid",
  "customer_id": "john_doe",
  "items": [
    {
      "product_id": "laptop",
      "quantity": 1,
      "price": 1286.58
    },
    {
      "product_id": "mouse",
      "quantity": 2,
      "price": 25.99
    }
  ],
  "total_amount": 1338.56,
  "status": "PENDING",
  "created_at": "2026-03-07T07:06:45Z"
}
```

---

## Test Case 3: Empty Items Array
**Request:**
```bash
POST http://localhost:8080/orders
Content-Type: application/json

{
  "customer_id": "test_user",
  "items": [],
  "total_amount": 0
}
```

**Expected Response (201 Created):**
```json
{
  "id": "generated-uuid",
  "customer_id": "test_user",
  "items": [],
  "total_amount": 0,
  "status": "PENDING",
  "created_at": "2026-03-07T07:10:00Z"
}
```

---

## Test Case 4: Invalid JSON
**Request:**
```bash
POST http://localhost:8080/orders
Content-Type: application/json

{
  "item": "Laptop",
  "quantity": "invalid",
  "price": 1286.58
}
```

**Expected Response (400 Bad Request):**
```
Bad request
```

---

## Test Case 5: Missing Required Fields (Legacy)
**Request:**
```bash
POST http://localhost:8080/orders
Content-Type: application/json

{
  "item": "Laptop"
}
```

**Expected Response (201 Created):**
```json
{
  "id": "generated-uuid",
  "customer_id": "legacy-customer",
  "items": [
    {
      "product_id": "Laptop",
      "quantity": 0,
      "price": 0
    }
  ],
  "total_amount": 0,
  "item": "Laptop",
  "quantity": 0,
  "price": 0,
  "status": "PENDING",
  "created_at": "2026-03-07T07:12:00Z"
}
```

---

## Test Case 6: Method Not Allowed
**Request:**
```bash
GET http://localhost:8080/orders
```

**Expected Response (405 Method Not Allowed):**
```
Method not allowed
```

---

## Test Case 7: Large Order
**Request:**
```bash
POST http://localhost:8080/orders
Content-Type: application/json

{
  "customer_id": "bulk_buyer",
  "items": [
    {
      "product_id": "laptop",
      "quantity": 10,
      "price": 999.99
    },
    {
      "product_id": "monitor",
      "quantity": 10,
      "price": 299.99
    },
    {
      "product_id": "keyboard",
      "quantity": 20,
      "price": 49.99
    }
  ],
  "total_amount": 18499.60
}
```

**Expected Response (201 Created):**
```json
{
  "id": "generated-uuid",
  "customer_id": "bulk_buyer",
  "items": [
    {
      "product_id": "laptop",
      "quantity": 10,
      "price": 999.99
    },
    {
      "product_id": "monitor",
      "quantity": 10,
      "price": 299.99
    },
    {
      "product_id": "keyboard",
      "quantity": 20,
      "price": 49.99
    }
  ],
  "total_amount": 18499.60,
  "status": "PENDING",
  "created_at": "2026-03-07T07:15:00Z"
}
```

---

## cURL Commands

### Test Case 1 (Legacy):
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"item": "Laptop", "quantity": 1, "price": 1286.58}'
```

### Test Case 2 (New Format):
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id": "john_doe", "items": [{"product_id": "laptop", "quantity": 1, "price": 1286.58}, {"product_id": "mouse", "quantity": 2, "price": 25.99}], "total_amount": 1338.56}'
```

### Test Case 6 (GET - Should Fail):
```bash
curl -X GET http://localhost:8080/orders
```

---

## PowerShell Commands

### Test Case 1 (Legacy):
```powershell
Invoke-RestMethod -Uri http://localhost:8080/orders -Method POST -ContentType "application/json" -Body '{"item": "Laptop", "quantity": 1, "price": 1286.58}'
```

### Test Case 2 (New Format):
```powershell
Invoke-RestMethod -Uri http://localhost:8080/orders -Method POST -ContentType "application/json" -Body '{"customer_id": "john_doe", "items": [{"product_id": "laptop", "quantity": 1, "price": 1286.58}, {"product_id": "mouse", "quantity": 2, "price": 25.99}], "total_amount": 1338.56}'
```

---

## Postman Collection

You can import these test cases into Postman:

1. Create a new collection called "Order Service API"
2. Add the following requests:

### Request 1: Create Order (Legacy)
- **Method**: POST
- **URL**: http://localhost:8080/orders
- **Headers**: Content-Type: application/json
- **Body** (raw JSON):
```json
{
  "item": "Laptop",
  "quantity": 1,
  "price": 1286.58
}
```

### Request 2: Create Order (New Format)
- **Method**: POST
- **URL**: http://localhost:8080/orders
- **Headers**: Content-Type: application/json
- **Body** (raw JSON):
```json
{
  "customer_id": "john_doe",
  "items": [
    {
      "product_id": "laptop",
      "quantity": 1,
      "price": 1286.58
    },
    {
      "product_id": "mouse",
      "quantity": 2,
      "price": 25.99
    }
  ],
  "total_amount": 1338.56
}
```

---

## Expected Behavior

✅ **Success Cases**: Return 201 Created with order details
✅ **Legacy Support**: Old format automatically converted to new structure
✅ **Error Handling**: Invalid JSON returns 400 Bad Request
✅ **Method Validation**: Non-POST methods return 405 Method Not Allowed
✅ **Auto-generated Fields**: ID, status, and timestamp are automatically added
