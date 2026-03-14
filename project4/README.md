# AI Recruiter Assistant

A local AI-powered application that processes recruiter emails and generates professional replies using RAG (Retrieval-Augmented Generation).

## Features

- 🤖 **Local AI Processing** - Uses Ollama for private, offline AI operations
- 📧 **Email Ingestion** - IMAP support + manual upload
- 🎯 **Smart Classification** - AI-powered recruiter email detection
- 📝 **Auto-Generated Replies** - Professional email responses
- 📊 **Application Tracking** - Track all job applications
- 🔍 **Semantic Search** - PGVector-powered RAG system
- 📈 **Monitoring** - Prometheus + Grafana observability

## Tech Stack

**Backend**: Go + Gin + PostgreSQL + PGVector + Redis  
**Frontend**: Angular 18 + TailwindCSS  
**AI**: Ollama (Llama 3.1, Phi-3, Nomic Embed)  
**Infrastructure**: Docker + Prometheus + Grafana

## Quick Start

1. Start services:
```bash
docker-compose up -d
```

2. Install dependencies:
```bash
cd go-backend && go mod tidy
cd ../angular-ui && npm install
```

3. Run applications:
```bash
# Backend
cd go-backend && go run cmd/server/main.go

# Frontend  
cd angular-ui && ng serve
```

Access at: http://localhost:4200

## Architecture

Microservices architecture with:
- Authentication (JWT)
- Email processing pipeline
- AI service integration
- Vector database search
- Application tracking

## Documentation

See `/docs` for detailed API docs and deployment guides.
