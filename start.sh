#!/bin/bash

# B2B Procurement Marketplace Platform - Start Script
# This script starts all services, runs migrations, and seeds the database

set -e

echo "=========================================="
echo "B2B Procurement Marketplace Platform"
echo "Starting System..."
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if Docker is running
echo -e "${YELLOW}Checking Docker connection...${NC}"

# Source the Docker check script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "${SCRIPT_DIR}/scripts/check-docker.sh" ]; then
    source "${SCRIPT_DIR}/scripts/check-docker.sh"
else
    # Fallback if script not found
    check_docker_ready() {
        docker info > /dev/null 2>&1 || docker version > /dev/null 2>&1
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
            sleep 1
            waited=$((waited + 1))
        done
        return 1
    }
    show_docker_diagnostics() {
        echo "Docker diagnostics:"
        docker --version 2>&1 || echo "Docker CLI not found"
        docker context ls 2>&1 || echo "Failed to list contexts"
        docker context show 2>&1 || echo "Failed to show context"
        docker info 2>&1 || echo "docker info failed"
    }
fi

# Function to check if Docker daemon is accessible (for compatibility)
check_docker() {
    check_docker_ready
}

# Function to free port 3002 safely (only kills Node.js/Next.js processes)
free_port_3002() {
    local port=3002
    
    # Use the safe free-port script if available
    if [ -f "frontend/scripts/free-port.sh" ]; then
        if AUTO_KILL=true bash frontend/scripts/free-port.sh; then
            return 0
        else
            return 1
        fi
    fi
    
    # Fallback: safe manual check
    local listening_pids=$(lsof -nP -iTCP:$port -sTCP:LISTEN -t 2>/dev/null)
    
    if [ -z "$listening_pids" ]; then
        echo -e "${GREEN}✅ Port ${port} is free${NC}"
        return 0
    fi
    
    echo -e "${YELLOW}Port ${port} is in use. Checking processes...${NC}"
    echo ""
    echo -e "${YELLOW}Diagnostics:${NC}"
    echo "--- Docker containers using port $port ---"
    docker ps --format "table {{.Names}}\t{{.Ports}}" 2>/dev/null | grep -E "3002|$port" || echo "  (none)"
    echo ""
    echo "--- Processes listening on port $port ---"
    lsof -nP -iTCP:$port -sTCP:LISTEN 2>/dev/null || echo "  (none found)"
    echo ""
    
    # Check each process
    local safe_to_kill=()
    local unsafe_processes=()
    local allowed=("node" "next" "npm" "yarn" "pnpm")
    
    for pid in $listening_pids; do
        local cmd=$(ps -p $pid -o comm= 2>/dev/null | xargs basename 2>/dev/null)
        if [ -z "$cmd" ]; then
            cmd=$(ps -p $pid -o command= 2>/dev/null | awk '{print $1}' | xargs basename 2>/dev/null)
        fi
        
        if [ -z "$cmd" ]; then
            continue
        fi
        
        local is_allowed=false
        for allowed_cmd in "${allowed[@]}"; do
            if [ "$cmd" = "$allowed_cmd" ] || [[ "$cmd" == *"$allowed_cmd"* ]]; then
                is_allowed=true
                break
            fi
        done
        
        if [ "$is_allowed" = true ]; then
            safe_to_kill+=("$pid")
            echo -e "  ${GREEN}✅${NC} PID $pid ($cmd) - safe to kill"
        else
            unsafe_processes+=("$pid:$cmd")
            echo -e "  ${RED}❌${NC} PID $pid ($cmd) - NOT safe to kill"
        fi
    done
    
    echo ""
    
    # If unsafe processes found, don't proceed
    if [ ${#unsafe_processes[@]} -gt 0 ]; then
        echo -e "${RED}❌ ERROR: Port ${port} is occupied by processes that cannot be safely killed:${NC}"
        for unsafe in "${unsafe_processes[@]}"; do
            IFS=':' read -r pid cmd <<< "$unsafe"
            echo -e "   - PID $pid: $cmd"
        done
        echo ""
        echo -e "${YELLOW}Please close that application or stop the container using port ${port}.${NC}"
        echo ""
        echo -e "${YELLOW}To see what's using the port:${NC}"
        echo "  lsof -nP -iTCP:$port -sTCP:LISTEN"
        echo "  docker ps --format 'table {{.Names}}\t{{.Ports}}' | grep $port"
        return 1
    fi
    
    # Kill safe processes
    if [ ${#safe_to_kill[@]} -gt 0 ]; then
        echo -e "${YELLOW}Killing safe processes on port ${port}...${NC}"
        for pid in "${safe_to_kill[@]}"; do
            kill -9 $pid 2>/dev/null || true
        done
        sleep 1
        
        # Verify port is free
        if lsof -nP -iTCP:$port -sTCP:LISTEN > /dev/null 2>&1; then
            echo -e "${RED}❌ Failed to free port ${port}${NC}"
            return 1
        else
            echo -e "${GREEN}✅ Port ${port} is now free${NC}"
            return 0
        fi
    else
        # Check if port is still in use
        if lsof -nP -iTCP:$port -sTCP:LISTEN > /dev/null 2>&1; then
            echo -e "${RED}❌ Port ${port} is still in use but no safe processes found${NC}"
            return 1
        else
            echo -e "${GREEN}✅ Port ${port} is now free${NC}"
            return 0
        fi
    fi
}

if ! check_docker_ready; then
    echo -e "${YELLOW}Docker daemon is not accessible.${NC}"
    echo ""
    
    if wait_for_docker 90; then
        echo -e "${GREEN}✅ Docker is now ready!${NC}"
    else
        echo -e "${RED}❌ Error: Cannot connect to Docker daemon after 90 seconds.${NC}"
        echo ""
        show_docker_diagnostics
        exit 1
    fi
else
    echo -e "${GREEN}✅ Docker is ready!${NC}"
fi

# Check Docker memory/resources
echo -e "${YELLOW}Checking Docker resources...${NC}"
docker_memory=$(docker info 2>/dev/null | grep -i "Total Memory" | awk '{print $3 $4}' || echo "unknown")
if [ "$docker_memory" != "unknown" ]; then
    echo -e "  Docker Memory: ${docker_memory}"
fi

# Warn about potential memory issues
echo -e "${YELLOW}Note: If you encounter 'cannot allocate memory' errors:${NC}"
echo -e "${YELLOW}  1. Increase Docker Desktop memory (Settings > Resources > Memory)${NC}"
echo -e "${YELLOW}  2. Recommended: At least 8GB allocated to Docker${NC}"
echo -e "${YELLOW}  3. Clean build cache: docker builder prune -f${NC}"
echo ""

# Check if Make is available
if ! command -v make &> /dev/null; then
    echo -e "${YELLOW}Warning: Make is not installed. Some commands may fail.${NC}"
fi

# Check if curl is available
if ! command -v curl &> /dev/null; then
    echo -e "${YELLOW}Warning: curl is not installed. Health checks will be limited.${NC}"
    CURL_AVAILABLE=false
else
    CURL_AVAILABLE=true
fi

# Function to check if containers are already running
check_containers_running() {
    local running_count=$(docker compose -f docker-compose.all.yml ps --format json 2>/dev/null | grep -c '"State":"running"' || echo "0")
    local total_count=$(docker compose -f docker-compose.all.yml config --services 2>/dev/null | wc -l | tr -d ' ' || echo "0")
    
    if [ "$running_count" -gt 0 ] && [ "$running_count" -ge "$((total_count / 2))" ]; then
        return 0  # Most containers are running
    else
        return 1  # Containers are not running
    fi
}

# Free port 3002 before starting (frontend port)
echo -e "${YELLOW}Checking and freeing port 3002 for frontend...${NC}"

# Check if port is in use
if lsof -nP -iTCP:3002 -sTCP:LISTEN > /dev/null 2>&1; then
    echo -e "${YELLOW}Port 3002 is in use. Checking what's using it...${NC}"
    
    # Check if it's a local Node.js process (local dev server)
    LOCAL_PID=$(lsof -nP -iTCP:3002 -sTCP:LISTEN -t 2>/dev/null | head -1)
    if [ -n "$LOCAL_PID" ]; then
        LOCAL_CMD=$(ps -p $LOCAL_PID -o command= 2>/dev/null | head -1)
        if echo "$LOCAL_CMD" | grep -qE "next dev|node.*3002|npm.*dev"; then
            echo -e "${YELLOW}⚠️  Found local frontend dev server running (PID $LOCAL_PID)${NC}"
            echo -e "${YELLOW}   This will conflict with Docker frontend container!${NC}"
            echo -e "${YELLOW}   Stopping local dev server...${NC}"
            kill -9 $LOCAL_PID 2>/dev/null || true
            sleep 2
            echo -e "${GREEN}✅ Local dev server stopped${NC}"
        fi
    fi
    
    # Check if a Docker container is using the port
    DOCKER_CONTAINER=$(docker ps --format "{{.Names}}\t{{.Ports}}" 2>/dev/null | grep -E ":3002->|3002:" | awk '{print $1}' | head -1)
    if [ -n "$DOCKER_CONTAINER" ]; then
        echo -e "${YELLOW}Found Docker container '$DOCKER_CONTAINER' using port 3002${NC}"
        echo -e "${YELLOW}   (This is expected - we'll restart it)${NC}"
        # Don't stop it here - let docker compose handle it
    fi
    
    # Check if port is still in use (might be a Node.js process)
    if lsof -nP -iTCP:3002 -sTCP:LISTEN > /dev/null 2>&1; then
        echo -e "${YELLOW}Port still in use. Checking for Node.js processes...${NC}"
        
        # Try the safe free-port script if available
        if [ -f "frontend/scripts/free-port.sh" ]; then
            if AUTO_KILL=true bash frontend/scripts/free-port.sh; then
                echo -e "${GREEN}✅ Port 3002 freed successfully${NC}"
            else
                # Direct kill as fallback (only Node.js processes)
                PIDS=$(lsof -nP -iTCP:3002 -sTCP:LISTEN -t 2>/dev/null)
                if [ -n "$PIDS" ]; then
                    for pid in $PIDS; do
                        CMD=$(ps -p $pid -o comm= 2>/dev/null | xargs basename 2>/dev/null)
                        if [[ "$CMD" == "node" ]] || [[ "$CMD" == "next" ]] || [[ "$CMD" == *"node"* ]] || [[ "$CMD" == "npm" ]] || [[ "$CMD" == "yarn" ]]; then
                            echo -e "${YELLOW}Killing Node.js process (PID $pid, $CMD)...${NC}"
                            kill -9 $pid 2>/dev/null || true
                        else
                            echo -e "${YELLOW}⚠️  Process $pid ($CMD) is using port 3002 but is not a Node.js process${NC}"
                            echo -e "${YELLOW}   This may be a system process. Please check manually.${NC}"
                        fi
                    done
                    sleep 2
                fi
                
                # Verify port is free
                if lsof -nP -iTCP:3002 -sTCP:LISTEN > /dev/null 2>&1; then
                    echo -e "${RED}❌ Port 3002 is still in use after kill attempt${NC}"
                    echo -e "${YELLOW}Please manually free port 3002 and try again:${NC}"
                    echo "  lsof -nP -iTCP:3002 -sTCP:LISTEN"
                    echo "  kill -9 <PID>"
                    exit 1
                else
                    echo -e "${GREEN}✅ Port 3002 freed successfully${NC}"
                fi
            fi
        else
            # Direct kill fallback
            PIDS=$(lsof -nP -iTCP:3002 -sTCP:LISTEN -t 2>/dev/null)
            if [ -n "$PIDS" ]; then
                for pid in $PIDS; do
                    CMD=$(ps -p $pid -o comm= 2>/dev/null | xargs basename 2>/dev/null)
                    if [[ "$CMD" == "node" ]] || [[ "$CMD" == "next" ]] || [[ "$CMD" == *"node"* ]] || [[ "$CMD" == "npm" ]] || [[ "$CMD" == "yarn" ]]; then
                        echo -e "${YELLOW}Killing Node.js process (PID $pid)...${NC}"
                        kill -9 $pid 2>/dev/null || true
                    fi
                done
                sleep 2
            fi
        fi
    else
        echo -e "${GREEN}✅ Port 3002 freed after stopping Docker container${NC}"
    fi
else
    echo -e "${GREEN}✅ Port 3002 is already free${NC}"
fi
echo ""

# Restart Docker containers if they're running
echo -e "${YELLOW}Checking Docker containers...${NC}"
if check_containers_running; then
    running_services=$(docker compose -f docker-compose.all.yml ps --services --filter "status=running" 2>/dev/null | wc -l | tr -d ' ')
    echo -e "${YELLOW}Found ${running_services} running containers. Restarting them...${NC}"
    
    # Check Docker health before restart
    if ! docker info > /dev/null 2>&1 || ! docker ps > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  Docker is unhealthy. Attempting recovery...${NC}"
        if [ -f "${SCRIPT_DIR}/scripts/docker-recover-macos.sh" ] && [[ "$(uname)" == "Darwin" ]]; then
            bash "${SCRIPT_DIR}/scripts/docker-recover-macos.sh" || {
                echo -e "${RED}❌ Docker recovery failed. Cannot proceed.${NC}"
                exit 1
            }
        else
            echo -e "${RED}❌ Docker is unhealthy. Please restart Docker Desktop manually.${NC}"
            exit 1
        fi
    fi
    
    # Try to restart containers with error handling
    ERROR_OUTPUT=$(docker compose -f docker-compose.all.yml restart 2>&1) || {
        EXIT_CODE=$?
        
        # Check if error is Docker connection related
        if echo "$ERROR_OUTPUT" | grep -qE "ECONNREFUSED.*com.docker.docker|Cannot connect to the Docker daemon.*docker.sock|docker.sock.*connect"; then
            echo -e "${YELLOW}⚠️  Docker connection error detected. Attempting recovery...${NC}"
            if [ -f "${SCRIPT_DIR}/scripts/docker-recover-macos.sh" ] && [[ "$(uname)" == "Darwin" ]]; then
                bash "${SCRIPT_DIR}/scripts/docker-recover-macos.sh" || {
                    echo -e "${RED}❌ Docker recovery failed. Cannot proceed.${NC}"
                    exit 1
                }
                echo ""
                echo -e "${YELLOW}Retrying container restart...${NC}"
                docker compose -f docker-compose.all.yml restart || {
                    echo -e "${RED}❌ Failed to restart containers after recovery.${NC}"
                    exit 1
                }
            else
                echo -e "${RED}❌ Docker connection error. Please restart Docker Desktop manually.${NC}"
                exit 1
            fi
        else
            echo -e "${RED}❌ Failed to restart containers:${NC}"
            echo "$ERROR_OUTPUT"
            exit $EXIT_CODE
        fi
    }
    
    echo -e "${GREEN}✅ Containers restarted${NC}"
    echo ""
    sleep 5  # Give containers time to restart
fi

# Check service versions (smart rebuild)
echo -e "${YELLOW}Checking service versions...${NC}"

# Helper function to safely check Docker before compose operations
safe_docker_check() {
    if ! docker info > /dev/null 2>&1 || ! docker ps > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  Docker is unhealthy. Attempting recovery...${NC}"
        if [ -f "${SCRIPT_DIR}/scripts/docker-recover-macos.sh" ] && [[ "$(uname)" == "Darwin" ]]; then
            bash "${SCRIPT_DIR}/scripts/docker-recover-macos.sh" || {
                echo -e "${RED}❌ Docker recovery failed. Cannot proceed.${NC}"
                return 1
            }
        else
            echo -e "${RED}❌ Docker is unhealthy. Please restart Docker Desktop manually.${NC}"
            return 1
        fi
    fi
    return 0
}

# Helper function to safely run docker compose with recovery
safe_docker_compose() {
    local cmd="$1"
    shift
    
    local error_output=$(docker compose -f docker-compose.all.yml $cmd "$@" 2>&1) || {
        local exit_code=$?
        
        if echo "$error_output" | grep -qE "ECONNREFUSED.*com.docker.docker|Cannot connect to the Docker daemon.*docker.sock|docker.sock.*connect"; then
            echo -e "${YELLOW}⚠️  Docker connection error detected. Attempting recovery...${NC}"
            if [ -f "${SCRIPT_DIR}/scripts/docker-recover-macos.sh" ] && [[ "$(uname)" == "Darwin" ]]; then
                bash "${SCRIPT_DIR}/scripts/docker-recover-macos.sh" || {
                    echo -e "${RED}❌ Docker recovery failed. Cannot proceed.${NC}"
                    return 1
                }
                echo ""
                echo -e "${YELLOW}Retrying docker compose operation...${NC}"
                docker compose -f docker-compose.all.yml $cmd "$@" || {
                    echo -e "${RED}❌ Failed after recovery.${NC}"
                    return 1
                }
            else
                echo -e "${RED}❌ Docker connection error. Please restart Docker Desktop manually.${NC}"
                return 1
            fi
        else
            echo -e "${RED}❌ Docker compose failed:${NC}"
            echo "$error_output"
            return $exit_code
        fi
    }
    return 0
}

# Function to check if frontend needs rebuild
check_frontend_changes() {
    local frontend_dir="frontend"
    local frontend_version_file=".frontend-version"
    
    # Calculate current frontend version
    # Priority: file hash (detects uncommitted changes) > git hash > modification time
    local current_version=""
    
    # First, try file hash (detects uncommitted changes)
    if command -v shasum &> /dev/null; then
        current_version=$(find "$frontend_dir" -type f \( -name "package.json" -o -name "next.config.js" -o -name "tsconfig.json" -o -path "*/app/layout.tsx" -o -path "*/app/app/*/page.tsx" \) 2>/dev/null | \
            head -20 | xargs cat 2>/dev/null | shasum -a 256 | cut -d' ' -f1 | head -c 12)
    elif command -v sha256sum &> /dev/null; then
        current_version=$(find "$frontend_dir" -type f \( -name "package.json" -o -name "next.config.js" -o -name "tsconfig.json" -o -path "*/app/layout.tsx" -o -path "*/app/app/*/page.tsx" \) 2>/dev/null | \
            head -20 | xargs cat 2>/dev/null | sha256sum | cut -d' ' -f1 | head -c 12)
    fi
    
    # Fallback to git hash if file hash failed
    if [ -z "$current_version" ] && command -v git &> /dev/null && [ -d "$frontend_dir" ]; then
        current_version=$(git log -1 --format=%H -- "$frontend_dir" 2>/dev/null | head -c 12)
    fi
    
    # Last resort: modification time
    if [ -z "$current_version" ] && [ -f "$frontend_dir/package.json" ]; then
        current_version=$(stat -f "%m" "$frontend_dir/package.json" 2>/dev/null || stat -c "%Y" "$frontend_dir/package.json" 2>/dev/null)
    fi
    
    # Get stored version
    local stored_version=""
    if [ -f "$frontend_version_file" ]; then
        stored_version=$(cat "$frontend_version_file" 2>/dev/null | tr -d ' \n')
    fi
    
    # Check if container exists and get its version
    local container_version=""
    if docker compose -f docker-compose.all.yml ps frontend 2>/dev/null | grep -q "running\|created"; then
        container_version=$(docker inspect $(docker compose -f docker-compose.all.yml ps -q frontend 2>/dev/null) 2>/dev/null | \
            grep -o '"service.version":"[^"]*"' | cut -d'"' -f4 | head -c 12)
    fi
    
    # Compare versions
    if [ -z "$stored_version" ] && [ -z "$container_version" ]; then
        echo "new"  # No version stored, needs build
        return 0
    fi
    
    local compare_version="${container_version:-$stored_version}"
    if [ "$current_version" != "$compare_version" ]; then
        echo "changed"  # Version changed, needs rebuild
        return 0
    fi
    
    echo "same"  # Version same, no rebuild needed
    return 0
}

# Check for changes in frontend and backend
echo -e "${YELLOW}Checking for code changes (frontend and backend)...${NC}"
echo ""

FRONTEND_NEEDS_REBUILD=false
BACKEND_NEEDS_REBUILD=false

# Check frontend
FRONTEND_STATUS=$(check_frontend_changes)
if [ "$FRONTEND_STATUS" = "new" ] || [ "$FRONTEND_STATUS" = "changed" ]; then
    FRONTEND_NEEDS_REBUILD=true
    if [ "$FRONTEND_STATUS" = "new" ]; then
        echo -e "  🔨 Frontend: NEW (needs initial build)"
    else
        echo -e "  🔄 Frontend: CHANGED (needs rebuild)"
    fi
else
    echo -e "  ✅ Frontend: UP TO DATE"
fi

# Check backend services
if bash scripts/check-versions.sh > /tmp/version-check.log 2>&1; then
    # Some backend services need rebuild
    BACKEND_NEEDS_REBUILD=true
    cat /tmp/version-check.log
elif [ -s /tmp/version-check.log ]; then
    # All backend services up to date
    cat /tmp/version-check.log
else
    # Version check failed - assume rebuild needed for safety
    echo -e "  ⚠️  Backend: Version check unavailable (will rebuild for safety)"
    BACKEND_NEEDS_REBUILD=true
fi

echo ""

# Check Docker health before proceeding
if ! safe_docker_check; then
    exit 1
fi

# Rebuild logic
if [ "$FRONTEND_NEEDS_REBUILD" = true ] || [ "$BACKEND_NEEDS_REBUILD" = true ]; then
    echo -e "${GREEN}Step 1: Building and starting changed services...${NC}"
    
    # Always rebuild frontend if it needs rebuild
    if [ "$FRONTEND_NEEDS_REBUILD" = true ]; then
        echo -e "${YELLOW}Rebuilding frontend (new version detected)...${NC}"
        # Calculate version for frontend (use same method as check_frontend_changes)
        VERSION=""
        if command -v shasum &> /dev/null; then
            VERSION=$(find frontend -type f \( -name "package.json" -o -name "next.config.js" -o -name "tsconfig.json" -o -path "*/app/layout.tsx" -o -path "*/app/app/*/page.tsx" \) 2>/dev/null | \
                head -20 | xargs cat 2>/dev/null | shasum -a 256 | cut -d' ' -f1 | head -c 12)
        elif command -v sha256sum &> /dev/null; then
            VERSION=$(find frontend -type f \( -name "package.json" -o -name "next.config.js" -o -name "tsconfig.json" -o -path "*/app/layout.tsx" -o -path "*/app/app/*/page.tsx" \) 2>/dev/null | \
                head -20 | xargs cat 2>/dev/null | sha256sum | cut -d' ' -f1 | head -c 12)
        fi
        
        # Fallback to git hash if file hash failed
        if [ -z "$VERSION" ] && command -v git &> /dev/null; then
            VERSION=$(git log -1 --format=%H -- frontend 2>/dev/null | head -c 12 || echo "unknown")
        fi
        
        # Last resort: timestamp
        if [ -z "$VERSION" ]; then
            VERSION=$(date +%s)
        fi
        
        echo -e "${YELLOW}Building frontend with version: ${VERSION:0:12}...${NC}"
        SERVICE_VERSION="$VERSION" safe_docker_compose build frontend || exit 1
        safe_docker_compose up -d frontend || exit 1
        
        # Save frontend version
        echo "$VERSION" > .frontend-version
        echo -e "${GREEN}✅ Frontend rebuilt and started${NC}"
        echo ""
    fi
    
    # Rebuild backend services if needed
    if [ "$BACKEND_NEEDS_REBUILD" = true ]; then
        # Use make dev-up which handles backend version checking
        echo -e "${YELLOW}Rebuilding backend services...${NC}"
        make dev-up
    fi
else
    # All services up to date
    echo -e "${GREEN}✅ All services (frontend + backend) are up to date${NC}"
    
    # Still check if containers are running
    if check_containers_running; then
        running_services=$(docker compose -f docker-compose.all.yml ps --services --filter "status=running" 2>/dev/null | wc -l | tr -d ' ')
        echo -e "${GREEN}✅ All services up to date and running (${running_services} services)${NC}"
        echo -e "${YELLOW}Containers already restarted above. Skipping additional restart.${NC}"
    else
        echo -e "${YELLOW}Services up to date but containers not running. Starting...${NC}"
        echo ""
        echo -e "${GREEN}Step 1: Starting all services...${NC}"
        safe_docker_compose up -d || exit 1
    fi
fi

echo ""
echo -e "${GREEN}Step 2: Waiting for services to be ready...${NC}"
echo -e "${YELLOW}Waiting for infrastructure services (PostgreSQL, Redis)...${NC}"
max_attempts=30
attempt=0
while [ $attempt -lt $max_attempts ]; do
    if docker compose -f docker-compose.all.yml ps postgres redis | grep -q "healthy\|running"; then
        break
    fi
    attempt=$((attempt + 1))
    echo -n "."
    sleep 2
done
echo ""

echo -e "${YELLOW}Waiting for backend services to start...${NC}"
sleep 20

echo ""
echo -e "${GREEN}Step 3: Verifying backend services are running...${NC}"
backend_services=(
    "identity-service:8001"
    "company-service:8002"
    "catalog-service:8003"
    "equipment-service:8004"
    "marketplace-service:8005"
    "procurement-service:8006"
    "logistics-service:8007"
    "collaboration-service:8008"
    "notification-service:8009"
    "billing-service:8010"
    "virtual-warehouse-service:8011"
    "diagnostics-service:8013"
)

all_backend_ok=true
for service_port in "${backend_services[@]}"; do
    IFS=':' read -r service port <<< "$service_port"
    if [ "$CURL_AVAILABLE" = true ]; then
        if curl -f -s -X GET "http://localhost:${port}/health" > /dev/null 2>&1; then
            echo -e "  ${GREEN}✅${NC} ${service} (port ${port}): Running"
        else
            echo -e "  ${RED}❌${NC} ${service} (port ${port}): Not responding"
            all_backend_ok=false
        fi
    else
        # Fallback: check if container is running
        if docker compose -f docker-compose.all.yml ps "$service" 2>/dev/null | grep -q "running"; then
            echo -e "  ${GREEN}✅${NC} ${service} (port ${port}): Container running"
        else
            echo -e "  ${RED}❌${NC} ${service} (port ${port}): Container not running"
            all_backend_ok=false
        fi
    fi
done

if [ "$all_backend_ok" = false ]; then
    echo -e "${YELLOW}Warning: Some backend services are not responding yet. They may still be starting up.${NC}"
    echo -e "${YELLOW}You can check logs with: docker compose -f docker-compose.all.yml logs <service-name>${NC}"
fi

echo ""
echo -e "${GREEN}Step 4: Verifying frontend is running...${NC}"
if [ "$CURL_AVAILABLE" = true ]; then
    if curl -f -s "http://localhost:3002" > /dev/null 2>&1; then
        echo -e "  ${GREEN}✅${NC} Frontend (port 3002): Running"
    else
        echo -e "  ${YELLOW}⚠️${NC}  Frontend (port 3002): Not responding yet (may still be starting)"
        echo -e "  ${YELLOW}   Check logs: docker compose -f docker-compose.all.yml logs frontend${NC}"
    fi
else
    # Fallback: check if container is running
    if docker compose -f docker-compose.all.yml ps frontend | grep -q "running"; then
        echo -e "  ${GREEN}✅${NC} Frontend (port 3002): Container running"
    else
        echo -e "  ${RED}❌${NC} Frontend (port 3002): Container not running"
        echo -e "  ${YELLOW}   Check logs: docker compose -f docker-compose.all.yml logs frontend${NC}"
    fi
fi

echo ""
echo -e "${GREEN}Step 5: Running database migrations...${NC}"
if make migrate-all; then
    echo -e "${GREEN}Migrations completed.${NC}"
else
    echo -e "${YELLOW}Some migrations may have already been applied (this is normal).${NC}"
fi

echo ""
echo -e "${GREEN}Step 6: Seeding database with demo data...${NC}"
if make seed-all; then
    echo -e "${GREEN}All services seeded.${NC}"
else
    echo -e "${YELLOW}Some services may have already been seeded (this is normal).${NC}"
fi

echo ""
echo -e "${GREEN}Step 7: Final health check...${NC}"
make health-check

echo ""
echo "=========================================="
echo -e "${GREEN}System is ready!${NC}"
echo "=========================================="
echo ""
# Get the actual frontend port from docker-compose file
FRONTEND_PORT=$(sed -n '/^  frontend:/,/^  [a-z]/p' docker-compose.all.yml | grep -E '^\s+-\s*"[0-9]+:' | sed 's/.*"\([0-9]*\):.*/\1/' | head -1)
if [ -z "$FRONTEND_PORT" ]; then
    # Fallback: try to get from running container port mapping  
    FRONTEND_PORT=$(docker compose -f docker-compose.all.yml ps frontend 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+->' | head -1 | cut -d: -f2 | cut -d- -f1)
fi
if [ -z "$FRONTEND_PORT" ]; then
    FRONTEND_PORT="3002"
fi
echo "Access the application:"
echo "  Frontend:  http://localhost:${FRONTEND_PORT}"
echo "  API Docs:  http://localhost:8001/health"
echo ""
echo "Demo Accounts (password: demo123456):"
echo "  - Platform Admin: admin@demo.com"
echo "  - Requester: buyer.requester@demo.com"
echo "  - Procurement: buyer.procurement@demo.com"
echo "  - Supplier: supplier@demo.com"
echo ""
echo "Note: OpenSearch is disabled by default. Use 'make dev-up-search' to enable search features."
echo ""
echo "Useful commands:"
echo "  make logs-all        - View all service logs"
echo "  make health-check    - Check service health"
echo "  make dev-down        - Stop all services"
echo "  make dev-up-search   - Start with OpenSearch enabled"
echo ""
