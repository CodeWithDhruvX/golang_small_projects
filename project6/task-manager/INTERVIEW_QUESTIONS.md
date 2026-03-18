# Security & Go API Interview Questions

This document contains frequently asked interview questions covering the security features and Go concepts implemented in our secure API.

## 📋 Table of Contents

1. [Go Fundamentals](#go-fundamentals)
2. [API Security](#api-security)
3. [Encryption & Cryptography](#encryption--cryptography)
4. [Authentication & Authorization](#authentication--authorization)
5. [Error Handling & Validation](#error-handling--validation)
6. [Middleware & HTTP](#middleware--http)
7. [SSL/TLS & HTTPS](#ssltls--https)
8. [Logging & Monitoring](#logging--monitoring)
9. [Design Patterns & Architecture](#design-patterns--architecture)
10. [System Design & Best Practices](#system-design--best-practices)

---

## Go Fundamentals

### Q1: What are Go interfaces and how do they enable polymorphism?
**Answer:** Go interfaces are implicit contracts that define a set of method signatures. They enable polymorphism by allowing different types to satisfy the same interface, enabling code to work with different types through a common interface.

```go
type Entity interface {
    GetID() interface{}
    Validate() error
}

type Task struct { /* ... */ }
type User struct { /* ... */ }

// Both Task and User can be used as Entity
func ProcessEntity(e Entity) {
    // Works with any type that implements Entity
}
```

### Q2: Explain Go's concurrency model and how goroutines differ from threads.
**Answer:** Go uses goroutines (lightweight threads managed by Go runtime) and channels for communication. Goroutines are cheaper than OS threads (start with 2KB stack vs 1MB+), are multiplexed onto OS threads, and communicate via channels rather than shared memory.

### Q3: What is the purpose of `sync.RWMutex` and when would you use it?
**Answer:** `sync.RWMutex` provides a read-write lock allowing multiple readers or one exclusive writer. Use it when data is read frequently but written rarely, like our in-memory data stores.

```go
var tasksMu sync.RWMutex

// Multiple readers can access simultaneously
tasksMu.RLock()
defer tasksMu.RUnlock()

// Exclusive access for writes
tasksMu.Lock()
defer tasksMu.Unlock()
```

### Q4: How does Go's garbage collection work?
**Answer:** Go uses a concurrent tri-color mark-and-sweep GC. It runs concurrently with the program, has low pause times, and automatically manages memory. Developers don't need to manually allocate/deallocate memory.

---

## API Security

### Q5: What are the most important security headers for REST APIs?
**Answer:** Key security headers include:
- `X-Content-Type-Options: nosniff` - Prevents MIME-type sniffing
- `X-Frame-Options: DENY` - Prevents clickjacking
- `X-XSS-Protection: 1; mode=block` - Enables XSS protection
- `Strict-Transport-Security` - Enforces HTTPS
- `Content-Security-Policy` - Controls resource loading

### Q6: What is rate limiting and why is it important?
**Answer:** Rate limiting controls the number of requests a client can make in a time period. It prevents DoS attacks, ensures fair resource usage, and protects against brute force attacks.

```go
limiter := rate.NewLimiter(rate.Limit(60), 10) // 60 req/min, burst 10
if !limiter.Allow() {
    return http.StatusTooManyRequests
}
```

### Q7: Explain the principle of input validation in APIs.
**Answer:** Input validation ensures all incoming data meets expected format, type, and constraints before processing. It prevents injection attacks, data corruption, and ensures system stability.

### Q8: What is CORS and how do you configure it securely?
**Answer:** CORS (Cross-Origin Resource Sharing) controls cross-origin requests. Secure configuration includes:
- Whitelisting specific origins instead of using "*"
- Limiting allowed methods and headers
- Enabling credentials only when necessary

---

## Encryption & Cryptography

### Q9: What is AES-256-GCM and why is it preferred over other encryption modes?
**Answer:** AES-256-GCM is an authenticated encryption mode that provides both confidentiality and integrity. It's preferred because:
- 256-bit key size provides strong security
- GCM mode includes authentication (prevents tampering)
- Efficient hardware acceleration on modern CPUs
- No padding oracle attacks

### Q10: How do you securely store passwords in a database?
**Answer:** Use bcrypt, scrypt, or Argon2 with:
- High work factor to slow down brute force attacks
- Unique salt per password
- Never store plain text or reversible encryption

```go
// Use bcrypt for password hashing
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

### Q11: What is the difference between encryption at rest and encryption in transit?
**Answer:** 
- **Encryption at rest**: Protects data stored on disk/database
- **Encryption in transit**: Protects data moving over networks (HTTPS/TLS)
- Both are essential for complete data protection

### Q12: How do you manage encryption keys securely?
**Answer:** Best practices include:
- Use hardware security modules (HSM) in production
- Implement key rotation policies
- Never store keys with encrypted data
- Use environment variables or secret management systems
- Implement key derivation functions (PBKDF2, scrypt)

---

## Authentication & Authorization

### Q13: What is JWT and how does it work?
**Answer:** JWT (JSON Web Token) is a compact, URL-safe token format containing claims. It consists of:
- Header: Algorithm and token type
- Payload: Claims (user data, permissions)
- Signature: Cryptographic signature for integrity

```go
type JWTClaims struct {
    UserID   string    `json:"user_id"`
    Username string    `json:"username"`
    Role     string    `json:"role"`
    Exp      time.Time `json:"exp"`
}
```

### Q14: What is the difference between authentication and authorization?
**Answer:**
- **Authentication**: Verifying who you are (login, credentials)
- **Authorization**: Verifying what you can do (permissions, roles)

### Q15: What are the advantages and disadvantages of JWT tokens?
**Advantages:**
- Stateless (no server-side session storage)
- Portable across services
- Contains user information
- Standardized format

**Disadvantages:**
- Cannot be easily revoked
- Larger than session tokens
- Requires secure key management
- Vulnerable to token theft

### Q16: How would you implement role-based access control (RBAC)?
**Answer:** Implement RBAC with:
- User roles (admin, user, guest)
- Permission matrix per role
- Middleware to check permissions
- Hierarchical role system

```go
func RequireRoleMiddleware(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := getUserRole(c)
        if !checkRole(userRole, requiredRole) {
            c.JSON(403, gin.H{"error": "insufficient permissions"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

---

## Error Handling & Validation

### Q17: What is structured error handling and why is it important?
**Answer:** Structured error handling provides consistent, machine-readable error responses with:
- Error codes for programmatic handling
- Human-readable messages
- Request correlation IDs
- Detailed validation errors

### Q18: How do you implement comprehensive input validation?
**Answer:** Use multiple layers:
1. **Struct tags** for basic validation
2. **Custom validators** for business logic
3. **Middleware** for request validation
4. **Database constraints** as final layer

### Q19: What is the difference between validation errors and business logic errors?
**Answer:**
- **Validation errors**: Invalid input format/type (HTTP 400)
- **Business logic errors**: Valid input but violates rules (HTTP 409, 422)

### Q20: How do you handle panics gracefully in a web server?
**Answer:** Use recovery middleware to:
- Catch panics in handlers
- Log detailed error information
- Return structured error responses
- Maintain server stability

```go
func RecoveryMiddleware() gin.HandlerFunc {
    return gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, recovered interface{}) {
        log.Printf("Panic recovered: %v", recovered)
        c.JSON(500, gin.H{"error": "internal server error"})
    })
}
```

---

## Middleware & HTTP

### Q21: What is middleware and how does it work in Go web frameworks?
**Answer:** Middleware is code that runs between request receipt and handler execution. In Gin, it forms a chain where each middleware can process, modify, or stop the request.

### Q22: What is the order of middleware execution and why is it important?
**Answer:** Middleware executes in the order they're registered. Important order:
1. Recovery (first to catch panics)
2. Logging (to capture all requests)
3. Security headers
4. Authentication/Authorization
5. Rate limiting
6. Validation

### Q23: How do you implement request timeout in Go HTTP servers?
**Answer:** Use `context.WithTimeout` to set deadlines:
```go
ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
defer cancel()
r = r.WithContext(ctx)
```

### Q24: What is CORS and how do you handle it in Go?
**Answer:** CORS allows cross-origin requests with proper headers. Handle by:
- Adding `Access-Control-Allow-Origin` header
- Handling preflight OPTIONS requests
- Configuring allowed methods and headers

---

## SSL/TLS & HTTPS

### Q25: What is the difference between SSL and TLS?
**Answer:** TLS is the successor to SSL. TLS 1.0 was based on SSL 3.0, but TLS provides:
- Stronger security algorithms
- Better handshake protocol
- Protection against known SSL vulnerabilities

### Q26: How do you implement HTTPS in a Go web server?
**Answer:** Use `http.ListenAndServeTLS()` with certificate files:
```go
httpServer := &http.Server{
    Addr:      ":8443",
    TLSConfig: tlsConfig,
}
httpServer.ListenAndServeTLS("cert.pem", "key.pem")
```

### Q27: What are the key components of TLS configuration?
**Answer:** Important TLS settings:
- `MinVersion`: TLS 1.2 minimum
- `CipherSuites`: Only secure ciphers
- `PreferServerCipherSuites`: true
- Certificate validation

### Q28: How do you generate self-signed certificates for development?
**Answer:** Use Go's crypto packages to generate:
```go
// Generate private key
privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

// Create certificate template
template := x509.Certificate{
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name{...},
    NotBefore: time.Now(),
    NotAfter: time.Now().Add(365 * 24 * time.Hour),
}

// Create certificate
certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
```

---

## Logging & Monitoring

### Q29: What is structured logging and why is it better than traditional logging?
**Answer:** Structured logging uses consistent formats (like JSON) with:
- Machine-parseable output
- Consistent field names
- Easy filtering and analysis
- Better integration with log management systems

### Q30: What security events should you log in an API?
**Answer:** Critical events to log:
- Failed authentication attempts
- Unauthorized access attempts
- Rate limit violations
- Administrative actions
- Certificate expirations
- System errors and panics

### Q31: How do you implement request correlation in distributed systems?
**Answer:** Use request IDs that flow through the system:
- Generate unique ID per request
- Pass in headers between services
- Include in all log entries
- Return in responses for debugging

### Q32: What is the difference between INFO, WARN, and ERROR log levels?
**Answer:**
- **INFO**: Normal operational information
- **WARN**: Unexpected but non-critical situations
- **ERROR**: Error conditions that need attention
- **DEBUG**: Detailed diagnostic information

---

## Design Patterns & Architecture

### Q33: What is the Repository pattern and how does it apply to our API?
**Answer:** Repository pattern abstracts data access logic, providing a clean interface between business logic and data storage. In our API, we implemented generic `Repository[T Entity]` interface.

### Q34: What is dependency injection and why is it useful?
**Answer:** Dependency injection provides dependencies to objects rather than having them create dependencies. Benefits:
- Easier testing (mock dependencies)
- Loose coupling between components
- Better code organization
- Easier configuration changes

### Q35: What is the Service Layer pattern?
**Answer:** Service layer encapsulates business logic, providing a clean separation between:
- Controllers (HTTP handling)
- Services (business logic)
- Repositories (data access)

### Q36: How do you implement the Factory pattern in Go?
**Answer:** Use factory functions to create objects:
```go
func NewEncryptionService(key string) *EncryptionService {
    return &EncryptionService{key: deriveKey(key)}
}
```

---

## System Design & Best Practices

### Q37: What are the key principles of secure API design?
**Answer:** Key principles:
- **Defense in depth**: Multiple security layers
- **Least privilege**: Minimal necessary permissions
- **Fail securely**: Default to secure behavior
- **Input validation**: Never trust client input
- **Audit logging**: Track security events

### Q38: How do you design for scalability in a REST API?
**Answer:** Design considerations:
- Stateless services
- Horizontal scaling capability
- Database connection pooling
- Caching strategies
- Load balancing
- Rate limiting

### Q39: What is the difference between stateful and stateless APIs?
**Answer:**
- **Stateful**: Server maintains session state between requests
- **Stateless**: Each request contains all necessary information
- Stateless is preferred for scalability and reliability

### Q40: How do you implement graceful shutdown in Go web servers?
**Answer:** Handle OS signals and close connections:
```go
func gracefulShutdown(server *http.Server) {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    server.Shutdown(ctx)
}
```

---

## Practical Coding Questions

### Q41: Implement a rate limiter using Go's time/rate package
```go
func RateLimitMiddleware(limiter *rate.Limiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "rate limit exceeded"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### Q42: Write a middleware that adds security headers
```go
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Next()
    }
}
```

### Q43: Implement JWT token validation
```go
func ValidateJWT(tokenString, secret string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, err
}
```

### Q44: Write a function to encrypt data using AES-256-GCM
```go
func Encrypt(plaintext, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}
```

---

## Behavioral Questions

### Q45: How do you stay updated with security best practices?
**Answer:** 
- Follow security blogs and newsletters
- Participate in security communities
- Attend security conferences
- Regular security training and certifications
- Review OWASP guidelines
- Follow CVE announcements

### Q46: Describe a time you found and fixed a security vulnerability.
**Answer:** Focus on:
- How you discovered the vulnerability
- The impact assessment
- The fix implementation
- Testing and validation
- Lessons learned

### Q47: How do you approach security code reviews?
**Answer:** 
- Check for input validation
- Verify authentication/authorization
- Review error handling
- Check for sensitive data exposure
- Verify encryption usage
- Review logging security

### Q48: How would you explain a complex security concept to a non-technical stakeholder?
**Answer:** 
- Use analogies and simple terms
- Focus on business impact
- Provide concrete examples
- Avoid jargon
- Emphasize risk and mitigation

---

## Advanced Topics

### Q49: What is OAuth 2.0 and how does it differ from JWT?
**Answer:** OAuth 2.0 is an authorization framework, JWT is a token format. OAuth defines the flow for obtaining access tokens, while JWT defines the token structure.

### Q50: What are the OWASP Top 10 and how do they apply to APIs?
**Answer:** OWASP Top 10 includes:
- Broken Access Control
- Cryptographic Failures
- Injection
- Insecure Design
- Security Misconfiguration
- Vulnerable Components
- Authentication Failures
- Software/Data Integrity Failures
- Logging/Monitoring Failures
- Server-Side Request Forgery

### Q51: What is a zero-trust architecture?
**Answer:** Zero-trust assumes no implicit trust and verifies every request regardless of source. Key principles:
- Never trust, always verify
- Least privilege access
- Micro-segmentation
- Continuous monitoring

### Q52: How do you implement API versioning?
**Answer:** Common approaches:
- URL path versioning (`/api/v1/`)
- Header versioning (`Accept: application/vnd.api.v1+json`)
- Query parameter versioning (`?version=1`)
- Subdomain versioning (`v1.api.example.com`)

---

## Go-Specific Security Questions

### Q53: What are some Go-specific security considerations?
**Answer:**
- Goroutine safety and race conditions
- Memory safety with slices and maps
- Proper error handling to avoid information leakage
- Safe use of `unsafe` package
- Regular expression DoS protection

### Q54: How does Go handle race conditions?
**Answer:** Go provides:
- `go test -race` for race detection
- `sync` package for synchronization
- Channel-based communication
- Atomic operations package

### Q55: What is the difference between `make` and `new` in Go?
**Answer:**
- `new(T)`: Returns pointer to zeroed T
- `make(T)`: Returns initialized T (for slices, maps, channels)
- Use `make` for reference types, `new` rarely needed

---

## Testing Questions

### Q56: How do you test security features in Go?
**Answer:** 
- Unit tests for validation logic
- Integration tests for authentication
- Security tests for encryption
- Load tests for rate limiting
- Penetration testing

### Q57: What are table-driven tests in Go?
**Answer:** Table-driven tests use slices of test cases:
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        email string
        valid bool
    }{
        {"test@example.com", true},
        {"invalid-email", false},
    }
    
    for _, tt := range tests {
        result := ValidateEmail(tt.email)
        if result != tt.valid {
            t.Errorf("ValidateEmail(%s) = %v; want %v", tt.email, result, tt.valid)
        }
    }
}
```

---

## Performance Questions

### Q58: How do you optimize API performance while maintaining security?
**Answer:**
- Use connection pooling
- Implement caching with proper invalidation
- Optimize database queries
- Use efficient encryption algorithms
- Implement proper rate limiting
- Monitor and profile performance

### Q59: What is the impact of security middleware on performance?
**Answer:** Security features add overhead:
- Encryption/decryption CPU cost
- Validation processing time
- Logging I/O overhead
- TLS handshake latency
- Rate limiting checks

Mitigation strategies:
- Use hardware acceleration
- Optimize middleware order
- Cache validation results
- Async logging

---

## Database Security

### Q60: How do you prevent SQL injection in Go?
**Answer:** 
- Use parameterized queries with database/sql
- Use ORM packages that handle escaping
- Validate all inputs
- Use prepared statements
- Implement least privilege database users

```go
// Safe parameterized query
rows, err := db.Query("SELECT * FROM users WHERE id = ?", userID)
```

---

This comprehensive list covers the most frequently asked interview questions for the security topics and Go concepts implemented in our secure API. Practice answering these questions with concrete examples from our codebase!
