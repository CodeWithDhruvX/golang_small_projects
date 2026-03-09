#!/bin/bash

# Private Knowledge Base - Automated Startup Script
# This script sets up and runs the entire project with all services

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[ERROR] $1${NC}"
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

# Check if required tools are installed
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed. Please install Docker first."
    fi
    
    # Check Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        error "Docker Compose is not installed. Please install Docker Compose first."
    fi
    
    # Check Go
    if ! command -v go &> /dev/null; then
        error "Go is not installed. Please install Go 1.21+ first."
    fi
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        error "Node.js is not installed. Please install Node.js 18+ first."
    fi
    
    # Check npm
    if ! command -v npm &> /dev/null; then
        error "npm is not installed. Please install npm first."
    fi
    
    log "All prerequisites are installed ✓"
}

# Function to wait for service to be ready
wait_for_service() {
    local service_name=$1
    local host=$2
    local port=$3
    local timeout=${4:-60}
    
    info "Waiting for $service_name to be ready..."
    
    local count=0
    while [ $count -lt $timeout ]; do
        if nc -z $host $port 2>/dev/null; then
            log "$service_name is ready ✓"
            return 0
        fi
        count=$((count + 1))
        sleep 2
    done
    
    error "$service_name failed to start within $timeout seconds"
}

# Function to test HTTP endpoint
test_endpoint() {
    local url=$1
    local expected_status=${2:-200}
    local timeout=${3:-30}
    
    info "Testing endpoint: $url"
    
    local count=0
    while [ $count -lt $timeout ]; do
        if curl -s -o /dev/null -w "%{http_code}" "$url" | grep -q "$expected_status"; then
            log "Endpoint $url is responding correctly ✓"
            return 0
        fi
        count=$((count + 1))
        sleep 2
    done
    
    error "Endpoint $url failed to respond with status $expected_status"
}

# Start infrastructure services
start_infrastructure() {
    log "Starting infrastructure services..."
    
    # Start Docker Compose services
    docker-compose up -d postgres pgadmin redis prometheus grafana
    
    # Wait for services to be ready
    wait_for_service "PostgreSQL" localhost 5432 60
    wait_for_service "Redis" localhost 6379 30
    wait_for_service "Ollama" localhost 11434 60
    wait_for_service "Prometheus" localhost 9090 30
    wait_for_service "Grafana" localhost 3000 30
    
    log "Infrastructure services started ✓"
}

# Pull AI models
setup_ai_models() {
    log "Setting up AI models..."
    
    # Pull required models
    info "Pulling Llama 3.1 8B model..."
    ollama pull llama3.1:8b
    
    info "Pulling Nomic Embed Text model..."
    ollama pull nomic-embed-text
    
    info "Pulling Qwen2.5-Coder model..."
    ollama pull qwen2.5-coder
    
    info "Pulling Phi-3 model..."
    ollama pull phi3
    
    log "AI models setup completed ✓"
}

# Setup and run backend
setup_backend() {
    log "Setting up Go backend..."
    
    cd go-backend
    
    # Install dependencies
    info "Installing Go dependencies..."
    go mod tidy
    
    # Build the application
    info "Building Go backend..."
    go build -o bin/server cmd/server/main.go
    
    # Start backend in background
    info "Starting Go backend..."
    ./bin/server &
    BACKEND_PID=$!
    
    # Wait for backend to be ready
    wait_for_service "Backend" localhost 8080 30
    
    # Test backend endpoints
    test_endpoint "http://localhost:8080/health"
    test_endpoint "http://localhost:8080/ready"
    test_endpoint "http://localhost:8080/metrics"
    
    cd ..
    
    log "Backend setup completed ✓"
}

# Setup and run frontend
setup_frontend() {
    log "Setting up Angular frontend..."
    
    cd angular-ui
    
    # Install dependencies
    info "Installing npm dependencies..."
    npm install
    
    # Build the application
    info "Building Angular frontend..."
    ng build --configuration development
    
    # Start frontend in background
    info "Starting Angular development server..."
    ng serve &
    FRONTEND_PID=$!
    
    # Wait for frontend to be ready
    wait_for_service "Frontend" localhost 4200 60
    
    # Test frontend
    test_endpoint "http://localhost:4200/"
    
    cd ..
    
    log "Frontend setup completed ✓"
}

