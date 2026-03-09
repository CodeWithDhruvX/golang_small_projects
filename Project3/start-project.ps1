# Private Knowledge Base - Automated Startup Script (PowerShell)
# This script sets up and runs the entire project with all services

param(
    [switch]$SkipTests,
    [switch]$DevMode
)

# Error handling
$ErrorActionPreference = "Stop"

# Colors for output
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Log-Info($message) {
    Write-ColorOutput Green "[$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')] $message"
}

function Log-Error($message) {
    Write-ColorOutput Red "[ERROR] $message"
    throw $message
}

function Log-Warning($message) {
    Write-ColorOutput Yellow "[WARNING] $message"
}

function Log-Info-Blue($message) {
    Write-ColorOutput Blue "[INFO] $message"
}

# Check if required tools are installed
function Test-Prerequisites {
    Log-Info "Checking prerequisites..."
    
    # Check Docker
    if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
        Log-Error "Docker is not installed. Please install Docker Desktop first."
    }
    
    # Check Docker Compose
    if (-not (Get-Command docker-compose -ErrorAction SilentlyContinue)) {
        Log-Error "Docker Compose is not installed. Please install Docker Compose first."
    }
    
    # Check Go
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        Log-Error "Go is not installed. Please install Go 1.21+ first."
    }
    
    # Check Node.js
    if (-not (Get-Command node -ErrorAction SilentlyContinue)) {
        Log-Error "Node.js is not installed. Please install Node.js 18+ first."
    }
    
    # Check npm
    if (-not (Get-Command npm -ErrorAction SilentlyContinue)) {
        Log-Error "npm is not installed. Please install npm first."
    }
    
    Log-Info "All prerequisites are installed ✓"
}

# Function to wait for service to be ready
function Wait-ForService($serviceName, $host, $port, $timeout = 60) {
    Log-Info-Blue "Waiting for $serviceName to be ready..."
    
    $count = 0
    $client = New-Object System.Net.Sockets.TcpClient
    
    while ($count -lt $timeout) {
        try {
            $client.Connect($host, $port)
            $client.Close()
            Log-Info "$serviceName is ready ✓"
            return
        }
        catch {
            $count++
            Start-Sleep 2
        }
    }
    
    Log-Error "$serviceName failed to start within $timeout seconds"
}

# Function to test HTTP endpoint
function Test-Endpoint($url, $expectedStatus = 200, $timeout = 30) {
    Log-Info-Blue "Testing endpoint: $url"
    
    $count = 0
    while ($count -lt $timeout) {
        try {
            $response = Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 5
            if ($response.StatusCode -eq $expectedStatus) {
                Log-Info "Endpoint $url is responding correctly ✓"
                return
            }
        }
        catch {
            # Continue trying
        }
        $count++
        Start-Sleep 2
    }
    
    Log-Error "Endpoint $url failed to respond with status $expectedStatus"
}

# Start infrastructure services
function Start-Infrastructure {
    Log-Info "Starting infrastructure services..."
    
    # Start Docker Compose services
    docker-compose up -d postgres pgadmin redis prometheus grafana
    
    try {
        Get-Process ollama -ErrorAction SilentlyContinue | Stop-Process -Force
    } catch {}
    $script:ollamaJob = Start-Job -ScriptBlock {
        $env:OLLAMA_KEEP_ALIVE = "30m"
        $env:OLLAMA_NUM_PARALLEL = "2"
        $env:OLLAMA_MAX_LOADED_MODELS = "2"
        $env:OLLAMA_FLASH_ATTENTION = "1"
        $env:OLLAMA_KV_CACHE_TYPE = "q8_0"
        $env:OLLAMA_SCHED_SPREAD = "1"
        ollama serve
    }
    
    # Wait for services to be ready
    Wait-ForService "PostgreSQL" "localhost" 5432 60
    Wait-ForService "Redis" "localhost" 6379 30
    Wait-ForService "Ollama" "localhost" 11434 60
    Wait-ForService "Prometheus" "localhost" 9090 30
    Wait-ForService "Grafana" "localhost" 3000 30
    
    Log-Info "Infrastructure services started ✓"
}

