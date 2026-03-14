# Ollama Startup Script with Performance Optimizations (PowerShell)
# Based on Project 3 optimization results

Write-Host "Starting Ollama with performance optimizations..." -ForegroundColor Green

# Set optimization environment variables
$env:OLLAMA_KEEP_ALIVE = "30m"
$env:OLLAMA_NUM_PARALLEL = "2"
$env:OLLAMA_FLASH_ATTENTION = "1"
$env:OLLAMA_KV_CACHE_TYPE = "q8_0"
$env:OLLAMA_MAX_LOADED_MODELS = "2"
$env:OLLAMA_HOST = "127.0.0.1:11434"
$env:OLLAMA_SCHED_SPREAD = "1"

Write-Host "Environment variables set:" -ForegroundColor Yellow
Write-Host "OLLAMA_KEEP_ALIVE=$env:OLLAMA_KEEP_ALIVE"
Write-Host "OLLAMA_NUM_PARALLEL=$env:OLLAMA_NUM_PARALLEL"
Write-Host "OLLAMA_FLASH_ATTENTION=$env:OLLAMA_FLASH_ATTENTION"
Write-Host "OLLAMA_KV_CACHE_TYPE=$env:OLLAMA_KV_CACHE_TYPE"
Write-Host "OLLAMA_MAX_LOADED_MODELS=$env:OLLAMA_MAX_LOADED_MODELS"
Write-Host "OLLAMA_HOST=$env:OLLAMA_HOST"
Write-Host "OLLAMA_SCHED_SPREAD=$env:OLLAMA_SCHED_SPREAD"

# Check if required models are available
Write-Host "Checking for required models..." -ForegroundColor Yellow

# Check for llama3.1:8b model
$llamaCheck = ollama list | Select-String "llama3.1:8b"
if (-not $llamaCheck) {
    Write-Host "Pulling llama3.1:8b model..." -ForegroundColor Cyan
    ollama pull llama3.1:8b
} else {
    Write-Host "llama3.1:8b model already available" -ForegroundColor Green
}

# Check for nomic-embed-text model
$nomicCheck = ollama list | Select-String "nomic-embed-text"
if (-not $nomicCheck) {
    Write-Host "Pulling nomic-embed-text model..." -ForegroundColor Cyan
    ollama pull nomic-embed-text
} else {
    Write-Host "nomic-embed-text model already available" -ForegroundColor Green
}

Write-Host "Starting Ollama server..." -ForegroundColor Green
ollama serve
