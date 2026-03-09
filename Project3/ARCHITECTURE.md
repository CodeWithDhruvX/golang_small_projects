# Private Knowledge Base - Architecture Documentation

## 🏗️ **System Architecture Overview**

The Private Knowledge Base is a **cloud-free, self-hosted** AI-powered document intelligence system built with modern microservices architecture.

```
┌─────────────────────────────────────────────────────────────────┐
│                    Private Knowledge Base                        │
├─────────────────────────────────────────────────────────────────┤
│  Frontend Layer (Angular 18)                                    │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │   Chat UI   │ │ Documents  │ │   Upload    │ │   Login     │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│  API Gateway Layer (Go/Gin)                                     │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │   Auth      │ │   Documents │ │    Chat     │ │  Ingestion  │ │
│  │ Middleware  │ │   Handler   │ │   Handler   │ │   Service   │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│  Business Logic Layer                                           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │   RAG       │ │  Embeddings │ │   Search    │ │  Chunking   │ │
│  │   Service   │ │   Service   │ │   Service   │ │   Service   │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│  Data Layer                                                     │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │ PostgreSQL  │ │   PGVector  │ │   Redis     │ │ File Storage│ │
│  │   + Tables │ │  + Vectors  │ │   Cache     │ │   System    │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│  AI/ML Layer                                                    │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │   Ollama    │ │  Llama 3.1  │ │ Nomic Embed │ │  Model Mgmt │ │
│  │   Service   │ │    LLM      │ │   Text      │ │   Service   │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│  Infrastructure Layer                                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │ Prometheus  │ │   Grafana   │ │   Docker    │ │ Kubernetes  │ │
│  │  Metrics    │ │ Dashboard   │ │  Containers │ │   Cluster   │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## 🏛️ **Architectural Principles**

### **1. Microservices Architecture**
- **Service Separation**: Each component is an independent service
- **API First**: All services communicate through well-defined APIs
- **Containerized**: Services run in isolated Docker containers
- **Scalable**: Horizontal scaling with Kubernetes HPA

### **2. Cloud-Native Design**
- **Stateless Services**: Application logic is stateless
- **External State**: State managed in databases and caches
- **Health Checks**: All services have liveness/readiness probes
- **Graceful Degradation**: Services handle failures gracefully

### **3. Security First**
- **Zero Trust**: All communications require authentication
- **JWT Tokens**: Stateless authentication with expiration
- **Input Validation**: All inputs validated and sanitized
- **Least Privilege**: Services run with minimal permissions

### **4. Observability**
- **Metrics**: Prometheus collects all service metrics
- **Logging**: Structured logging with correlation IDs
- **Tracing**: Request tracing across services
- **Alerting**: Proactive alerting on critical issues

## 🏗️ **Component Architecture**

### **Frontend Architecture (Angular 18)**

```
src/app/
├── components/           # UI Components
│   ├── auth/            # Authentication components
│   │   └── login/       # Login form
│   ├── chat/            # Chat interface
│   ├── documents/       # Document management
│   ├── upload/          # File upload
│   └── core/            # Shared components
│       ├── header/      # Application header
│       ├── sidebar/     # Navigation sidebar
│       └── theme-toggle/ # Theme switcher
├── services/            # Business logic
│   ├── auth.service.ts  # Authentication
│   ├── chat.service.ts  # Chat functionality
│   └── document.service.ts # Document management
├── models/              # Data models
│   ├── auth.model.ts    # Auth types
│   ├── chat.model.ts    # Chat types
│   └── document.model.ts # Document types
├── guards/              # Route guards
│   └── auth.guard.ts    # Authentication guard
├── interceptors/        # HTTP interceptors
│   ├── auth.interceptor.ts  # Add auth headers
│   └── error.interceptor.ts # Error handling
└── app.routes.ts        # Application routing
```

**Key Features:**
- **Standalone Components**: Angular 18 standalone components
- **Signals**: Reactive state management
- **RxJS**: Reactive programming for async operations
- **Tailwind CSS**: Modern utility-first styling
- **TypeScript**: Full type safety

### **Backend Architecture (Go)**

```
cmd/
└── server/
    └── main.go           # Application entry point

internal/
├── auth/                # Authentication module
│   ├── auth.go          # JWT token management
│   ├── middleware.go    # Auth middleware
│   └── handlers.go      # Auth HTTP handlers
├── config/              # Configuration management
│   ├── config.go        # Configuration loading
│   └── database.go      # Database config
├── ingestion/           # Document processing
│   ├── service.go       # Ingestion orchestration
│   ├── markdown.go      # Markdown parser
│   ├── pdf.go           # PDF processor
│   ├── txt.go           # Text processor
│   └── go.go            # Go source processor
├── rag/                 # RAG (Retrieval-Augmented Generation)
│   ├── service.go       # RAG orchestration
│   ├── retrieval.go     # Vector search
│   └── generation.go    # LLM integration
├── storage/             # Data persistence
│   ├── postgres.go      # PostgreSQL operations
│   ├── models.go        # Data models
│   └── migrations.go    # Database migrations
├── web/                 # HTTP layer
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   └── routes.go        # Route definitions
└── metrics/             # Observability
    └── metrics.go       # Prometheus metrics

