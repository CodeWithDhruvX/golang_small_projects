#!/bin/bash

# Ollama Startup Script with Performance Optimizations
# Based on Project 3 optimization results

echo "Starting Ollama with performance optimizations..."

# Set optimization environment variables
export OLLAMA_KEEP_ALIVE=30m
export OLLAMA_NUM_PARALLEL=2
export OLLAMA_FLASH_ATTENTION=1
export OLLAMA_KV_CACHE_TYPE=q8_0
export OLLAMA_MAX_LOADED_MODELS=2
export OLLAMA_HOST=127.0.0.1:11434
export OLLAMA_SCHED_SPREAD=1

echo "Environment variables set:"
echo "OLLAMA_KEEP_ALIVE=$OLLAMA_KEEP_ALIVE"
echo "OLLAMA_NUM_PARALLEL=$OLLAMA_NUM_PARALLEL"
echo "OLLAMA_FLASH_ATTENTION=$OLLAMA_FLASH_ATTENTION"
echo "OLLAMA_KV_CACHE_TYPE=$OLLAMA_KV_CACHE_TYPE"
echo "OLLAMA_MAX_LOADED_MODELS=$OLLAMA_MAX_LOADED_MODELS"
echo "OLLAMA_HOST=$OLLAMA_HOST"
echo "OLLAMA_SCHED_SPREAD=$OLLAMA_SCHED_SPREAD"

# Check if required models are available
echo "Checking for required models..."
if ! ollama list | grep -q "llama3.1:8b"; then
    echo "Pulling llama3.1:8b model..."
    ollama pull llama3.1:8b
fi

if ! ollama list | grep -q "nomic-embed-text"; then
    echo "Pulling nomic-embed-text model..."
    ollama pull nomic-embed-text
fi

echo "Starting Ollama server..."
ollama serve
