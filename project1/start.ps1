# start.ps1
Write-Host "Starting services..." -ForegroundColor Green

# Check if Docker is running and start Kafka if available
$dockerRunning = $false
try {
    $dockerVersion = docker --version 2>$null
    if ($dockerVersion) {
        # Test if Docker daemon is actually running
        $dockerInfo = docker info 2>$null
        if ($LASTEXITCODE -eq 0) {
            $dockerRunning = $true
            Write-Host "Docker is running. Starting Kafka infrastructure using Docker Compose..." -ForegroundColor Green
            docker compose up -d
            
            Write-Host "Waiting a few seconds for Kafka to initialize..." -ForegroundColor Yellow
            Start-Sleep -Seconds 10
            Write-Host "Kafka infrastructure started successfully!" -ForegroundColor Green
        } else {
            Write-Host "Docker command found but Docker daemon is not running." -ForegroundColor Yellow
        }
    } else {
        Write-Host "Docker not found." -ForegroundColor Yellow
    }
} catch {
    Write-Host "Docker is not available." -ForegroundColor Yellow
}

if (-not $dockerRunning) {
    Write-Host "Services will run with mock Kafka (for development/testing)." -ForegroundColor Yellow
}

Write-Host "Starting Notification Service (Consumer)..." -ForegroundColor Green
Start-Process -FilePath "go" -ArgumentList "run cmd/notification-service/main.go" -NoNewWindow

Write-Host "Starting Order Service (Producer API)..." -ForegroundColor Green
Start-Process -FilePath "go" -ArgumentList "run cmd/order-service/main.go" -NoNewWindow

Write-Host "All services started successfully!" -ForegroundColor Green
Write-Host "You can test the API using your preferred tool (e.g., Postman or curl on port 8080)." -ForegroundColor Cyan
Write-Host "Note: If Docker is not running, services will use mock Kafka for testing." -ForegroundColor Cyan