pkg/
├── logger/              # Logging utilities
│   └── logger.go
└── metrics/             # Metrics collection
    └── metrics.go
```

**Key Features:**
- **Clean Architecture**: Separation of concerns
- **Dependency Injection**: Wire for dependency management
- **Gin Framework**: High-performance HTTP framework
- **PGX**: PostgreSQL driver with connection pooling
- **Prometheus**: Metrics collection
- **Logrus**: Structured logging

## 🗄️ **Database Architecture**

### **PostgreSQL Schema**

```sql
-- Documents table
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    file_size INTEGER NOT NULL,
    upload_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed BOOLEAN DEFAULT FALSE,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Document chunks for RAG
CREATE TABLE document_chunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    embedding vector(768), -- For Nomic-embed-text
    page_number INTEGER,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Chat sessions
CREATE TABLE chat_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_name VARCHAR(255),
    model_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Chat messages
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    message_type VARCHAR(20) NOT NULL, -- 'user' or 'assistant'
    content TEXT NOT NULL,
    citations JSONB,
    model_used VARCHAR(100),
    tokens_used INTEGER,
    response_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Vector indexes for similarity search
CREATE INDEX ON document_chunks USING ivfflat (embedding vector_cosine_ops);
```

### **Redis Cache Architecture**

```
Redis Data Structures:
├── auth:tokens:{user_id}     # User session tokens
├── cache:embeddings:{hash}   # Document embedding cache
├── cache:search:{query_hash} # Search result cache
├── metrics:requests         # Request rate limiting
└── metrics:models:{model}    # Model usage statistics
```

## 🔐 **Security Architecture**

### **Authentication Flow**

```
1. User Login Request
   ↓
2. Validate Credentials
   ↓
3. Generate JWT Token
   ↓
4. Store Token in Redis
   ↓
5. Return Token to Client
   ↓
6. Client Includes Token in Requests
   ↓
7. Middleware Validates Token
   ↓
8. Process Request
```

### **Security Layers**

1. **Network Security**
   - Internal service network
   - TLS encryption for external traffic
   - Network policies in Kubernetes

2. **Application Security**
   - JWT authentication with expiration
   - Input validation and sanitization
   - SQL injection prevention
   - XSS protection

3. **Container Security**
   - Non-root user execution
   - Resource limits
   - Read-only filesystem where possible
   - Security scanning

4. **Data Security**
   - Encrypted data at rest (PostgreSQL)
   - Encrypted data in transit (TLS)
   - Access logging and auditing

## 🤖 **AI/ML Architecture**

### **Model Integration**

```
Ollama Service (localhost:11434)
├── Llama 3.1 8B (Chat/Generation)
│   ├── Input: User message + context
│   ├── Output: Generated response
│   └── Latency: ~3-5 seconds
└── Nomic Embed Text (Embeddings)
    ├── Input: Text chunks
    ├── Output: 768-dimension vectors
    └── Latency: ~500ms
```

### **RAG Pipeline**

```
1. User Query
   ↓
2. Query Embedding (Nomic-embed-text)
   ↓
3. Vector Search (PGVector)
   ↓
4. Context Retrieval (Top-K chunks)
   ↓
5. Prompt Engineering
   ↓
6. LLM Generation (Llama 3.1)
   ↓
7. Response with Citations
```

### **Document Processing Pipeline**

```
1. File Upload
   ↓
2. Content Type Detection
   ↓
3. Content Extraction
   ├── PDF → Text extraction
   ├── Markdown → AST parsing
   ├── Text → Direct processing
   └── Go → Source code parsing
   ↓
4. Text Chunking (20% overlap)
   ↓
5. Embedding Generation (Nomic-embed-text)
   ↓
6. Vector Storage (PGVector)
   ↓
7. Metadata Indexing
```

## 📊 **Monitoring Architecture**

### **Metrics Collection**

```
Prometheus Metrics:
├── Application Metrics
│   ├── http_requests_total
│   ├── http_request_duration_seconds
│   ├── documents_processed_total
│   └── chat_sessions_total
├── Database Metrics
│   ├── postgres_connections_active
│   ├── postgres_query_duration_seconds
│   └── pgvector_index_size_bytes
├── AI/ML Metrics
│   ├── ollama_requests_total
│   ├── ollama_request_duration_seconds
│   ├── embedding_generation_duration
│   └── model_tokens_used_total
└── System Metrics
    ├── cpu_usage_percent
    ├── memory_usage_bytes
    └── disk_usage_bytes
