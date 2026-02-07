#!/bin/bash

# Docker readiness check - CLI-based only
# Determines readiness ONLY by whether the CLI can reach the daemon
# Does NOT check for Docker Desktop process or socket paths

check_docker_ready() {
    # Try docker info first (preferred), fallback to docker version
    if docker info > /dev/null 2>&1; then
        return 0
    elif docker version > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

wait_for_docker() {
    local max_wait=${1:-90}
    local waited=0
    
    echo "Waiting for Docker daemon to be ready (max ${max_wait}s)..."
    
    while [ $waited -lt $max_wait ]; do
        if check_docker_ready; then
            echo "✅ Docker is ready!"
            return 0
        fi
        
        # Show progress every 5 seconds
        if [ $((waited % 5)) -eq 0 ]; then
            if [ $waited -eq 0 ]; then
                echo -n "Waiting"
            else
                echo -n "."
            fi
        fi
        
        # Every 15 seconds, show status
        if [ $((waited % 15)) -eq 0 ] && [ $waited -gt 0 ]; then
            echo ""
            echo "  Still waiting... (${waited}s / ${max_wait}s)"
        fi
        
        sleep 1
        waited=$((waited + 1))
    done
    
    echo ""
    return 1
}

show_docker_diagnostics() {
    echo "Docker daemon not running or context socket missing"
    echo ""
    
    # Check if docker command exists
    if ! command -v docker &> /dev/null; then
        echo "❌ Docker CLI is not installed"
        return 1
    fi
    
    echo "Docker CLI version:"
    docker --version 2>&1 || echo "  Failed to get version"
    echo ""
    
    echo "Current Docker context:"
    CURRENT_CONTEXT=$(docker context show 2>&1)
    if [ -n "$CURRENT_CONTEXT" ] && [ "$CURRENT_CONTEXT" != "Failed to show context" ]; then
        echo "$CURRENT_CONTEXT"
    else
        echo "  (unable to determine context)"
    fi
    echo ""
    
    echo "Docker context list:"
    docker context ls 2>&1 || echo "  Failed to list contexts"
    echo ""
    
    # Check if context is desktop-linux and provide specific instructions
    if [ -n "$CURRENT_CONTEXT" ] && echo "$CURRENT_CONTEXT" | grep -q "desktop-linux"; then
        echo "⚠️  Current context is 'desktop-linux' (Docker Desktop)"
        echo ""
        
        # Check for common Docker Desktop socket paths
        SOCKET_FOUND=false
        if [ -S "$HOME/.docker/run/docker.sock" ]; then
            echo "✅ Docker Desktop socket found: $HOME/.docker/run/docker.sock"
            SOCKET_FOUND=true
        elif [ -S "/var/run/docker.sock" ]; then
            echo "✅ Docker socket found: /var/run/docker.sock"
            SOCKET_FOUND=true
        fi
        
        if [ "$SOCKET_FOUND" = false ]; then
            echo "❌ Docker Desktop socket not found"
            echo ""
            echo "Start Docker Desktop (Docker.app) then retry"
            echo ""
            echo "On macOS:"
            echo "  open -a Docker"
            echo ""
            echo "On Windows:"
            echo "  Start Docker Desktop from Start Menu"
            echo ""
            echo "Wait 10-20 seconds after starting Docker Desktop, then run this script again."
            return 1
        fi
    fi
    
    echo "Docker info error output:"
    docker info 2>&1 | head -20 || echo "  docker info failed"
    echo ""
    
    echo "Docker version error output:"
    docker version 2>&1 | head -10 || echo "  docker version failed"
    echo ""
    
    echo "Troubleshooting steps:"
    echo "  1. Ensure Docker daemon is running"
    echo "  2. Check Docker context: docker context ls"
    echo "  3. Try switching context: docker context use default"
    echo "  4. Verify Docker daemon is accessible: docker ps"
}

# Main execution
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    if check_docker_ready; then
        echo "✅ Docker is ready"
        exit 0
    else
        if wait_for_docker 90; then
            exit 0
        else
            echo "❌ Error: Cannot connect to Docker daemon after 90 seconds"
            echo ""
            show_docker_diagnostics
            exit 1
        fi
    fi
fi
