# start-services.ps1 - One shot to start everything
param(
    [switch]$Rebuild,
    [switch]$Reset
)

Write-Host "=== Starting E-Commerce Microservices ===" -ForegroundColor Cyan

# Reset if requested
if ($Reset) {
    Write-Host "Resetting all data..." -ForegroundColor Yellow
    docker-compose down -v 2>$null
}

# Build if requested or first time
if ($Rebuild -or -not (docker images project2-user-service -q)) {
    Write-Host "Building Docker images..." -ForegroundColor Yellow
    docker-compose build --no-cache
}

# Start services
Write-Host "Starting services..." -ForegroundColor Green
docker-compose up -d

# Wait for initialization
Write-Host "Waiting for services to initialize..." -ForegroundColor Yellow
for ($i = 30; $i -gt 0; $i--) {
    Write-Host "`r$i seconds remaining..." -NoNewline -ForegroundColor Gray
    Start-Sleep -Seconds 1
}
Write-Host "`rServices starting...          " -ForegroundColor Green

# Health check
Write-Host "`n=== Health Checks ===" -ForegroundColor Cyan
$services = @(
    @{Name="User Service"; Url="http://localhost:8081/health"; Port=8081},
    @{Name="Order Service"; Url="http://localhost:8082/health"; Port=8082},
    @{Name="Payment Service"; Url="http://localhost:8083/health"; Port=8083}
)

$allHealthy = $true
foreach ($svc in $services) {
    try {
        $response = Invoke-RestMethod -Uri $svc.Url -TimeoutSec 5
        if ($response.status -eq "healthy") {
            Write-Host "OK: $($svc.Name)" -ForegroundColor Green
        } else {
            Write-Host "UNHEALTHY: $($svc.Name)" -ForegroundColor Yellow
            $allHealthy = $false
        }
    } catch {
        Write-Host "ERROR: $($svc.Name)" -ForegroundColor Red
        $allHealthy = $false
    }
}

if ($allHealthy) {
    Write-Host "`n=== All Services Ready! ===" -ForegroundColor Green
    Write-Host "URLs:"
    Write-Host "  User Service:    http://localhost:8081"
    Write-Host "  Order Service:   http://localhost:8082"
    Write-Host "  Payment Service: http://localhost:8083"
    Write-Host "`nQuick test:"
    Write-Host '  Invoke-RestMethod -Uri "http://localhost:8081/users" -Method POST -ContentType "application/json" -Body '\''{"name":"Test","email":"test@test.com"}'\''
} else {
    Write-Host "`nSome services are not ready yet. Check logs with: docker-compose logs -f" -ForegroundColor Yellow
}

Write-Host "`nTo stop: docker-compose down" -ForegroundColor Gray