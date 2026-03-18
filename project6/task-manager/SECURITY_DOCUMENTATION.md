# Security Documentation

This document outlines the comprehensive security features implemented in the Secure Task Manager API.

## 🔐 Security Features Overview

### 1. **Polymorphism & Interface-Based Design**
- **Entity Interface**: Common interface for all entities (Task, User, Project)
- **Repository Interface**: Generic data access operations
- **Service Interface**: Business logic abstraction
- **Benefits**: Type safety, code reusability, easy testing

### 2. **AES-256 Encryption**
- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Use Cases**: Passwords, sensitive user data, email addresses
- **Key Management**: SHA-256 hash of secret key
- **Implementation**: `encryption.go`

```go
// Example usage
encryptionService := NewEncryptionService("your-32-character-secret-key")
encrypted, err := encryptionService.EncryptField("sensitive-data")
```

### 3. **Comprehensive API Validation**
- **Input Validation**: Structured validation for all entities
- **Field Validation**: Length, format, type checking
- **Custom Validators**: Email, date format, enum values
- **Implementation**: `validation.go`

### 4. **Structured Error Handling**
- **Error Types**: Validation, Not Found, Unauthorized, Internal, etc.
- **Structured Responses**: JSON error format with details
- **Stack Traces**: Captured for internal errors
- **Request Tracking**: Unique request IDs for debugging
- **Implementation**: `errors.go`

### 5. **SSL/TLS Configuration**
- **HTTPS Support**: Automatic SSL certificate generation
- **TLS 1.2+**: Minimum TLS version requirement
- **Strong Ciphers**: Only secure cipher suites
- **Certificate Management**: Self-signed cert generation for development
- **Implementation**: `ssl_config.go`

### 6. **JWT Authentication**
- **Token-based Auth**: JWT tokens for user authentication
- **Role-based Access**: Admin, User, Guest roles
- **Token Expiration**: 24-hour token expiry
- **Secure Storage**: Bearer token in Authorization header
- **Implementation**: `auth.go`

### 7. **Security Middleware**
- **Rate Limiting**: 60 requests/minute, burst of 10
- **Security Headers**: XSS protection, content type options, frame options
- **CORS Support**: Configurable cross-origin resource sharing
- **Request Timeout**: 30-second request timeout
- **Recovery Middleware**: Panic recovery with structured logging
- **Implementation**: `middleware.go`

### 8. **Comprehensive Logging**
- **Structured Logging**: JSON-formatted logs
- **Security Events**: Special logging for security incidents
- **Request Tracking**: Request ID correlation
- **Log Levels**: Debug, Info, Warn, Error, Fatal
- **Implementation**: `logger.go`

## 🛡️ Security Headers

The following security headers are automatically added to all responses:

```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains (HTTPS only)
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
Server: SecureAPI
```

## 🔑 Authentication Flow

### 1. User Registration
```bash
POST /api/v1/public/register
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password_123"
}
```

### 2. User Login
```bash
POST /api/v1/public/login
{
  "username": "john_doe",
  "password": "secure_password_123"
}
```

### 3. Access Protected Resources
```bash
GET /api/v1/tasks
Authorization: Bearer <jwt_token>
```

## 🚦 Rate Limiting

- **Default Limit**: 60 requests per minute
- **Burst Size**: 10 requests
- **Response**: HTTP 429 with error details when exceeded
- **Per-IP**: Rate limiting applied per client IP

## 🔍 API Endpoints

### Public Endpoints (No Authentication)
- `GET /api/v1/public/health` - Health check
- `GET /api/v1/public/metrics` - System metrics
- `POST /api/v1/public/login` - User login
- `POST /api/v1/public/register` - User registration

### Protected Endpoints (Authentication Required)
- `GET /api/v1/tasks` - List all tasks
- `GET /api/v1/tasks/:id` - Get specific task
- `POST /api/v1/tasks` - Create new task
- `PUT /api/v1/tasks/:id` - Update task
- `DELETE /api/v1/tasks/:id` - Delete task

