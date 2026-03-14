# AI Recruiter Assistant - Startup Script for Windows PowerShell

Write-Host "🚀 Starting AI Recruiter Assistant..." -ForegroundColor Green

# Check if Docker is running
try {
    docker version > $null 2>&1
    Write-Host "✅ Docker is running" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker is not running. Please start Docker Desktop first." -ForegroundColor Red
    exit 1
}

# Start Docker services
Write-Host "🐳 Starting Docker services..." -ForegroundColor Blue
docker-compose up -d

# Wait for services to be ready
Write-Host "⏳ Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 30

# Check if PostgreSQL is ready
Write-Host "🔍 Checking PostgreSQL..." -ForegroundColor Blue
$pgReady = $false
for ($i = 1; $i -le 10; $i++) {
    try {
        docker exec ai-recruiter-postgres pg_isready -U postgres > $null 2>&1
        $pgReady = $true
        break
    } catch {
        Write-Host "   Attempt $i/10..." -ForegroundColor Yellow
        Start-Sleep -Seconds 5
    }
}

if ($pgReady) {
    Write-Host "✅ PostgreSQL is ready" -ForegroundColor Green
} else {
    Write-Host "❌ PostgreSQL failed to start" -ForegroundColor Red
    exit 1
}

# Pull Ollama models
Write-Host "🤖 Pulling Ollama models..." -ForegroundColor Blue
docker exec ai-recruiter-ollama ollama pull llama3.1:8b
docker exec ai-recruiter-ollama ollama pull phi3
docker exec ai-recruiter-ollama ollama pull nomic-embed-text

# Install Go dependencies
Write-Host "📦 Installing Go dependencies..." -ForegroundColor Blue
Set-Location go-backend
go mod tidy

# Install Angular dependencies
Write-Host "📦 Installing Angular dependencies..." -ForegroundColor Blue
Set-Location ../angular-ui
npm install

# Start backend (in background)
Write-Host "🔧 Starting backend server..." -ForegroundColor Blue
Set-Location ../go-backend
$backendJob = Start-Job -ScriptBlock {
    Set-Location $using:PWD
    go run cmd/server/main.go
}

# Start frontend (in background)
Write-Host "🎨 Starting frontend server..." -ForegroundColor Blue
Set-Location ../angular-ui
$frontendJob = Start-Job -ScriptBlock {
    Set-Location $using:PWD
    npm start
}

Write-Host ""
Write-Host "🎉 AI Recruiter Assistant is starting up!" -ForegroundColor Green
Write-Host ""
Write-Host "📊 Access URLs:" -ForegroundColor Cyan
Write-Host "   Frontend:     http://localhost:4200" -ForegroundColor White
Write-Host "   Backend API:  http://localhost:8080" -ForegroundColor White
Write-Host "   API Docs:     http://localhost:8080/swagger/index.html" -ForegroundColor White
Write-Host "   Grafana:      http://localhost:3000 (admin/admin123)" -ForegroundColor White
Write-Host "   Prometheus:   http://localhost:9090" -ForegroundColor White
Write-Host "   pgAdmin:      http://localhost:5050 (admin@ai-recruiter.local/admin123)" -ForegroundColor White
Write-Host ""
Write-Host "🛑 To stop: Press Ctrl+C and run 'docker-compose down'" -ForegroundColor Yellow
Write-Host ""

# Wait for user input to stop
try {
    while ($true) {
        Start-Sleep -Seconds 1
    }
} finally {
    Write-Host "🛑 Shutting down..." -ForegroundColor Red
    Stop-Job $backendJob -ErrorAction SilentlyContinue
    Stop-Job $frontendJob -ErrorAction SilentlyContinue
    Remove-Job $backendJob -ErrorAction SilentlyContinue
    Remove-Job $frontendJob -ErrorAction SilentlyContinue
    docker-compose down
    Write-Host "✅ Shutdown complete" -ForegroundColor Green
}