```

### **Grafana Dashboards**

1. **Application Overview**
   - Request rate and latency
   - Error rates
   - Active users
   - Document processing stats

2. **Database Performance**
   - Connection pool metrics
   - Query performance
   - Vector index efficiency
   - Storage usage

3. **AI/ML Metrics**
   - Model usage statistics
   - Generation latency
   - Token consumption
   - Embedding performance

4. **Infrastructure Health**
   - CPU and memory usage
   - Disk I/O
   - Network traffic
   - Container health

## 🚀 **Deployment Architecture**

### **Docker Compose (Development)**

```
Services:
├── postgres (Database)
├── pgadmin (Database UI)
├── redis (Cache)
├── ollama (AI Models)
├── prometheus (Metrics)
├── grafana (Dashboards)
└── go-backend (Application)
```

### **Kubernetes (Production)**

```
Namespace: knowledge-base
├── postgres-stateful.yaml
│   ├── PostgreSQL StatefulSet
│   ├── PGVector extension
│   └── Persistent Volumes
├── backend-deploy.yaml
│   ├── Go Backend Deployment
│   ├── Horizontal Pod Autoscaler
│   └── Services
├── frontend-deploy.yaml
│   ├── Angular Frontend Deployment
│   ├── Horizontal Pod Autoscaler
│   └── Services
├── ingress.yaml
│   ├── NGINX Ingress Controller
│   ├── SSL Termination
│   └── Routing Rules
└── monitoring/
    ├── prometheus.yaml
    └── grafana.yaml
```

## 🔄 **Data Flow Architecture**

### **Document Upload Flow**

```
Frontend → API Gateway → Ingestion Service → Storage Layer
    ↓           ↓              ↓                ↓
File Upload → Auth Check → Content Process → Database Store
    ↓           ↓              ↓                ↓
Progress   → JWT Valid → Chunk & Embed → Vector Index
```

### **Chat Query Flow**

```
Frontend → API Gateway → RAG Service → AI Service → Storage Layer
    ↓           ↓              ↓            ↓              ↓
User Query → Auth Check → Vector Search → LLM Call → Context Retrieval
    ↓           ↓              ↓            ↓              ↓
Streaming  → JWT Valid → Similarity    → Generation → Citations
```

## 🏛️ **Scalability Architecture**

### **Horizontal Scaling**

1. **Stateless Application Services**
   - Multiple instances behind load balancer
   - Session state in Redis
   - Database connection pooling

2. **Database Scaling**
   - Read replicas for read-heavy workloads
   - Connection pooling
   - Query optimization

3. **AI Model Scaling**
   - Multiple Ollama instances
   - Model caching
   - Request queuing

### **Performance Optimization**

1. **Caching Strategy**
   - Redis for session data
   - Embedding cache for repeated queries
   - Search result caching

2. **Database Optimization**
   - Vector indexes for similarity search
   - Query optimization
   - Connection pooling

3. **Frontend Optimization**
   - Lazy loading
   - Code splitting
   - Asset optimization

## 🛡️ **Reliability Architecture**

### **High Availability**

1. **Service Redundancy**
   - Multiple service instances
   - Health checks and auto-restart
   - Graceful shutdown

2. **Data Redundancy**
   - Database backups
   - Persistent volume replication
   - Multi-zone deployment

3. **Monitoring and Alerting**
   - Service health monitoring
   - Performance alerting
   - Automated recovery

### **Error Handling**

1. **Graceful Degradation**
   - Fallback mechanisms
   - Error boundaries
   - User-friendly error messages

2. **Retry Logic**
   - Exponential backoff
   - Circuit breakers
   - Request queuing

---

## 📋 **Architecture Decision Records (ADRs)**

### **ADR-001: Technology Stack Selection**
- **Decision**: Go + Angular + PostgreSQL + Ollama
- **Rationale**: Performance, ecosystem, cloud-free requirements
- **Status**: Implemented

### **ADR-002: Vector Database Choice**
- **Decision**: PostgreSQL + PGVector
- **Rationale**: Single database solution, ACID compliance, existing expertise
- **Status**: Implemented

### **ADR-003: Authentication Strategy**
- **Decision**: JWT tokens with Redis session storage
- **Rationale**: Stateless services, easy scaling, secure
- **Status**: Implemented

### **ADR-004: AI Model Integration**
- **Decision**: Ollama for local AI model serving
- **Rationale**: Cloud-free, privacy, cost-effective
- **Status**: Implemented

### **ADR-005: Monitoring Strategy**
- **Decision**: Prometheus + Grafana
- **Rationale**: Industry standard, scalable, rich ecosystem
- **Status**: Implemented

---

This architecture provides a **production-ready, scalable, and secure** foundation for the Private Knowledge Base system while maintaining the **cloud-free** requirement and ensuring **data privacy** and **control**.