# Pull AI models
function Set-AIModels {
    Log-Info "Setting up AI models..."
    
    # Pull required models
    Log-Info-Blue "Pulling Llama 3.1 8B model..."
    ollama pull llama3.1:8b
    
    Log-Info-Blue "Pulling Nomic Embed Text model..."
    ollama pull nomic-embed-text
    
    Log-Info-Blue "Pulling Qwen2.5-Coder model..."
    ollama pull qwen2.5-coder
    
    Log-Info-Blue "Pulling Phi-3 model..."
    ollama pull phi3
    
    Log-Info "AI models setup completed ✓"
}

# Setup and run backend
function Set-Backend {
    Log-Info "Setting up Go backend..."
    
    Push-Location go-backend
    
    try {
        # Install dependencies
        Log-Info-Blue "Installing Go dependencies..."
        go mod tidy
        
        # Build the application
        Log-Info-Blue "Building Go backend..."
        go build -o bin/server cmd/server/main.go
        
        # Start backend in background
        Log-Info-Blue "Starting Go backend..."
        $backendJob = Start-Job -ScriptBlock {
            Set-Location $using:PWD
            .\bin\server
        }
        
        # Wait for backend to be ready
        Wait-ForService "Backend" "localhost" 8080 30
        
        # Test backend endpoints
        Test-Endpoint "http://localhost:8080/health"
        Test-Endpoint "http://localhost:8080/ready"
        Test-Endpoint "http://localhost:8080/metrics"
        
        Log-Info "Backend setup completed ✓"
    }
    finally {
        Pop-Location
    }
}

# Setup and run frontend
function Set-Frontend {
    Log-Info "Setting up Angular frontend..."
    
    Push-Location angular-ui
    
    try {
        # Install dependencies
        Log-Info-Blue "Installing npm dependencies..."
        npm install
        
        if ($DevMode) {
            # Start frontend in development mode
            Log-Info-Blue "Starting Angular development server..."
            $frontendJob = Start-Job -ScriptBlock {
                Set-Location $using:PWD
                npm start
            }
        } else {
            # Build the application
            Log-Info-Blue "Building Angular frontend..."
            ng build --configuration development
            
            # Start frontend in background
            Log-Info-Blue "Starting Angular development server..."
            $frontendJob = Start-Job -ScriptBlock {
                Set-Location $using:PWD
                ng serve
            }
        }
        
        # Wait for frontend to be ready
        Wait-ForService "Frontend" "localhost" 4200 60
        
        # Test frontend
        Test-Endpoint "http://localhost:4200/"
        
        Log-Info "Frontend setup completed ✓"
    }
    finally {
        Pop-Location
    }
}

