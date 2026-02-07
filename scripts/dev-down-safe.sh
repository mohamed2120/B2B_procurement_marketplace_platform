#!/bin/bash

# Safe dev-down script with automatic Docker recovery

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"

echo "=========================================="
echo "Safe Development Shutdown"
echo "=========================================="
echo ""

# Try to stop services
echo "Attempting to stop services..."
ERROR_OUTPUT=$(docker compose -f docker-compose.all.yml down 2>&1) || {
    EXIT_CODE=$?
    
    # Check if the error is Docker-related

    if echo "$ERROR_OUTPUT" | grep -qE "ECONNREFUSED.*com.docker.docker|Cannot connect to the Docker daemon.*docker.sock|docker.sock.*connect"; then
        echo "⚠️  Docker connection error detected. Attempting recovery..."
        bash "$SCRIPT_DIR/docker-recover-macos.sh" || {
            echo "⚠️  Docker recovery failed, but continuing with force stop..."
        }
        echo ""
        
        # Retry once after recovery
        echo "Retrying service stop..."
        if docker compose -f docker-compose.all.yml down 2>&1; then
            echo "✅ Services stopped successfully after recovery"
            exit 0
        else
            echo "⚠️  Still unable to stop services. Docker may need manual intervention."
            echo "You may need to manually stop containers or restart Docker Desktop."
            exit 1
        fi
    else
        # Non-Docker error, just report it
        echo "❌ Failed to stop services:"
        echo "$ERROR_OUTPUT"
        exit $EXIT_CODE
    fi
}

# Success case
echo "✅ Services stopped successfully"
exit 0
