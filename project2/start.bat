@echo off
echo Starting E-Commerce Microservices...
docker-compose up -d
echo.
echo Waiting 30 seconds for initialization...
timeout /t 30 /nobreak >nul
echo.
echo Checking health...
powershell -Command "try { $r=Invoke-RestMethod -Uri 'http://localhost:8081/health' -TimeoutSec 5; if($r.status -eq 'healthy') { Write-Host 'User Service: OK' -ForegroundColor Green } } catch { Write-Host 'User Service: ERROR' -ForegroundColor Red }"
powershell -Command "try { $r=Invoke-RestMethod -Uri 'http://localhost:8082/health' -TimeoutSec 5; if($r.status -eq 'healthy') { Write-Host 'Order Service: OK' -ForegroundColor Green } } catch { Write-Host 'Order Service: ERROR' -ForegroundColor Red }"
powershell -Command "try { $r=Invoke-RestMethod -Uri 'http://localhost:8083/health' -TimeoutSec 5; if($r.status -eq 'healthy') { Write-Host 'Payment Service: OK' -ForegroundColor Green } } catch { Write-Host 'Payment Service: ERROR' -ForegroundColor Red }"
echo.
echo All services started!
echo Test: powershell -Command "Invoke-RestMethod -Uri 'http://localhost:8081/users' -Method POST -ContentType 'application/json' -Body '{\\"name\\":\\"John\\",\\"email\\":\\"john@test.com\\"}'"
pause