# Run comprehensive tests
function Invoke-Tests {
    if ($SkipTests) {
        Log-Warning "Skipping functionality tests as requested"
        return
    }
    
    Log-Info "Running comprehensive functionality tests..."
    
    # Test backend API endpoints
    Log-Info-Blue "Testing backend API endpoints..."
    
    # Test health endpoints
    Test-Endpoint "http://localhost:8080/health"
    Test-Endpoint "http://localhost:8080/ready"
    Test-Endpoint "http://localhost:8080/metrics"
    
    # Test authentication
    Log-Info-Blue "Testing authentication..."
    try {
        $authResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" `
            -Method POST `
            -ContentType "application/json" `
            -Body '{"username": "testuser", "password": "testpass"}' `
            -TimeoutSec 10
        
        if ($authResponse.token) {
            Log-Info "Authentication test passed ✓"
        } else {
            Log-Warning "Authentication test may need manual verification"
        }
    }
    catch {
        Log-Warning "Authentication test failed - may need manual verification"
    }
    
    # Test Ollama integration
    Log-Info-Blue "Testing Ollama integration..."
    try {
        $ollamaResponse = Invoke-RestMethod -Uri "http://localhost:11434/api/tags" -TimeoutSec 10
        
        if ($ollamaResponse.models -and ($ollamaResponse.models | Where-Object { $_.name -like "*llama3.1:8b*" })) {
            Log-Info "Ollama integration test passed ✓"
        } else {
            Log-Warning "Ollama models may still be loading"
        }
    }
    catch {
        Log-Warning "Ollama integration test failed - models may still be loading"
    }
    
    # Test database connection
    Log-Info-Blue "Testing database connection..."
    try {
        $dbTest = docker exec postgres psql -U postgres -d knowledge_base -c "SELECT COUNT(*) FROM documents;" 2>$null
        if ($dbTest -and $dbTest -match "\d+") {
            Log-Info "Database connection test passed ✓"
        } else {
            Log-Warning "Database connection may need verification"
        }
    }
    catch {
        Log-Warning "Database connection test failed - may need manual verification"
    }
    
    Log-Info "Functionality tests completed ✓"
}

# Display service URLs and status
function Show-Status {
    Log-Info "🎉 Project startup completed successfully!"
    Write-Host ""
    Write-Host "=== 🚀 Service URLs ===" -ForegroundColor Cyan
    Write-Host "Frontend Application: http://localhost:4200" -ForegroundColor Green
    Write-Host "Backend API: http://localhost:8080" -ForegroundColor Green
    Write-Host "API Documentation: http://localhost:8080/swagger/index.html" -ForegroundColor Green
    Write-Host "pgAdmin: http://localhost:5050 (admin@knowledge-base.local / admin123)" -ForegroundColor Green
    Write-Host "Grafana: http://localhost:3000 (admin / admin123)" -ForegroundColor Green
    Write-Host "Prometheus: http://localhost:9090" -ForegroundColor Green
    Write-Host "Ollama: http://localhost:11434" -ForegroundColor Green
    Write-Host ""
    Write-Host "=== 📊 Monitoring ===" -ForegroundColor Cyan
    Write-Host "Grafana Dashboards: Application metrics and performance" -ForegroundColor Blue
    Write-Host "Prometheus Metrics: Raw metrics data" -ForegroundColor Blue
    Write-Host ""
    Write-Host "=== 🤖 AI Models Status ===" -ForegroundColor Cyan
    try {
        ollama list
    }
    catch {
        Write-Host "Models may still be loading..." -ForegroundColor Yellow
    }
    Write-Host ""
    Write-Host "=== 🔧 Management Commands ===" -ForegroundColor Cyan
    Write-Host "View logs: docker-compose logs -f" -ForegroundColor Blue
    Write-Host "Stop services: docker-compose down" -ForegroundColor Blue
    Write-Host "Restart services: docker-compose restart" -ForegroundColor Blue
    Write-Host ""
    Write-Host "=== 📝 Next Steps ===" -ForegroundColor Cyan
    Write-Host "1. Open http://localhost:4200 in your browser" -ForegroundColor Yellow
    Write-Host "2. Upload some documents (PDF, TXT, MD, or Go files)" -ForegroundColor Yellow
    Write-Host "3. Start chatting with your documents!" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Press Ctrl+C to stop all services" -ForegroundColor Yellow
}

# Cleanup function
function Stop-Services {
    Log-Info "Cleaning up..."
    
    # Stop background jobs
    Get-Job | Stop-Job
    Get-Job | Remove-Job
    
    # Stop Docker services
    docker-compose down
    
    Log-Info "Cleanup completed"
}

# Main execution
function Main {
    Write-Host "========================================" -ForegroundColor Blue
    Write-Host "🚀 Private Knowledge Base Startup" -ForegroundColor Blue
    Write-Host "========================================" -ForegroundColor Blue
    Write-Host ""
    
    try {
        # Check prerequisites
        Test-Prerequisites
        
        # Start infrastructure
        Start-Infrastructure
        
        # Setup AI models
        Set-AIModels
        
        # Setup backend
        Set-Backend
        
        # Setup frontend
        Set-Frontend
        
        # Run tests
        Invoke-Tests
        
        # Display status
        Show-Status
        
        # Keep script running
        Write-Host "Monitoring services... Press Ctrl+C to stop" -ForegroundColor Green
        while ($true) {
            Start-Sleep 10
            
            # Check if services are still running
            try {
                $backendHealth = Invoke-WebRequest -Uri "http://localhost:8080/health" -UseBasicParsing -TimeoutSec 5
                $frontendHealth = Invoke-WebRequest -Uri "http://localhost:4200/" -UseBasicParsing -TimeoutSec 5
            }
            catch {
                Log-Error "One or more services stopped unexpectedly"
            }
        }
    }
    catch {
        Log-Error "Startup failed: $($_.Exception.Message)"
    }
    finally {
        Stop-Services
    }
}

# Set up signal handlers
$null = Register-EngineEvent -SourceIdentifier PowerShell.Exiting -Action { Stop-Services }

# Run main function
Main
