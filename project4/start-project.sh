#!/bin/bash

# AI Recruiter Assistant - Startup Script for Linux/Mac

echo "🚀 Starting AI Recruiter Assistant..."

# Check if Docker is running
if ! docker version > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

echo "✅ Docker is running"

# Start Docker services
echo "🐳 Starting Docker services..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Check if PostgreSQL is ready
echo "🔍 Checking PostgreSQL..."
for i in {1..10}; do
    if docker exec ai-recruiter-postgres pg_isready -U postgres > /dev/null 2>&1; then
        echo "✅ PostgreSQL is ready"
        break
    fi
    echo "   Attempt $i/10..."
    sleep 5
done

# Pull Ollama models
echo "🤖 Pulling Ollama models..."
docker exec ai-recruiter-ollama ollama pull llama3.1:8b
docker exec ai-recruiter-ollama ollama pull phi3
docker exec ai-recruiter-ollama ollama pull nomic-embed-text

# Install Go dependencies
echo "📦 Installing Go dependencies..."
cd go-backend
go mod tidy

# Install Angular dependencies
echo "📦 Installing Angular dependencies..."
cd ../angular-ui
npm install

# Start backend
echo "🔧 Starting backend server..."
cd ../go-backend
go run cmd/server/main.go &
BACKEND_PID=$!

# Start frontend
echo "🎨 Starting frontend server..."
cd ../angular-ui
npm start &
FRONTEND_PID=$!

echo ""
echo "🎉 AI Recruiter Assistant is starting up!"
echo ""
echo "📊 Access URLs:"
echo "   Frontend:     http://localhost:4200"
echo "   Backend API:  http://localhost:8080"
echo "   API Docs:     http://localhost:8080/swagger/index.html"
echo "   Grafana:      http://localhost:3000 (admin/admin123)"
echo "   Prometheus:   http://localhost:9090"
echo "   pgAdmin:      http://localhost:5050 (admin@ai-recruiter.local/admin123)"
echo ""
echo "🛑 To stop: Press Ctrl+C and run 'docker-compose down'"
echo ""

# Wait for user input to stop
trap 'echo "🛑 Shutting down..."; kill $BACKEND_PID $FRONTEND_PID; docker-compose down; echo "✅ Shutdown complete"; exit 0' INT

wait
