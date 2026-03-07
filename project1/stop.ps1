# stop.ps1

Write-Host "Stopping Go microservices (Consumer & Producer)..." -ForegroundColor Yellow
Stop-Process -Name "go" -Force -ErrorAction SilentlyContinue
Stop-Process -Name "main" -Force -ErrorAction SilentlyContinue

Write-Host "Stopping Kafka infrastructure using Docker Compose..." -ForegroundColor Yellow
docker compose down

Write-Host "All services stopped successfully!" -ForegroundColor Green
