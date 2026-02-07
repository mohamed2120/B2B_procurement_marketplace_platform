#!/bin/bash

# Update container versions in .service-versions.json after build

VERSIONS_FILE=".service-versions.json"
COMPOSE_FILE="docker-compose.all.yml"
GET_VERSION="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/get-service-version.sh"

# Initialize versions file if it doesn't exist
if [ ! -f "$VERSIONS_FILE" ]; then
    echo "{}" > "$VERSIONS_FILE"
fi

# Get all services
SERVICES=$(docker compose -f "$COMPOSE_FILE" config --services 2>/dev/null)

if [ -z "$SERVICES" ]; then
    echo "Error: Could not read services from $COMPOSE_FILE"
    exit 1
fi

# Update versions for each service
for SERVICE in $SERVICES; do
    # Skip infrastructure services
    if [[ "$SERVICE" =~ ^(postgres|redis|minio|opensearch|search-indexer-service)$ ]]; then
        continue
    fi
    
    # Get current source version
    SOURCE_VERSION=$("$GET_VERSION" "$SERVICE" 2>/dev/null)
    
    if [ -z "$SOURCE_VERSION" ]; then
        continue
    fi
    
    # Get container ID
    CONTAINER_ID=$(docker compose -f "$COMPOSE_FILE" ps -q "$SERVICE" 2>/dev/null | head -1)
    
    if [ -n "$CONTAINER_ID" ]; then
        # Try to get from Docker label
        CONTAINER_VERSION=$(docker inspect "$CONTAINER_ID" \
            --format='{{index .Config.Labels "service.version"}}' 2>/dev/null)
        
        # If no label, use source version
        if [ -z "$CONTAINER_VERSION" ] || [ "$CONTAINER_VERSION" = "<no value>" ]; then
            CONTAINER_VERSION="$SOURCE_VERSION"
        fi
    else
        CONTAINER_VERSION="$SOURCE_VERSION"
    fi
    
    # Update versions file
    if command -v jq &> /dev/null; then
        TEMP=$(mktemp)
        jq ".\"${SERVICE}\".source = \"${SOURCE_VERSION}\" | .\"${SERVICE}\".container = \"${CONTAINER_VERSION}\"" \
            "$VERSIONS_FILE" > "$TEMP" && mv "$TEMP" "$VERSIONS_FILE"
    fi
done

echo "✅ Updated container versions"
