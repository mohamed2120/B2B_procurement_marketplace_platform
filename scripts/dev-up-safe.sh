#!/bin/bash

# Safe dev-up script with automatic Docker recovery

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"

echo "=========================================="
echo "Safe Development Startup"
echo "=========================================="
echo ""

# Check Docker health and recover if needed
echo "Checking Docker health..."
if ! docker info > /dev/null 2>&1 || ! docker ps > /dev/null 2>&1; then
    echo "⚠️  Docker is unhealthy. Attempting recovery..."
    bash "$SCRIPT_DIR/docker-recover-macos.sh" || {
        echo "❌ Docker recovery failed. Cannot proceed."
        exit 1
    }
    echo ""
fi

# Check for Docker compose errors and recover
check_docker_compose() {
    docker compose -f docker-compose.all.yml config > /dev/null 2>&1
}

if ! check_docker_compose; then
    ERROR_OUTPUT=$(docker compose -f docker-compose.all.yml config 2>&1 || true)
    
    if echo "$ERROR_OUTPUT" | grep -qE "ECONNREFUSED.*com.docker.docker|Cannot connect to the Docker daemon.*docker.sock"; then
        echo "⚠️  Docker compose connection error detected. Attempting recovery..."
        bash "$SCRIPT_DIR/docker-recover-macos.sh" || {
            echo "❌ Docker recovery failed. Cannot proceed."
            exit 1
        }
        echo ""
    fi
fi

# Proceed with normal dev-up
echo "Docker is healthy. Proceeding with dev-up..."
echo ""

# Call make dev-up
make dev-up