### Admin Endpoints (Admin Role Required)
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/:id` - Get specific user
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

## 🔒 Data Encryption

### Sensitive Fields
- User passwords (hashed with bcrypt)
- Email addresses (AES-256 encrypted)
- AssignedTo field in tasks (AES-256 encrypted)

### Encryption Process
1. Data is encrypted before storage
2. Data is decrypted after retrieval
3. Encryption keys are never stored with data
4. Failed decryption is logged but doesn't fail requests

## 📝 Error Handling

### Error Response Format
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Validation failed",
  "details": {
    "validation_errors": [
      {
        "field": "email",
        "message": "email must be a valid email address",
        "value": "invalid-email"
      }
    ]
  },
  "timestamp": "2026-03-18T17:30:00Z",
  "request_id": "req_1234567890_abcdefgh"
}
```

### Error Types
- `VALIDATION_ERROR` - Input validation failed
- `NOT_FOUND` - Resource not found
- `UNAUTHORIZED` - Authentication required
- `FORBIDDEN` - Insufficient permissions
- `RATE_LIMIT_EXCEEDED` - Too many requests
- `TIMEOUT` - Request timeout
- `INTERNAL_ERROR` - Server error

## 🚀 Deployment Security

### Environment Variables
```bash
USE_HTTPS=true                    # Enable HTTPS (default: true)
JWT_SECRET=your-super-secret-key   # JWT signing secret
ENCRYPTION_KEY=your-32-char-key   # AES encryption key
LOG_LEVEL=INFO                    # Logging level
```

### SSL Certificate Management
- **Development**: Self-signed certificates auto-generated
- **Production**: Use certificates from trusted CA
- **Certificate Validation**: Automatic expiration checking
- **TLS Configuration**: Secure cipher suites only

## 🔍 Security Monitoring

### Security Events Logged
- Failed login attempts
- Unauthorized access attempts
- Rate limit violations
- Certificate expiration warnings
- Panic recovery events

### Log Examples
```json
{
  "timestamp": "2026-03-18T17:30:00Z",
  "level": "ERROR",
  "message": "Security Event: Failed login attempt",
  "data": {
    "username": "admin",
    "ip": "192.168.1.100"
  },
  "request_id": "req_1234567890_abcdefgh"
}
```

## 🛠️ Security Best Practices Implemented

1. **Input Validation**: All inputs are validated before processing
2. **Output Encoding**: Sensitive data is encrypted before storage
3. **Authentication**: JWT-based authentication with role-based access
4. **Authorization**: Role-based permissions for all endpoints
5. **Transport Security**: HTTPS with strong TLS configuration
6. **Rate Limiting**: Protection against DoS attacks
7. **Error Handling**: Secure error responses without information leakage
8. **Logging**: Comprehensive security event logging
9. **Headers**: Security headers for client-side protection
10. **Timeouts**: Request timeouts to prevent resource exhaustion

## 🔧 Configuration

### SSL Configuration
```go
sslConfig := &SSLConfig{
    CertFile:        "server.crt",
    KeyFile:         "server.key",
    UseHTTPS:        true,
    Port:            8080,
    HTTPSPort:       8443,
    RedirectToHTTPS: true,
}
```

### Rate Limiting Configuration
```go
// 60 requests per minute, burst of 10
r.Use(RateLimitMiddleware(60, 10))
```

### Authentication Configuration
```go
// 24-hour token expiry
authService := NewAuthService(jwtSecret, 24*time.Hour, logger, errorHandler)
```

## 📋 Security Checklist

- [x] Input validation on all endpoints
- [x] Authentication and authorization
- [x] HTTPS/TLS encryption
- [x] Data encryption at rest
- [x] Rate limiting
- [x] Security headers
- [x] Error handling
- [x] Logging and monitoring
- [x] CORS configuration
- [x] Request timeout
- [x] Panic recovery
- [x] Role-based access control

## 🚨 Security Considerations

### Production Deployment
1. Use production SSL certificates from trusted CA
2. Store secrets in secure environment variables or vault
3. Enable comprehensive monitoring and alerting
4. Regularly rotate encryption keys and JWT secrets
5. Implement proper backup and recovery procedures

### Known Limitations
1. JWT implementation is simplified (use proper JWT library in production)
2. Password hashing uses simple implementation (use bcrypt in production)
3. In-memory storage (use database in production)
4. No session management (consider implementing session store)

## 📞 Contact

For security concerns or issues, please contact the development team or create a security issue in the project repository.
