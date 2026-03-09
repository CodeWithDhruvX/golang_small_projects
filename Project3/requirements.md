Building a completely **in-house, cloud-free** AI Knowledge Base requires a professional-grade engineering approach. This project will transform your local hardware into a private "intelligence" hub using the models you already have installed.

Below are the **full requirements** for a **Go-based implementation** organized by functional, technical, and architectural categories.

---

## 📋 1. Functional Requirements

### **A. Intelligent Knowledge Ingestion**

* **Multi-Format Support:** Users must be able to upload PDF, TXT, Markdown, and `.go` files.
* **Automated Parsing:** The system should automatically extract text and structure (headers, paragraphs) using **Go libraries** (unidoc/unipdf, go-markdown, etc.).
* **Smart Chunking:** Documents must be split into manageable "chunks" (e.g., 500-1000 tokens) with a **20% overlap** to ensure the AI doesn't lose context between pieces.
* **Metadata Tagging:** Automatically tag chunks with source file name, page number (for PDFs), and upload timestamp.

### **B. Private Semantic Search (RAG)**

* **Vector Retrieval:** When a user asks a question, the system must use **Nomic-Embed-Text** to search the database for the most relevant document passages.
* **Hallucination Control:** The AI must be strictly instructed (via system prompts) to answer *only* based on the retrieved document context.
* **Citations:** The UI should display which document and page the AI used to generate its answer.
* **Database Monitoring:** Real-time monitoring of vector store performance, query latency, and storage utilization.

### **C. User Interface & Experience (Angular 18)**

* **Real-time Streaming:** The chat must feel responsive by using **Server-Sent Events (SSE)** to stream text as it's generated.
* **Expert Mode Toggle:** 
  * **Llama 3.1 8B:** Default for reasoning.
  * **Qwen2.5-Coder:** Specialized mode for code analysis.
  * **Phi-3:** "Fast Mode" for simple summaries.

* **Admin Dashboard:** A view to see all indexed files, manually trigger a re-index, or delete data from the vector store.

---

## 🛠️ 2. Technical Requirements

### **Backend (Go 1.21+)**

