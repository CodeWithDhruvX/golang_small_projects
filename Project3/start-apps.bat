@echo off
echo Starting all applications locally...

echo.
echo [1/3] Starting Docker services...
start "Docker Services" cmd /k "cd /d %~dp0 && docker-compose up -d"
echo Waiting for Docker services to start...
timeout /t 30 /nobreak

echo.
echo [2/3] Starting Angular Frontend...
start "Angular Frontend" cmd /k "cd /d %~dp0\angular-ui && npm start"

echo.
echo [3/3] Starting Go Backend...
start "Go Backend" cmd /k "cd /d %~dp0\go-backend && go run cmd/server/main.go"

echo.
echo All applications are starting in separate terminals...
echo - Docker: PostgreSQL, Redis, Prometheus, Grafana
echo - Frontend: http://localhost:4200
echo - Backend: http://localhost:8080
echo.
pause
