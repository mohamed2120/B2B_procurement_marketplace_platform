#!/bin/bash

# Get version for a service
# Usage: get-service-version.sh <service-name>

SERVICE=$1

if [ -z "$SERVICE" ]; then
    echo "Usage: $0 <service-name>"
    exit 1
fi

SERVICE_DIR="services/${SERVICE}"

# Method 1: Check for VERSION file
if [ -f "${SERVICE_DIR}/VERSION" ]; then
    cat "${SERVICE_DIR}/VERSION" | tr -d ' \n'
    exit 0
fi

# Method 2: Git commit hash of service directory
if [ -d "$SERVICE_DIR" ] && command -v git &> /dev/null; then
    HASH=$(git log -1 --format=%H -- "$SERVICE_DIR" 2>/dev/null | head -c 12)
    if [ -n "$HASH" ]; then
        echo "$HASH"
        exit 0
    fi
fi

# Method 3: Hash of main.go and Dockerfile
if [ -f "${SERVICE_DIR}/cmd/server/main.go" ] || [ -f "${SERVICE_DIR}/Dockerfile" ]; then
    FILES=""
    [ -f "${SERVICE_DIR}/cmd/server/main.go" ] && FILES="${FILES} ${SERVICE_DIR}/cmd/server/main.go"
    [ -f "${SERVICE_DIR}/Dockerfile" ] && FILES="${FILES} ${SERVICE_DIR}/Dockerfile"
    
    if command -v shasum &> /dev/null; then
        HASH=$(cat $FILES 2>/dev/null | shasum -a 256 | cut -d' ' -f1 | head -c 12)
    elif command -v sha256sum &> /dev/null; then
        HASH=$(cat $FILES 2>/dev/null | sha256sum | cut -d' ' -f1 | head -c 12)
    else
        # Fallback: use modification time
        if [ -f "${SERVICE_DIR}/cmd/server/main.go" ]; then
            stat -f "%m" "${SERVICE_DIR}/cmd/server/main.go" 2>/dev/null || stat -c "%Y" "${SERVICE_DIR}/cmd/server/main.go" 2>/dev/null
            exit 0
        fi
    fi
    
    if [ -n "$HASH" ]; then
        echo "$HASH"
        exit 0
    fi
fi

# Method 4: Last modification time of service directory
if [ -d "$SERVICE_DIR" ]; then
    if command -v find &> /dev/null; then
        find "$SERVICE_DIR" -type f -name "*.go" -o -name "Dockerfile" 2>/dev/null | \
            xargs stat -f "%m" 2>/dev/null | sort -n | tail -1 || \
            find "$SERVICE_DIR" -type f -name "*.go" -o -name "Dockerfile" 2>/dev/null | \
            xargs stat -c "%Y" 2>/dev/null | sort -n | tail -1
        exit 0
    fi
fi

# Fallback: current timestamp
date +%s