* **Web Framework:** **Gin** or **Echo** for high-performance HTTP routing
* **Concurrency:** **Goroutines** and **channels** for handling high-concurrency AI requests efficiently
* **Security:** Local **JWT-based authentication** using `github.com/golang-jwt/jwt` (since it's in-house, you can integrate with a local LDAP if needed)
* **Storage:** **PostgreSQL + PGVector** extension with `github.com/jackc/pgx/v5` driver
* **Vector Operations:** **pgvector** Go bindings for similarity search
* **Configuration:** **Viper** for configuration management
* **Logging:** **Logrus** or **Zap** for structured logging
* **Document Parsing:**
  * PDF: `github.com/unidoc/unipdf/v3`
  * Markdown: `github.com/gomarkdown/markdown`
  * Text: Standard library `bufio` and `io`
  * Go files: `go/parser` and `go/token`

### **Frontend (Angular 18)**

* **Standalone Components:** Modern, modular architecture
* **RxJS & Signals:** For reactive state management and handling real-time AI streams
* **Tailwind CSS:** For a clean, dark-themed "Developer Assistant" interface
* **TypeScript:** Full type safety with Angular's built-in support
* **Angular CLI:** Development tooling and build system

---

## ☁️ 3. Infrastructure Requirements (Cloud-Free K8s)

To run this in-house without the cloud, you will use a **Local Kubernetes Cluster** (MicroK8s, K3s, or Minikube).

### **A. Orchestration & Persistence**

* **K8s Distribution:** **K3s** (Lightweight, perfect for local servers)
* **Ingress Controller:** **NGINX Ingress** to manage local DNS (e.g., accessing the app at `http://private.ai.local`)
* **Storage Class:** **Local Path Provisioner** to map your physical SSD to the PostgreSQL pod so data persists after restarts

### **B. Monitoring & Observability (The "O11y" Stack)**

* **Prometheus:** Configured to scrape metrics from Go app (HTTP metrics, AI latency) and PostgreSQL database
* **PostgreSQL Exporter:** Sidecar container providing database metrics (connections, query performance, storage usage)
* **Grafana:** Custom dashboards for:
  * **Application Performance:** HTTP requests, response times, Go runtime metrics
  * **Database Metrics:** Connection pool, query latency, vector store size
  * **Ollama Health:** GPU/CPU usage per model
  * **Vector Stats:** Number of indexed chunks in PGVector
* **pgAdmin:** Web-based PostgreSQL administration interface for database management
* **Loki + Alloy:** For log aggregation (Centralized logs for K8s pods)

---

## 🐳 4. Docker Development Environment

For local development and testing, the project includes a comprehensive Docker Compose setup:

### **A. Core Services**

* **PostgreSQL + PGVector:** Database with vector extension for semantic search
* **PostgreSQL Exporter:** Metrics collection for database monitoring
* **pgAdmin:** Web-based PostgreSQL management UI (http://localhost:5050)
* **Prometheus:** Metrics collection and storage (http://localhost:9090)
* **Grafana:** Visualization and dashboarding (http://localhost:3000)
* **Redis:** Caching layer for improved performance
* **Ollama:** Local AI model serving (http://localhost:11434)

### **B. Development Workflow**

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f postgres

# Build and run Go backend
cd go-backend
go mod tidy
go run main.go

# Build and run frontend (Angular)
cd angular-ui
npm install
ng serve

# Access services
# PostgreSQL: localhost:5432
# pgAdmin: http://localhost:5050 (admin@knowledge-base.local / admin123)
# Prometheus: http://localhost:9090
# Grafana: http://localhost:3000 (admin / admin123)
# Redis: localhost:6379
# Ollama: http://localhost:11434
# Go Backend: http://localhost:8080
# Angular UI: http://localhost:4200
```

---

## 📁 5. Project Blueprint (Folder Structure)

```text
private-knowledge-base-go/
├── docker-compose.yml          # Local development environment
├── docker-config/              # Docker configuration files
│   ├── postgres/
│   │   └── init-pgvector.sql   # Database initialization script
│   ├── prometheus/
│   │   └── prometheus.yml       # Prometheus configuration
│   └── grafana/
│       ├── provisioning/        # Auto-provisioning configs
│       │   ├── datasources/
│       │   └── dashboards/
│       └── dashboards/          # Pre-built dashboards
├── k8s/                       # Kubernetes manifests
│   ├── monitoring/            # Prometheus/Grafana configs
│   ├── postgres-stateful.yaml # DB with PGVector + exporter
│   ├── pgadmin-deploy.yaml    # PostgreSQL UI management
│   ├── backend-deploy.yaml    # Go backend config
│   └── frontend-deploy.yaml   # React/Nginx config
├── go-backend/                # Go 1.21+ Backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # Application entry point
│   ├── internal/
│   │   ├── ingestion/         # Document parsing logic
│   │   │   ├── pdf.go
│   │   │   ├── markdown.go
│   │   │   ├── text.go
│   │   │   └── gofiles.go
│   │   ├── rag/               # Retrieval & Prompt logic
│   │   │   ├── service.go
│   │   │   ├── vector.go
│   │   │   └── ai.go
│   │   ├── web/               # HTTP Handlers
│   │   │   ├── chat.go
│   │   │   ├── documents.go
│   │   │   └── health.go
│   │   ├── storage/           # Database layer
│   │   │   ├── postgres.go
│   │   │   ├── models.go
│   │   │   └── migrations/
│   │   ├── auth/              # Authentication
│   │   │   ├── jwt.go
│   │   │   └── middleware.go
│   │   └── config/            # Configuration
│   │       ├── config.go
│   │       └── viper.go
│   ├── pkg/                   # Public packages
│   │   ├── logger/
│   │   └── metrics/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── Makefile
├── angular-ui/                # Angular 18 Frontend
│   ├── src/
│   │   ├── app/
│   │   │   ├── components/
│   │   │   │   ├── chat/        # Streaming chat UI
│   │   │   │   ├── upload/      # File management
│   │   │   │   ├── documents/   # Document listing
│   │   │   │   └── core/        # Shared components
│   │   │   ├── services/
│   │   │   │   ├── api.service.ts    # API client
│   │   │   │   ├── sse.service.ts    # SSE handling
│   │   │   │   └── auth.service.ts   # Authentication
│   │   │   ├── models/        # TypeScript interfaces
│   │   │   ├── app.config.ts  # Angular configuration
│   │   │   └── app.component.ts
│   │   ├── assets/            # Static assets
│   │   └── styles/            # Global styles
│   ├── angular.json           # Angular CLI configuration
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.js
│   └── Dockerfile
└── tests/                     # E2E tests and documentation
    ├── e2e/                   # End-to-end test scenarios
    ├── integration/           # Integration tests
    └── TESTING_GUIDE.md       # Testing documentation
```

---

## 🚀 6. Key Go Implementation Details

### **A. Go Backend Architecture**

```go
// Main application structure
type Server struct {
    config      *config.Config
    db          *storage.PostgresDB
    ragService  *rag.Service
    httpServer  *gin.Engine
    logger      *logrus.Logger
}

// Document ingestion service
type IngestionService struct {
    pdfParser     *pdf.Parser
    markdownParser *markdown.Parser
    textParser    *text.Parser
    goParser      *gofiles.Parser
    vectorStore   *storage.VectorStore
}

// RAG service implementation
type RAGService struct {
    vectorStore   *storage.VectorStore
    aiClient      *ai.OllamaClient
    promptManager *prompt.Manager
}
```

### **B. Core Go Packages**

```go
// go.mod dependencies
module private-knowledge-base-go

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/jackc/pgx/v5 v5.4.3
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/spf13/viper v1.16.0
    github.com/sirupsen/logrus v1.9.3
    github.com/unidoc/unipdf/v3 v3.1.1
    github.com/gomarkdown/markdown v0.0.0-20230723190231-5d916f9d1c31
    github.com/google/uuid v1.3.0
    github.com/prometheus/client_golang v1.16.0
    golang.org/x/net v0.12.0
)
```

### **C. Concurrent Processing Pattern**

```go
// Example: Concurrent document processing
func (s *IngestionService) ProcessDocument(ctx context.Context, doc *Document) error {
    // Use worker pool for concurrent chunking
    chunks := make(chan Chunk, 100)
    results := make(chan ProcessedChunk, 100)
    
    // Start workers
    for i := 0; i < runtime.NumCPU(); i++ {
        go s.chunkWorker(ctx, chunks, results)
    }
    
    // Send chunks to workers
    go func() {
        defer close(chunks)
        for _, chunk := range s.splitIntoChunks(doc.Content) {
            select {
            case chunks <- chunk:
            case <-ctx.Done():
                return
            }
        }
    }()
    
    // Collect results
    var processedChunks []ProcessedChunk
    for i := 0; i < len(chunks); i++ {
        select {
        case result := <-results:
            processedChunks = append(processedChunks, result)
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    
    return s.storeChunks(ctx, processedChunks)
}
```

---

## 📊 7. Performance & Scalability Considerations

### **Go Advantages**

* **Lower Memory Footprint:** Go typically uses 50-70% less memory than Java
* **Faster Startup:** Go applications start in milliseconds vs seconds for Java
* **Better Concurrency:** Goroutines are lighter than Java threads
* **Single Binary:** No JVM dependency, easier containerization
* **Compile-Time Optimization:** Better performance for CPU-bound tasks

### **Benchmark Expectations**

| Metric | Java Spring Boot | Go Gin | Improvement |
|--------|------------------|--------|-------------|
| Startup Time | 3-5 seconds | 50-200ms | 10-50x faster |
| Memory Usage | 512MB-1GB | 64-128MB | 4-8x less |
| Request Latency | 10-50ms | 5-20ms | 2-2.5x faster |
| Concurrent Connections | 1000-5000 | 10000+ | 2-10x more |

---

## 🔧 8. Development Workflow

### **Local Development**

```bash
# Backend development
cd go-backend
go mod tidy
go run cmd/server/main.go

# Run tests
go test ./...
go test -race ./...  # Race condition detection
go test -cover ./...  # Coverage report

# Build for production
go build -o bin/server cmd/server/main.go

# Frontend development (Angular)
cd angular-ui
npm install
ng serve

# Run E2E tests
ng e2e
```

### **Docker Development**

```bash
# Build Go backend Docker image
docker build -t knowledge-base-go-backend go-backend/

# Build frontend Docker image
docker build -t knowledge-base-go-frontend angular-ui/

# Run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f go-backend
```

---

## 🚀 9. Deployment Strategy

### **Kubernetes Deployment**

```yaml
# Go Backend Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-backend
  template:
    metadata:
      labels:
        app: go-backend
    spec:
      containers:
      - name: go-backend
        image: knowledge-base-go-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: GIN_MODE
          value: "release"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

---

## 🎯 10. Recommended Next Step

Since you are doing this in-house, the most critical part is setting up the **PostgreSQL + PGVector** container within Kubernetes to store your AI's "memory."

**Would you like me to generate the Kubernetes `StatefulSet` YAML for the PGVector database and the Go backend deployment to get your storage layer ready?**

---

## 📝 11. Migration Notes (Java to Go)

### **Key Differences**

| Aspect | Java Spring Boot | Go |
|--------|------------------|-----|
| Dependency Injection | Spring IoC Container | Constructor injection (manual) |
| ORM | Spring Data JPA | SQL queries with pgx |
| Validation | Bean Validation | Struct tags + custom validation |
| Testing | JUnit + Mockito | Go testing + testify |
| Configuration | application.yml | Viper + environment variables |
| Logging | SLF4J + Logback | Logrus/Zap |
| Metrics | Micrometer | Prometheus client |

### **Migration Strategy**

1. **Database Schema:** Keep the same PostgreSQL schema
2. **API Endpoints:** Maintain identical REST API contracts
3. **Frontend:** **No changes needed** - Angular app works seamlessly with Go backend
4. **Authentication:** JWT tokens remain compatible
5. **Docker Images:** Replace Java image with Go scratch image

### **Benefits of Go Migration**

- **Reduced Resource Usage:** Lower memory and CPU requirements
- **Faster Deployment:** Smaller Docker images, faster startup
- **Better Performance:** Improved latency and throughput
- **Simpler Operations:** Single binary, no JVM management
- **Cost Efficiency:** Lower cloud infrastructure costs

---

This Go-based implementation provides the same functionality as the Java version while leveraging Go's performance advantages, simpler deployment model, and more efficient resource usage. The architecture maintains compatibility with the existing infrastructure while providing a more modern and efficient backend implementation.
