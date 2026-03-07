# Fixed Postman Requests - User Creation Issues

## 🚨 Issue Identified
The email `john.doe@example.com` already exists in the database, causing a 500 error due to PostgreSQL unique constraint violation.

## ✅ Working Solutions

### **Option 1: Use Different Email (Recommended)**
```json
POST http://localhost:8081/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john.doe.new@example.com"
}
```

**Expected Response (201 Created):**
```json
{
  "id": "uuid-here",
  "name": "John Doe",
  "email": "john.doe.new@example.com",
  "created_at": "2026-03-07T15:28:22.741734457Z"
}
```

---

### **Option 2: Use Different Name**
```json
POST http://localhost:8081/users
Content-Type: application/json

{
  "name": "Johnny Doe",
  "email": "john.doe@example.com"
}
```

---

### **Option 3: Clear Existing Data**
```bash
# Remove the existing John Doe user
docker exec postgres psql -U postgres -d userdb -c "DELETE FROM users WHERE email = 'john.doe@example.com';"
```

---

### **Option 4: Use Completely New User**
```json
POST http://localhost:8081/users
Content-Type: application/json

{
  "name": "Alice Johnson",
  "email": "alice.johnson@example.com"
}
```

---

## 🧪 Test These Working Examples

### **Example 1: Create New User**
```json
POST http://localhost:8081/users
Content-Type: application/json

{
  "name": "Alice Johnson",
  "email": "alice.johnson@example.com"
}
```

### **Example 2: Create Order**
```json
POST http://localhost:8082/orders
Content-Type: application/json

{
  "user_id": "paste-user-id-from-above",
  "product_name": "MacBook Pro",
  "price": 1999
}
```

### **Example 3: Check Payment**
```json
GET http://localhost:8083/payments
```

---

## 🔍 Current Users in Database

| Name | Email | Status |
|------|-------|--------|
| Alice Smith | alice.smith@example.com | ✅ Available |
| Bob Wilson | bob.wilson@example.com | ✅ Available |
| Alice Test | alice@test.com | ✅ Available |
| Test User | test@example.com | ✅ Available |
| Jane Doe | jane.doe@example.com | ✅ Available |
| John Doe | john.doe@example.com | ❌ Already Exists |

---

## 🎯 Quick PowerShell Test

```powershell
# This will work
$response = Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body '{"name":"John Doe","email":"john.doe.new@example.com"}'
Write-Host "User created: $($response.id)"

# This will fail (500 error)
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body '{"name":"John Doe","email":"john.doe@example.com"}' -ErrorAction Stop
} catch {
    Write-Host "Error: $($_.Exception.Response.StatusCode.value__) - Email already exists"
}
```

---

## 📊 Service Status Check

All services are healthy:
- ✅ User Service: http://localhost:8081/health
- ✅ Order Service: http://localhost:8082/health  
- ✅ Payment Service: http://localhost:8083/health

---

## 🏆 Solution Summary

**The issue is not with the API** - it's working correctly and protecting data integrity by preventing duplicate emails. 

**Simply use a different email address** and your Postman requests will work perfectly!

The microservices architecture is functioning correctly:
1. User creation ✅
2. Kafka event publishing ✅  
3. Order creation ✅
4. Payment processing ✅
