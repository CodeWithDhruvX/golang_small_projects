# Interview Questions - Project7 Multi-Database Application

This document contains frequently asked interview questions based on the concepts and technologies used in this project.

## 🗄️ Database & ORM Questions

### Q1: What is GORM and why would you use it over raw SQL?
**Answer**: GORM is Go's ORM library that provides a developer-friendly way to interact with databases. 

**Advantages:**
- Type-safe operations with Go structs
- Automatic migrations
- Relationships and associations
- Built-in validation and hooks
- Query builder for dynamic queries
- Database agnostic (switch between PostgreSQL, MySQL, etc.)

**When to use raw SQL instead:**
- Complex queries with multiple JOINs
- Performance-critical operations
- Database-specific features not supported by ORM
- Aggregation queries with complex grouping

### Q2: How do you handle database migrations in GORM?
**Answer**: 
```go
// Auto migration
err = db.AutoMigrate(
    &models.User{},
    &models.Post{},
    &models.Tag{},
)
```

**Best practices:**
- Use AutoMigrate in development
- Create separate migration files for production
- Version your migrations
- Test migrations on staging first
- Handle rollback scenarios

### Q3: Explain the difference between GORM's `Raw()` and `Exec()` methods.
**Answer**: 
- **`db.Raw()`**: Executes raw SQL and scans results into structs, used for SELECT queries
- **`db.Exec()`**: Executes raw SQL without returning results, used for INSERT/UPDATE/DELETE

**Example:**
```go
// Raw - for SELECT
db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)

// Exec - for INSERT/UPDATE/DELETE
result := db.Exec("UPDATE users SET active = ? WHERE id = ?", false, userID)
```

## 🐳 Docker & Containerization Questions

### Q4: Why use Docker Compose for database setup?
**Answer**: Docker Compose provides:
- **Consistency**: Same environment across development, testing, production
- **Isolation**: Database doesn't affect local system
- **Version Control**: Database versions defined in code
- **Easy Setup**: One command to start entire stack
- **Portability**: Works on any machine with Docker

### Q5: Explain the Docker Compose file in this project.
**Answer**: 
```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: project7_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
```

**Key concepts:**
- **Services**: Each container (postgres, mongodb, adminer)
- **Volumes**: Persistent data storage
- **Networks**: Inter-container communication
- **Environment variables**: Configuration
- **Port mapping**: Host to container port access

## 🚀 Architecture & Design Questions

### Q6: Why use both PostgreSQL and MongoDB in the same application?
**Answer**: 
- **PostgreSQL (Relational)**: Structured data with relationships
  - Users, Posts, Categories (ACID compliance)
  - Complex transactions and constraints
  
- **MongoDB (NoSQL)**: Flexible schema for unstructured data
  - Logs, Analytics, User Profiles (horizontal scaling)
  - Fast writes and schema flexibility

**This is called Polyglot Persistence**

### Q7: Explain the project structure and why it's organized this way.
**Answer**: 
```
project7/
├── cmd/           # Application entry points
├── internal/      # Private application code
│   ├── config/    # Configuration management
│   ├── database/  # Database connections
│   ├── handlers/  # HTTP request handlers
│   └── models/    # Data models
├── docker/        # Docker configurations
└── docker-compose.yml
```

**Benefits:**
- **Separation of concerns**: Each package has single responsibility
- **Testability**: Easy to unit test individual components
- **Maintainability**: Clear organization for future development
- **Scalability**: Easy to add new features

## 🔌 API & REST Questions

### Q8: What are RESTful principles and how does this API follow them?
**Answer**: 
**REST Principles:**
1. **Stateless**: Each request contains all information
2. **Client-Server**: Clear separation of concerns
3. **Uniform Interface**: Standard HTTP methods (GET, POST, PUT, DELETE)
4. **Resource-based**: URLs identify resources (e.g., `/api/v1/users`)

**Implementation:**
```go
// GET /api/v1/users - Get all users
// GET /api/v1/users/:id - Get specific user
// POST /api/v1/users - Create user
// PUT /api/v1/users/:id - Update user
// DELETE /api/v1/users/:id - Delete user
```

### Q9: How do you handle error responses in a REST API?
**Answer**: 
```go
if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}

if err := h.PostgresDB.CreateUser(&user); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}
```

**Best practices:**
- Use appropriate HTTP status codes
- Return consistent error format
- Include descriptive error messages
- Log errors for debugging

## 🔍 Raw SQL & Performance Questions

### Q10: When would you use raw SQL instead of GORM?
**Answer**: 
**Use Raw SQL for:**
- Complex JOIN operations
- Database-specific functions
- Performance-critical queries
- Aggregation with complex grouping
- Window functions

**Example from project:**
```go
query := `
    SELECT u.id, u.username, u.email, COUNT(p.id) as post_count 
    FROM users u 
    LEFT JOIN posts p ON u.id = p.user_id 
    WHERE u.deleted_at IS NULL 
    GROUP BY u.id, u.username, u.email 
    ORDER BY post_count DESC`
```

### Q11: How do you prevent SQL injection in raw queries?
**Answer**: 
**Always use parameterized queries:**
```go
// ❌ VULNERABLE
query := fmt.Sprintf("SELECT * FROM users WHERE age = %d", age)

// ✅ SAFE
query := "SELECT * FROM users WHERE age = ?"
err := db.Raw(query, age).Scan(&users).Error
```

