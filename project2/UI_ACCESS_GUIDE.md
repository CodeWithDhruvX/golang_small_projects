# Database and Kafka UI Access Guide

## 🎯 Overview
This project now includes web-based UI tools for managing and monitoring PostgreSQL, MongoDB, and Kafka.

## 📊 UI Tools Available

### 1. **pgAdmin** - PostgreSQL Management
**URL:** http://localhost:5050  
**Login:** 
- Email: `admin@admin.com`
- Password: `admin123`

**Features:**
- Visual database browser
- SQL query editor
- Table design and management
- User management
- Performance monitoring

**Connection Details:**
- Host: `postgres`
- Port: `5432`
- Username: `postgres`
- Password: `password`
- Databases: `userdb`, `orderdb`

---

### 2. **MongoDB Express** - MongoDB Management
**URL:** http://localhost:8084  
**Login:** 
- Username: `admin`
- Password: `pass`

**Features:**
- Database and collection browser
- Document viewer and editor
- Index management
- Query builder
- Statistics dashboard

**Connection Details:**
- Host: `mongodb`
- Port: `27017`
- Database: `paymentdb`
- Collection: `payments`

---

### 3. **Kafka UI** - Kafka Cluster Management
**URL:** http://localhost:8085  
**Login:** No authentication required

**Features:**
- Topic management and browsing
- Consumer group monitoring
- Message inspection
- Broker status
- Partition details
- Real-time metrics

**Connection Details:**
- Cluster Name: `local`
- Bootstrap Servers: `kafka:29092`
- Zookeeper: `zookeeper:2181`

---

## 🚀 Quick Access Commands

### Start All Services with UI
```bash
docker-compose up -d
```

### Start Only UI Services
```bash
docker-compose up -d pgadmin mongo-express kafka-ui
```

### Check UI Services Status
```bash
docker-compose ps | grep -E "(pgadmin|mongo-express|kafka-ui)"
```

---

## 📋 Service URLs Summary

| Service | URL | Port | Purpose |
|---------|-----|------|---------|
| **pgAdmin** | http://localhost:5050 | 5050 | PostgreSQL GUI |
| **MongoDB Express** | http://localhost:8084 | 8084 | MongoDB GUI |
| **Kafka UI** | http://localhost:8085 | 8085 | Kafka GUI |
| **User Service** | http://localhost:8081 | 8081 | User API |
| **Order Service** | http://localhost:8082 | 8082 | Order API |
| **Payment Service** | http://localhost:8083 | 8083 | Payment API |

---

## 🔍 What You Can Do with Each UI

### PostgreSQL (pgAdmin)
```sql
-- View all users
SELECT * FROM userdb.users ORDER BY created_at DESC;

-- View all orders
SELECT * FROM orderdb.orders ORDER BY created_at DESC;

-- Check user-order relationships
SELECT u.name, u.email, o.product_name, o.price, o.status
FROM userdb.users u
JOIN orderdb.orders o ON u.id = o.user_id;
```

### MongoDB (Mongo Express)
- Browse `paymentdb` → `payments` collection
- View payment documents with structure:
```json
{
  "_id": "payment-uuid",
  "order_id": "order-uuid", 
  "user_id": "user-uuid",
  "amount": 1200,
  "payment_status": "success|failed",
  "created_at": "timestamp"
}
```

### Kafka (Kafka UI)
- **Topics:** `user.created`, `order.created`, `payment.completed`
- **Consumer Groups:** `order-service-group`, `payment-service-group`
- **Messages:** View real-time event flow
- **Lag Monitoring:** Check consumer processing status

---

## 🛠️ Troubleshooting

### UI Not Accessible
```bash
# Check if UI containers are running
docker-compose ps pgadmin mongo-express kafka-ui

# Restart UI services
docker-compose restart pgadmin mongo-express kafka-ui

# Check logs
docker-compose logs pgadmin
docker-compose logs mongo-express
docker-compose logs kafka-ui
```

### Connection Issues
- Ensure all services are running: `docker-compose up -d`
- Check network connectivity between containers
- Verify port conflicts (none should overlap)

### pgAdmin Setup Issues
```bash
# Reset pgAdmin
docker-compose down pgadmin
docker-compose up -d pgadmin
```

---

## 📊 Monitoring Dashboard Example

### Kafka UI Dashboard Shows:
- **Active Topics:** 3 (user.created, order.created, payment.completed)
- **Consumer Groups:** 2 active groups
- **Message Throughput:** Real-time message rates
- **Broker Health:** Kafka cluster status

### MongoDB Express Dashboard Shows:
- **Database Size:** Storage usage
- **Collection Stats:** Document counts
- **Index Performance:** Query optimization
- **Recent Operations:** Activity log

### pgAdmin Dashboard Shows:
- **Database Connections:** Active sessions
- **Query Performance:** Slow queries
- **Storage Usage:** Table sizes
- **Backup Status:** Maintenance info

---

## 🔐 Security Notes

⚠️ **Development Environment Only**
- Default credentials are for development
- No SSL/TLS encryption
- Open network access
- Consider security hardening for production

**Production Recommendations:**
- Change default passwords
- Enable SSL/TLS
- Use network isolation
- Implement proper authentication
- Set up audit logging

---

## 📱 Mobile Access

All UI tools are responsive and work on mobile devices:
- **pgAdmin:** Full mobile support
- **MongoDB Express:** Optimized for tablets
- **Kafka UI:** Mobile-friendly dashboard

Access from any device on the same network using your machine's IP address instead of `localhost`.
