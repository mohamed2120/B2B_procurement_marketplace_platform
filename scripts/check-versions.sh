#!/bin/bash

# Check service versions and determine what needs rebuilding
# Usage: check-versions.sh [--rebuild-list]

VERSIONS_FILE=".service-versions.json"
COMPOSE_FILE="docker-compose.all.yml"
REBUILD_LIST=false

if [ "$1" = "--rebuild-list" ]; then
    REBUILD_LIST=true
fi

# Source the version getter
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GET_VERSION="$SCRIPT_DIR/get-service-version.sh"

# Initialize versions file if it doesn't exist
if [ ! -f "$VERSIONS_FILE" ]; then
    echo "{}" > "$VERSIONS_FILE"
fi

# Get all services from docker-compose
SERVICES=$(docker compose -f "$COMPOSE_FILE" config --services 2>/dev/null)

if [ -z "$SERVICES" ]; then
    echo "Error: Could not read services from $COMPOSE_FILE"
    exit 1
fi

NEEDS_REBUILD=()
NO_REBUILD=()
NEW_SERVICES=()

if [ "$REBUILD_LIST" != true ]; then
    echo "Checking service versions..."
    echo ""
fi

for SERVICE in $SERVICES; do
    # Skip infrastructure services (they use pre-built images)
    if [[ "$SERVICE" =~ ^(postgres|redis|minio|opensearch|search-indexer-service)$ ]]; then
        continue
    fi
    
    # Get current source version
    CURRENT_VERSION=$("$GET_VERSION" "$SERVICE" 2>/dev/null)
    
    if [ -z "$CURRENT_VERSION" ]; then
        CURRENT_VERSION="unknown"
    fi
    
    # Get container version (from running container or stored version)
    CONTAINER_VERSION=""
    CONTAINER_ID=$(docker compose -f "$COMPOSE_FILE" ps -q "$SERVICE" 2>/dev/null | head -1)
    
    if [ -n "$CONTAINER_ID" ]; then
        # Try to get from Docker label
        CONTAINER_VERSION=$(docker inspect "$CONTAINER_ID" \
            --format='{{index .Config.Labels "service.version"}}' 2>/dev/null)
        
        if [ -z "$CONTAINER_VERSION" ] || [ "$CONTAINER_VERSION" = "<no value>" ]; then
            # Try to get from stored versions file
            if command -v jq &> /dev/null && [ -f "$VERSIONS_FILE" ]; then
                CONTAINER_VERSION=$(jq -r ".\"${SERVICE}\".container // empty" "$VERSIONS_FILE" 2>/dev/null)
            fi
        fi
    else
        # No container exists, check stored version
        if command -v jq &> /dev/null && [ -f "$VERSIONS_FILE" ]; then
            CONTAINER_VERSION=$(jq -r ".\"${SERVICE}\".container // empty" "$VERSIONS_FILE" 2>/dev/null)
        fi
    fi
    
    # Determine status
    if [ -z "$CONTAINER_VERSION" ] || [ "$CONTAINER_VERSION" = "null" ] || [ "$CONTAINER_VERSION" = "" ]; then
        STATUS="new"
        NEW_SERVICES+=("$SERVICE")
        if [ "$REBUILD_LIST" != true ]; then
            echo "  🔨 $SERVICE: NEW (no version found)"
        fi
    elif [ "$CURRENT_VERSION" != "$CONTAINER_VERSION" ]; then
        STATUS="changed"
        NEEDS_REBUILD+=("$SERVICE")
        if [ "$REBUILD_LIST" != true ]; then
            echo "  🔄 $SERVICE: CHANGED (${CONTAINER_VERSION:0:8} → ${CURRENT_VERSION:0:8})"
        fi
    else
        STATUS="same"
        NO_REBUILD+=("$SERVICE")
        if [ "$REBUILD_LIST" != true ]; then
            echo "  ✅ $SERVICE: SAME (${CURRENT_VERSION:0:8})"
        fi
    fi
    
    # Update source version in file
    if command -v jq &> /dev/null; then
        TEMP=$(mktemp)
        jq ".\"${SERVICE}\".source = \"${CURRENT_VERSION}\"" "$VERSIONS_FILE" > "$TEMP" && mv "$TEMP" "$VERSIONS_FILE"
    fi
done

if [ "$REBUILD_LIST" != true ]; then
    echo ""
    echo "Summary:"
    echo "  New services: ${#NEW_SERVICES[@]}"
    echo "  Changed services: ${#NEEDS_REBUILD[@]}"
    echo "  Unchanged services: ${#NO_REBUILD[@]}"
    echo ""
fi

if [ "$REBUILD_LIST" = true ]; then
    # Output list of services that need rebuild (one per line, no other output)
    ALL_REBUILD=("${NEW_SERVICES[@]}" "${NEEDS_REBUILD[@]}")
    if [ ${#ALL_REBUILD[@]} -gt 0 ]; then
        for svc in "${ALL_REBUILD[@]}"; do
            echo "$svc"
        done
        exit 0
    else
        exit 1
    fi
fi

if [ "$REBUILD_LIST" != true ]; then
    if [ ${#NEEDS_REBUILD[@]} -eq 0 ] && [ ${#NEW_SERVICES[@]} -eq 0 ]; then
        echo "✅ All services are up to date. No rebuild needed."
        exit 1  # No rebuild needed
    else
        echo "🔨 Services that need rebuild:"
        for svc in "${NEW_SERVICES[@]}" "${NEEDS_REBUILD[@]}"; do
            echo "  - $svc"
        done
        exit 0  # Rebuild needed
    fi
fi