**GORM automatically handles parameter binding**

## 🔄 Concurrency & Performance Questions

### Q12: How would you handle database connection pooling?
**Answer**: 
```go
// GORM automatically handles connection pooling
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})

// Configure pool settings
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Q13: How do you optimize database performance?
**Answer**: 
**Indexing:**
```go
// MongoDB indexes
_, err = logCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{{Key: "timestamp", Value: -1}},
})

// PostgreSQL indexes (handled by GORM)
type User struct {
    Email string `gorm:"uniqueIndex"`
    Age   int    `gorm:"index"`
}
```

**Query Optimization:**
- Use appropriate indexes
- Avoid N+1 queries with preloading
- Use raw SQL for complex operations
- Implement pagination
- Cache frequently accessed data

## 🛡️ Security Questions

### Q14: How do you secure database connections?
**Answer**: 
**Environment Variables:**
```go
func LoadConfig() *AppConfig {
    return &AppConfig{
        Database: DatabaseConfig{
            PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
            MongoPassword:    getEnv("MONGO_PASSWORD", ""),
        },
    }
}
```

**Best practices:**
- Never hardcode credentials
- Use environment variables
- Implement connection encryption
- Use least privilege principle
- Regular password rotation

### Q15: How do you handle CORS in a Go API?
**Answer**: 
```go
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}
```

## 🧪 Testing Questions

### Q16: How would you test this application?
**Answer**: 
**Unit Tests:**
```go
func TestCreateUser(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    
    // Test user creation
    user := &models.User{Username: "test", Email: "test@example.com"}
    err := db.Create(user)
    
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

**Integration Tests:**
- Use testcontainers for real databases
- Test API endpoints with test server
- Mock external dependencies

## 📊 Monitoring & Observability Questions

### Q17: How do you monitor application health?
**Answer**: 
**Health Check Endpoint:**
```go
func (h *Handler) HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":    "healthy",
        "timestamp": time.Now(),
        "databases": gin.H{
            "postgres": "connected",
            "mongodb":  "connected",
        },
    })
}
```

**Logging Strategy:**
- Structured logging with context
- Different log levels (info, warn, error)
- Centralized log collection
- Performance metrics

## 🚀 Deployment Questions

### Q18: How would you deploy this application to production?
**Answer**: 
**Deployment Strategy:**
1. **Containerize**: Docker image with multi-stage build
2. **Database**: Managed database services (RDS, DocumentDB)
3. **Load Balancer**: nginx or cloud load balancer
4. **Monitoring**: Prometheus + Grafana
5. **CI/CD**: GitHub Actions or Jenkins

**Dockerfile:**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## 🔧 Troubleshooting Questions

### Q19: How would you debug a slow database query?
**Answer**: 
**Steps:**
1. **Enable query logging**: `logger.Default.LogMode(logger.Info)`
2. **Use EXPLAIN ANALYZE**: Analyze query execution plan
3. **Check indexes**: Ensure proper indexing
4. **Monitor connections**: Check for connection leaks
5. **Profile application**: Use pprof for CPU/memory analysis

### Q20: How do you handle database schema changes in production?
**Answer**: 
**Migration Strategy:**
1. **Version control migrations**: Use migration tools
2. **Backward compatibility**: Support old and new schema
3. **Rollback plan**: Always have rollback strategy
4. **Blue-green deployment**: Zero downtime deployment
5. **Test thoroughly**: Staging environment validation

## 💡 Advanced Questions

### Q21: Explain the difference between SQL and NoSQL databases.
**Answer**: 
**SQL (Relational):**
- Structured data with predefined schema
- ACID properties for data consistency
- Complex relationships and transactions
- Vertical scaling

**NoSQL (MongoDB):**
- Flexible schema, unstructured data
- BASE properties for availability
- Horizontal scaling
- Document-based storage

**When to use each:**
- SQL: Financial data, user accounts, inventory
- NoSQL: Analytics, logs, content management, real-time data

### Q22: How would you implement caching in this application?
**Answer**: 
**Caching Strategy:**
```go
// Redis cache example
func (h *Handler) GetUserWithCache(c *gin.Context) {
    cacheKey := fmt.Sprintf("user:%s", c.Param("id"))
    
    // Check cache
    if cached, err := redis.Get(cacheKey).Result(); err == nil {
        c.JSON(http.StatusOK, cached)
        return
    }
    
    // Get from database
    user, err := h.PostgresDB.GetUser(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    // Set cache
    redis.Set(cacheKey, user, 5*time.Minute)
    c.JSON(http.StatusOK, user)
}
```

**Cache Patterns:**
- Cache-aside: Most common pattern
- Write-through: Update cache on write
- Write-behind: Async cache updates
- Cache invalidation: TTL and manual invalidation

---

## 🎯 Quick Tips for Interviews

1. **Explain your design decisions** - Why you chose certain technologies
2. **Discuss trade-offs** - Performance vs. maintainability
3. **Show practical examples** - Reference your actual code
4. **Mention best practices** - Security, testing, monitoring
5. **Talk about scalability** - How would you handle growth
6. **Be honest about limitations** - Nobody writes perfect code

Remember: Interviewers want to see your thought process and problem-solving approach, not just the "right" answer!
