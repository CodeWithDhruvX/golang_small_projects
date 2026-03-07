# start.ps1 - Simple one-shot start
Write-Host "Starting E-Commerce Microservices..." -ForegroundColor Green

# Stop any existing containers first
docker-compose down 2>$null

# Start services
docker-compose up -d

Write-Host "Waiting 30 seconds for initialization..." -ForegroundColor Yellow
Start-Sleep -Seconds 30

Write-Host "`nChecking health..." -ForegroundColor Cyan

# Check User Service
try {
    $r = Invoke-RestMethod -Uri "http://localhost:8081/health" -TimeoutSec 5
    if ($r.status -eq "healthy") { Write-Host "User Service:    OK" -ForegroundColor Green }
    else { Write-Host "User Service:    UNHEALTHY" -ForegroundColor Yellow }
} catch { Write-Host "User Service:    ERROR" -ForegroundColor Red }

# Check Order Service
try {
    $r = Invoke-RestMethod -Uri "http://localhost:8082/health" -TimeoutSec 5
    if ($r.status -eq "healthy") { Write-Host "Order Service:   OK" -ForegroundColor Green }
    else { Write-Host "Order Service:   UNHEALTHY" -ForegroundColor Yellow }
} catch { Write-Host "Order Service:   ERROR" -ForegroundColor Red }

# Check Payment Service
try {
    $r = Invoke-RestMethod -Uri "http://localhost:8083/health" -TimeoutSec 5
    if ($r.status -eq "healthy") { Write-Host "Payment Service: OK" -ForegroundColor Green }
    else { Write-Host "Payment Service: UNHEALTHY" -ForegroundColor Yellow }
} catch { Write-Host "Payment Service: ERROR" -ForegroundColor Red }

Write-Host "`nTest command:" -ForegroundColor Cyan
Write-Host 'Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body (ConvertTo-Json @{name="Test";email="test@test.com"})'
Write-Host "`nTo stop: docker-compose down" -ForegroundColor Gray
