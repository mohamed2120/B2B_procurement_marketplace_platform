#!/bin/bash

# Docker Desktop recovery script for macOS
# Automatically restarts Docker Desktop when it becomes unhealthy

set -e

# Check if running on macOS
if [[ "$(uname)" != "Darwin" ]]; then
    echo "This script is for macOS only. Skipping Docker recovery."
    exit 0
fi

# Function to check if Docker is healthy
check_docker_healthy() {
    docker info > /dev/null 2>&1 && docker ps > /dev/null 2>&1
}

# Function to wait for Docker to become healthy
wait_for_docker_healthy() {
    local max_wait=${1:-180}
    local poll_interval=${2:-2}
    local waited=0
    
    echo "Waiting for Docker to become healthy (max ${max_wait}s)..."
    
    while [ $waited -lt $max_wait ]; do
        if check_docker_healthy; then
            echo "✅ Docker is now healthy!"
            return 0
        fi
        
        if [ $((waited % 10)) -eq 0 ] && [ $waited -gt 0 ]; then
            echo "  Still waiting... (${waited}s / ${max_wait}s)"
        fi
        
        sleep $poll_interval
        waited=$((waited + poll_interval))
    done
    
    return 1
}

# Check if Docker is already healthy
if check_docker_healthy; then
    echo "✅ Docker is healthy. No recovery needed."
    exit 0
fi

echo "⚠️  Docker is unhealthy. Attempting recovery..."
echo ""

# Recovery Step A: Soft restart
echo "Step A: Soft restart of Docker Desktop..."
echo "  Quitting Docker Desktop application..."
osascript -e 'tell application "Docker" to quit' 2>/dev/null || true
sleep 3

echo "  Killing Docker Desktop processes..."
pkill -f "Docker Desktop" 2>/dev/null || true
pkill -f "com.docker.backend" 2>/dev/null || true
pkill -f "com.docker.supervisor" 2>/dev/null || true
sleep 2

echo "  Starting Docker Desktop..."
open -a Docker 2>/dev/null || {
    echo "❌ Failed to start Docker Desktop application"
    exit 1
}

echo "  Waiting for Docker to become healthy..."
if wait_for_docker_healthy 180 2; then
    echo "✅ Docker recovered successfully (soft restart)"
    exit 0
fi

echo ""
echo "⚠️  Soft restart failed. Attempting hard restart..."
echo ""

# Recovery Step B: Hard restart
echo "Step B: Hard restart of Docker Desktop..."
echo "  Force killing all Docker processes..."
killall Docker 2>/dev/null || true
pkill -f com.docker 2>/dev/null || true
sleep 3

echo "  Starting Docker Desktop..."
open -a Docker 2>/dev/null || {
    echo "❌ Failed to start Docker Desktop application"
    exit 1
}

echo "  Waiting for Docker to become healthy..."
if wait_for_docker_healthy 180 2; then
    echo "✅ Docker recovered successfully (hard restart)"
    exit 0
fi

# If still failing
echo ""
echo "❌ Docker Desktop failed to recover after soft and hard restart attempts."
echo ""
echo "Please try one of the following:"
echo "  1. Restart your Mac"
echo "  2. Manually restart Docker Desktop from Applications"
echo "  3. Reinstall Docker Desktop"
echo ""
echo "To manually restart Docker Desktop:"
echo "  open -a Docker"
echo ""
exit 1
