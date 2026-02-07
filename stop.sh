#!/bin/bash

# B2B Procurement Marketplace Platform - Stop Script
# This script stops all services and cleans up resources

set -e

echo "=========================================="
echo "B2B Procurement Marketplace Platform"
echo "Stopping System..."
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if Docker is running
echo -e "${YELLOW}Checking Docker connection...${NC}"
if ! docker info > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Docker is not accessible. Some cleanup steps may be skipped.${NC}"
    echo ""
else
    echo -e "${GREEN}✅ Docker is accessible${NC}"
    echo ""
fi

# Stop local frontend dev server if running
echo -e "${YELLOW}Checking for local frontend dev server...${NC}"
LOCAL_PID=$(lsof -nP -iTCP:3002 -sTCP:LISTEN -t 2>/dev/null | head -1)
if [ -n "$LOCAL_PID" ]; then
    LOCAL_CMD=$(ps -p $LOCAL_PID -o command= 2>/dev/null | head -1)
    if echo "$LOCAL_CMD" | grep -qE "next dev|node.*3002|npm.*dev"; then
        echo -e "${YELLOW}Found local frontend dev server (PID $LOCAL_PID). Stopping...${NC}"
        kill -9 $LOCAL_PID 2>/dev/null || true
        sleep 1
        echo -e "${GREEN}✅ Local dev server stopped${NC}"
    else
        echo -e "${YELLOW}Port 3002 is in use by non-Node.js process (PID $LOCAL_PID)${NC}"
        echo -e "${YELLOW}   Skipping local dev server stop${NC}"
    fi
else
    echo -e "${GREEN}✅ No local frontend dev server running${NC}"
fi
echo ""

# Stop all Docker containers
echo -e "${YELLOW}Stopping all Docker containers...${NC}"
if docker info > /dev/null 2>&1; then
    # Check if containers are running
    RUNNING_COUNT=$(docker compose -f docker-compose.all.yml ps --format json 2>/dev/null | grep -c '"State":"running"' || echo "0")
    
    if [ "$RUNNING_COUNT" -gt 0 ]; then
        echo -e "${YELLOW}Found $RUNNING_COUNT running container(s)${NC}"
        
        # Try graceful stop first
        if docker compose -f docker-compose.all.yml stop 2>&1; then
            echo -e "${GREEN}✅ Containers stopped gracefully${NC}"
        else
            # If graceful stop fails, try force stop
            echo -e "${YELLOW}Graceful stop failed. Attempting force stop...${NC}"
            docker compose -f docker-compose.all.yml down --remove-orphans 2>&1 || {
                echo -e "${YELLOW}⚠️  Some containers may still be running${NC}"
            }
        fi
        
        # Verify containers are stopped
        sleep 2
        STILL_RUNNING=$(docker compose -f docker-compose.all.yml ps --format json 2>/dev/null | grep -c '"State":"running"' || echo "0")
        if [ "$STILL_RUNNING" -gt 0 ]; then
            echo -e "${YELLOW}⚠️  $STILL_RUNNING container(s) still running. Force stopping...${NC}"
            docker compose -f docker-compose.all.yml down --remove-orphans --timeout 5 2>&1 || true
        fi
    else
        echo -e "${GREEN}✅ No containers are running${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Docker is not accessible. Cannot stop containers.${NC}"
fi
echo ""

# Optional: Remove containers (commented out by default)
# Uncomment the following lines if you want to remove containers on stop
# echo -e "${YELLOW}Removing containers...${NC}"
# docker compose -f docker-compose.all.yml down --remove-orphans 2>&1 || true
# echo -e "${GREEN}✅ Containers removed${NC}"
# echo ""

# Free port 3002
echo -e "${YELLOW}Verifying port 3002 is free...${NC}"
if lsof -nP -iTCP:3002 -sTCP:LISTEN > /dev/null 2>&1; then
    PIDS=$(lsof -nP -iTCP:3002 -sTCP:LISTEN -t 2>/dev/null)
    if [ -n "$PIDS" ]; then
        echo -e "${YELLOW}Port 3002 is still in use. Cleaning up...${NC}"
        for pid in $PIDS; do
            CMD=$(ps -p $pid -o comm= 2>/dev/null | xargs basename 2>/dev/null)
            if [[ "$CMD" == "node" ]] || [[ "$CMD" == "next" ]] || [[ "$CMD" == *"node"* ]] || [[ "$CMD" == "npm" ]] || [[ "$CMD" == "yarn" ]]; then
                echo -e "${YELLOW}Killing Node.js process (PID $pid)...${NC}"
                kill -9 $pid 2>/dev/null || true
            fi
        done
        sleep 1
    fi
fi

# Verify port is free
if lsof -nP -iTCP:3002 -sTCP:LISTEN > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Port 3002 is still in use${NC}"
    echo -e "${YELLOW}   Run manually: lsof -nP -iTCP:3002 -sTCP:LISTEN${NC}"
else
    echo -e "${GREEN}✅ Port 3002 is free${NC}"
fi
echo ""

# Summary
echo "=========================================="
echo -e "${GREEN}System stopped successfully!${NC}"
echo "=========================================="
echo ""
echo "All services have been stopped."
echo ""
echo "To start again, run:"
echo "  ./start.sh"
echo ""
echo "To remove containers and volumes, run:"
echo "  docker compose -f docker-compose.all.yml down -v"
echo ""
