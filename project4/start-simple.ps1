# Simple Startup Script for AI Recruiter Assistant

Write-Host "Starting AI Recruiter Assistant..." -ForegroundColor Green

# Check Docker
try {
    docker version > $null 2>&1
    Write-Host "Docker is running" -ForegroundColor Green
} catch {
    Write-Host "Please start Docker Desktop first" -ForegroundColor Red
    exit 1
}

# Start services
Write-Host "Starting Docker services..." -ForegroundColor Blue
docker-compose up -d

# Wait for services
Write-Host "Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 30

# Install dependencies
Write-Host "Installing Go dependencies..." -ForegroundColor Blue
Set-Location go-backend
go mod tidy

Write-Host "Installing Angular dependencies..." -ForegroundColor Blue
Set-Location ../angular-ui
npm install

# Start backend
Write-Host "Starting backend server..." -ForegroundColor Blue
Set-Location ../go-backend
Start-Process powershell -ArgumentList "-Command", "go run cmd/server/main.go" -WindowStyle Minimized

# Start frontend
Write-Host "Starting frontend server..." -ForegroundColor Blue
Set-Location ../angular-ui
Start-Process powershell -ArgumentList "-Command", "npm start" -WindowStyle Minimized

Write-Host ""
Write-Host "AI Recruiter Assistant is starting!" -ForegroundColor Green
Write-Host "Frontend: http://localhost:4200"
Write-Host "Backend: http://localhost:8080"
Write-Host "Grafana: http://localhost:3000 (admin/admin123)"
Write-Host "Press Ctrl+C to stop"
