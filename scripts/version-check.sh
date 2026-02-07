#!/bin/bash

# Version checking script for services
# Tracks service versions and determines if rebuild is needed

VERSIONS_FILE=".service-versions.json"
COMPOSE_FILE="docker-compose.all.yml"

# Function to get service version
get_service_version() {
    local service=$1
    
    # Try multiple methods to get version
    # 1. Check if service has a version file
    if [ -f "services/${service}/VERSION" ]; then
        cat "services/${service}/VERSION" | tr -d ' \n'
        return 0
    fi
    
    # 2. Check git commit hash of service directory
    if [ -d "services/${service}" ] && command -v git &> /dev/null; then
        git log -1 --format=%H -- "services/${service}" 2>/dev/null | head -c 8
        return 0
    fi
    
    # 3. Check last modified time of main.go or Dockerfile
    if [ -f "services/${service}/cmd/server/main.go" ]; then
        stat -f "%m" "services/${service}/cmd/server/main.go" 2>/dev/null || stat -c "%Y" "services/${service}/cmd/server/main.go" 2>/dev/null
        return 0
    fi
    
    # 4. Fallback: use current timestamp
    date +%s
}

# Function to get running container version
get_container_version() {
    local service=$1
    
    # Check Docker image label
    local version=$(docker inspect "$(docker compose -f "$COMPOSE_FILE" ps -q "$service" 2>/dev/null | head -1)" \
        --format='{{index .Config.Labels "service.version"}}' 2>/dev/null)
    
    if [ -n "$version" ] && [ "$version" != "<no value>" ]; then
        echo "$version"
        return 0
    fi
    
    # Check from stored versions file
    if [ -f "$VERSIONS_FILE" ]; then
        jq -r ".\"${service}\".container // empty" "$VERSIONS_FILE" 2>/dev/null
    fi
}

# Function to save service version
save_service_version() {
    local service=$1
    local version=$2
    local type=$3  # "source" or "container"
    
    # Create versions file if it doesn't exist
    if [ ! -f "$VERSIONS_FILE" ]; then
        echo "{}" > "$VERSIONS_FILE"
    fi
    
    # Update version using jq
    if command -v jq &> /dev/null; then
        local temp=$(mktemp)
        jq ".\"${service}\".${type} = \"${version}\"" "$VERSIONS_FILE" > "$temp" && mv "$temp" "$VERSIONS_FILE"
    else
        # Fallback: simple JSON manipulation (basic)
        echo "Warning: jq not installed, version tracking limited"
    fi
}

# Function to check if service needs rebuild
service_needs_rebuild() {
    local service=$1
    local current_version=$(get_service_version "$service")
    local container_version=$(get_container_version "$service")
    
    if [ -z "$container_version" ]; then
        echo "new"  # No container exists, needs build
        return 0
    fi
    
    if [ "$current_version" != "$container_version" ]; then
        echo "changed"  # Version changed, needs rebuild
        return 0
    fi
    
    echo "same"  # Version same, no rebuild needed
    return 0
}

# Main function: check all services
check_all_services() {
    local services=$(docker compose -f "$COMPOSE_FILE" config --services 2>/dev/null)
    local needs_rebuild=()
    local no_rebuild=()
    local new_services=()
    
    echo "Checking service versions..."
    echo ""
    
    for service in $services; do
        # Skip infrastructure services
        if [[ "$service" =~ ^(postgres|redis|minio|opensearch)$ ]]; then
            continue
        fi
        
        local status=$(service_needs_rebuild "$service")
        local current_version=$(get_service_version "$service")
        local container_version=$(get_container_version "$service")
        
        case "$status" in
            "new")
                echo "  🔨 $service: NEW (no container) - version: ${current_version:0:8}"
                new_services+=("$service")
                ;;
            "changed")
                echo "  🔄 $service: CHANGED (${container_version:0:8} → ${current_version:0:8})"
                needs_rebuild+=("$service")
                ;;
            "same")
                echo "  ✅ $service: SAME (${current_version:0:8})"
                no_rebuild+=("$service")
                ;;
        esac
    done
    
    echo ""
    echo "Summary:"
    echo "  New services: ${#new_services[@]}"
    echo "  Changed services: ${#needs_rebuild[@]}"
    echo "  Unchanged services: ${#no_rebuild[@]}"
    echo ""
    
    if [ ${#needs_rebuild[@]} -eq 0 ] && [ ${#new_services[@]} -eq 0 ]; then
        echo "✅ All services are up to date. No rebuild needed."
        return 1  # No rebuild needed
    else
        echo "🔨 Services that need rebuild:"
        for svc in "${new_services[@]}" "${needs_rebuild[@]}"; do
            echo "  - $svc"
        done
        return 0  # Rebuild needed
    fi
}

# Export functions for use in other scripts
export -f get_service_version
export -f get_container_version
export -f save_service_version
export -f service_needs_rebuild

# If run directly, execute check
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    check_all_services
fi
