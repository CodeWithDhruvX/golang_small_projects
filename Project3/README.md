# Private Knowledge Base - AI-Powered Document Intelligence

A completely **in-house, cloud-free** AI knowledge base built with Go and Angular 18. Transform your local hardware into a private "intelligence" hub using local AI models for document processing and intelligent search.

## 🚀 Features

### **Core Capabilities**
- **Multi-Format Document Ingestion**: PDF, TXT, Markdown, and Go source files
- **Intelligent Chunking**: Smart document segmentation with 20% overlap for context preservation
- **Vector Search**: Semantic similarity search using Nomic-Embed-Text
- **RAG-Powered Chat**: Retrieval-Augmented Generation with source citations
- **Real-Time Streaming**: Server-Sent Events for responsive chat experience
- **Multiple AI Models**: Llama 3.1 8B, Qwen2.5-Coder, and Phi-3 support

### **Technical Stack**
- **Backend**: Go 1.21+ with Gin framework
- **Frontend**: Angular 18 with Tailwind CSS
- **Database**: PostgreSQL + PGVector for vector storage
- **AI**: Local Ollama integration
- **Monitoring**: Prometheus + Grafana
- **Deployment**: Kubernetes (K3s/MicroK8s ready)

## 📋 Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Angular 18    │    │   Go Backend    │    │   PostgreSQL    │
│                 │◄──►│                 │◄──►│   + PGVector    │
│  - Chat UI      │    │  - RAG Service  │    │                 │
│  - Upload UI    │    │  - Ingestion    │    │  - Documents    │
│  - Dark Theme   │    │  - Auth JWT     │    │  - Chunks       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │     Ollama      │
                    │                 │
                    │ - Llama 3.1 8B  │
                    │ - Qwen2.5-Coder │
                    │ - Phi-3         │
                    └─────────────────┘
```

## 🛠️ Quick Start

### Prerequisites

- **Go 1.21+**
- **Node.js 18+**
- **Docker & Docker Compose**
- **Ollama** (for local AI models)
- **PostgreSQL** (or use Docker)

### 1. Clone and Setup

```bash
git clone <repository-url>
cd private-knowledge-base-go
```

### 2. Start Infrastructure

```bash
# Start all services with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f
```

Services started:
- **PostgreSQL**: `localhost:5432`
- **pgAdmin**: `http://localhost:5050` (admin@knowledge-base.local / admin123)
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (admin / admin123)
- **Redis**: `localhost:6379`
- **Ollama**: `http://localhost:11434`

### 3. Install AI Models

```bash
# Pull required models
docker exec ollama ollama pull llama3.1:8b
docker exec ollama ollama pull nomic-embed-text
docker exec ollama ollama pull qwen2.5-coder
docker exec ollama ollama pull phi3
```

### 4. Run Backend

```bash
cd go-backend
go mod tidy
go run cmd/server/main.go
```

Backend runs on: `http://localhost:8080`

### 5. Run Frontend

```bash
cd angular-ui
npm install
ng serve
```

Frontend runs on: `http://localhost:4200`

### 6. Access the Application

Open your browser and navigate to `http://localhost:4200`

- **Default Login**: Any username/password (demo mode)
- **Upload Documents**: PDF, TXT, MD, or Go files
- **Start Chatting**: Ask questions about your documents

## 🐳 Docker Development

### Build Images

```bash
# Build Go backend
docker build -t knowledge-base-go-backend go-backend/

# Build Angular frontend  
docker build -t knowledge-base-go-frontend angular-ui/
```

### Production Deployment

```bash
# Deploy to Kubernetes
kubectl apply -f k8s/

# Or use Docker Compose for production
docker-compose -f docker-compose.prod.yml up -d
```

## ☸️ Kubernetes Deployment

### Prerequisites

- **K3s**, **MicroK8s**, or **Minikube**
- **Ingress Controller** (NGINX)
- **Persistent Storage** (Local Path Provisioner)

### Deploy

```bash
# Create namespace
kubectl create namespace knowledge-base

# Deploy all components
kubectl apply -f k8s/postgres-stateful.yaml
kubectl apply -f k8s/backend-deploy.yaml  
kubectl apply -f k8s/frontend-deploy.yaml
kubectl apply -f k8s/ingress.yaml
kubectl apply -f k8s/monitoring/

# Check status
kubectl get pods -n knowledge-base
```

### Access Services

- **Application**: `http://knowledge-base.local`
- **Grafana**: `http://knowledge-base.local/grafana`
- **Prometheus**: `http://knowledge-base.local/prometheus`

## 📊 Monitoring & Observability

### Metrics Collected

- **Application Metrics**: HTTP requests, response times, error rates
- **Database Metrics**: Connections, query performance, storage usage  
- **AI Metrics**: Request latency, token usage, model performance
- **System Metrics**: CPU, memory, disk usage