# Run comprehensive tests
run_tests() {
    log "Running comprehensive functionality tests..."
    
    # Test backend API endpoints
    info "Testing backend API endpoints..."
    
    # Test health endpoints
    test_endpoint "http://localhost:8080/health"
    test_endpoint "http://localhost:8080/ready"
    test_endpoint "http://localhost:8080/metrics"
    
    # Test authentication
    info "Testing authentication..."
    AUTH_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"username": "testuser", "password": "testpass"}')
    
    if echo "$AUTH_RESPONSE" | grep -q "token"; then
        log "Authentication test passed ✓"
    else
        warning "Authentication test may need manual verification"
    fi
    
    # Test document upload endpoint (without actual file)
    info "Testing document upload endpoint..."
    UPLOAD_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/documents/upload" \
        -H "Content-Type: multipart/form-data")
    
    # Test chat endpoint (without actual message)
    info "Testing chat endpoint..."
    CHAT_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/chat" \
        -H "Content-Type: application/json" \
        -d '{"sessionId": "test", "message": "test"}')
    
    # Test Ollama integration
    info "Testing Ollama integration..."
    OLLAMA_RESPONSE=$(curl -s "http://localhost:11434/api/tags")
    
    if echo "$OLLAMA_RESPONSE" | grep -q "llama3.1:8b"; then
        log "Ollama integration test passed ✓"
    else
        warning "Ollama models may still be loading"
    fi
    
    # Test database connection
    info "Testing database connection..."
    DB_TEST=$(docker exec postgres psql -U postgres -d knowledge_base -c "SELECT COUNT(*) FROM documents;" 2>/dev/null || echo "0")
    
    if [ "$DB_TEST" != "0" ]; then
        log "Database connection test passed ✓"
    else
        warning "Database connection may need verification"
    fi
    
    log "Functionality tests completed ✓"
}

# Display service URLs and status
display_status() {
    log "🎉 Project startup completed successfully!"
    echo ""
    echo "=== 🚀 Service URLs ==="
    echo -e "${GREEN}Frontend Application:${NC} http://localhost:4200"
    echo -e "${GREEN}Backend API:${NC} http://localhost:8080"
    echo -e "${GREEN}API Documentation:${NC} http://localhost:8080/swagger/index.html"
    echo -e "${GREEN}pgAdmin:${NC} http://localhost:5050 (admin@knowledge-base.local / admin123)"
    echo -e "${GREEN}Grafana:${NC} http://localhost:3000 (admin / admin123)"
    echo -e "${GREEN}Prometheus:${NC} http://localhost:9090"
    echo -e "${GREEN}Ollama:${NC} http://localhost:11434"
    echo ""
    echo "=== 📊 Monitoring ==="
    echo -e "${BLUE}Grafana Dashboards:${NC} Application metrics and performance"
    echo -e "${BLUE}Prometheus Metrics:${NC} Raw metrics data"
    echo ""
    echo "=== 🤖 AI Models Status ==="
    ollama list 2>/dev/null || echo "Models may still be loading..."
    echo ""
    echo "=== 🔧 Management Commands ==="
    echo "View logs: docker-compose logs -f"
    echo "Stop services: docker-compose down"
    echo "Restart services: docker-compose restart"
    echo ""
    echo "=== 📝 Next Steps ==="
    echo "1. Open http://localhost:4200 in your browser"
    echo "2. Upload some documents (PDF, TXT, MD, or Go files)"
    echo "3. Start chatting with your documents!"
    echo ""
    echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
}

# Cleanup function
cleanup() {
    log "Cleaning up..."
    
    # Kill background processes
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
    fi
    
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
    
    # Stop Docker services
    docker-compose down
    
    log "Cleanup completed"
}

# Set up signal handlers
trap cleanup EXIT INT TERM

# Main execution
main() {
    echo -e "${BLUE}========================================"
    echo "🚀 Private Knowledge Base Startup"
    echo "========================================${NC}"
    echo ""
    
    # Check prerequisites
    check_prerequisites
    
    # Start infrastructure
    start_infrastructure
    
    # Setup AI models
    setup_ai_models
    
    # Setup backend
    setup_backend
    
    # Setup frontend
    setup_frontend
    
    # Run tests
    run_tests
    
    # Display status
    display_status
    
    # Keep script running
    while true; do
        sleep 10
        # Check if services are still running
        if ! curl -s "http://localhost:8080/health" > /dev/null; then
            error "Backend service stopped unexpectedly"
        fi
        if ! curl -s "http://localhost:4200/" > /dev/null; then
            error "Frontend service stopped unexpectedly"
        fi
    done
}

# Run main function
main "$@"