### Grafana Dashboards

1. **Knowledge Base Overview**: Application health and performance
2. **Database Monitoring**: PostgreSQL performance and storage
3. **AI/ML Metrics**: Model usage and response times

### Alerts

- **Service Health**: Backend/frontend downtime
- **Performance**: High latency or error rates
- **Resources**: Memory/CPU threshold breaches
- **Database**: Connection issues or storage alerts

## 🔧 Configuration

### Environment Variables

**Backend (Go)**:
```bash
PORT=8080
ENVIRONMENT=production
DATABASE_URL=postgres://user:pass@localhost:5432/knowledge_base
OLLAMA_URL=http://localhost:11434
JWT_SECRET=your-secret-key
```

**Frontend (Angular)**:
```bash
NG_APP_API_URL=http://localhost:8080/api
NG_APP_OLLAMA_URL=http://localhost:11434
```

### Database Schema

The system uses PostgreSQL with PGVector extension:

- **documents**: File metadata and processing status
- **document_chunks**: Vector embeddings for semantic search
- **chat_sessions**: Conversation management
- **chat_messages**: Message history with citations

## 🧪 Development

### Backend Development

```bash
cd go-backend

# Run tests
go test ./...
go test -race ./...
go test -cover ./...

# Build
go build -o bin/server cmd/server/main.go

# Lint
golangci-lint run

# Format
go fmt ./...
```

### Frontend Development

```bash
cd angular-ui

# Install dependencies
npm install

# Run development server
ng serve

# Run tests
ng test
ng e2e

# Build for production
ng build --configuration production

# Lint
ng lint

# Format
npm run format
```

### API Documentation

Swagger documentation is available at:
- **Development**: `http://localhost:8080/swagger/index.html`
- **Production**: `http://your-domain/swagger/index.html`

## 🔒 Security

### Authentication

- **JWT-based** authentication with configurable secrets
- **Token expiration** and refresh mechanisms
- **Role-based** access control ready

### Data Protection

- **Local-only** processing (no cloud dependencies)
- **Encrypted** connections (HTTPS in production)
- **Container security** with non-root users
- **Network policies** for Kubernetes isolation

### Best Practices

- **Input validation** on all endpoints
- **SQL injection** prevention with parameterized queries
- **XSS protection** with content security policies
- **Rate limiting** for API endpoints

## 📈 Performance

### Benchmarks

| Metric | Expected Performance |
|--------|---------------------|
| Backend Startup | 50-200ms |
| Memory Usage | 64-128MB |
| Request Latency | 5-20ms |
| Concurrent Connections | 10,000+ |
| Document Processing | 1-5MB/s |

### Optimization

- **Concurrent processing** with Go goroutines
- **Vector indexing** with PGVector IVFFlat
- **Caching** with Redis for frequent queries
- **Lazy loading** for large document sets

## 🚀 Production Deployment

### System Requirements

**Minimum**:
- CPU: 2 cores
- Memory: 4GB RAM
- Storage: 20GB SSD
- Network: 100Mbps

**Recommended**:
- CPU: 4+ cores
- Memory: 8GB+ RAM
- Storage: 50GB+ SSD
- GPU: For AI model acceleration

### Scaling

- **Horizontal scaling** with Kubernetes HPA
- **Database sharding** for large document sets
- **Load balancing** with multiple backend instances
- **CDN integration** for static assets

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go and Angular best practices
- Write tests for new features
- Update documentation
- Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

### Troubleshooting

**Common Issues**:

1. **Ollama Connection Failed**
   ```bash
   # Check if Ollama is running
   docker-compose ps ollama
   
   # Restart Ollama
   docker-compose restart ollama
   ```

2. **Database Connection Issues**
   ```bash
   # Check PostgreSQL logs
   docker-compose logs postgres
   
   # Verify database exists
   docker exec postgres psql -U postgres -d knowledge_base
   ```

3. **Frontend Build Errors**
   ```bash
   # Clear node modules
   rm -rf node_modules package-lock.json
   npm install
   ```

### Getting Help

- **Documentation**: Check this README and inline code comments
- **Issues**: Create an issue on GitHub with detailed description
- **Community**: Join our discussions for questions and ideas

## 🗺️ Roadmap

### Upcoming Features

- [ ] **Multi-tenant support** for organizations
- [ ] **Advanced search** with filters and facets
- [ ] **Document versioning** and change tracking
- [ ] **Webhook integrations** for automation
- [ ] **Mobile app** (React Native)
- [ ] **Advanced analytics** and reporting
- [ ] **Plugin system** for custom processors

### Technical Improvements

- [ ] **GraphQL API** for flexible queries
- [ ] **Event sourcing** for audit trails
- [ ] **Microservices** architecture
- [ ] **Edge computing** support
- [ ] **Advanced caching** strategies

---

**Built with ❤️ for private, secure, and intelligent document management**
