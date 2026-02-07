#!/bin/bash

# B2B Procurement Marketplace Platform - Smoke Test Script
# This script performs basic health checks and smoke tests

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=========================================="
echo "B2B Procurement Marketplace Platform"
echo "Smoke Test"
echo "=========================================="
echo ""

FAILED=0

# Check Docker containers
echo -e "${YELLOW}1. Checking Docker containers...${NC}"
if ! docker compose -f docker-compose.all.yml ps --format json 2>/dev/null | grep -q '"State":"running"'; then
    echo -e "${RED}❌ No containers are running${NC}"
    FAILED=1
else
    RUNNING=$(docker compose -f docker-compose.all.yml ps --format json 2>/dev/null | grep -c '"State":"running"' || echo "0")
    echo -e "${GREEN}✅ $RUNNING container(s) running${NC}"
fi
echo ""

# Check backend services health
echo -e "${YELLOW}2. Checking backend service health endpoints...${NC}"
SERVICES=(
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
    "diagnostics-service:8012"
)

for service_port in "${SERVICES[@]}"; do
    service=$(echo $service_port | cut -d: -f1)
    port=$(echo $service_port | cut -d: -f2)
    
    if curl -sf "http://localhost:${port}/health" > /dev/null 2>&1; then
        echo -e "  ${GREEN}✅${NC} $service (port $port)"
    elif curl -sf "http://localhost:${port}/ready" > /dev/null 2>&1; then
        echo -e "  ${GREEN}✅${NC} $service (port $port) - /ready"
    else
        echo -e "  ${RED}❌${NC} $service (port $port) - not responding"
        FAILED=1
    fi
done
echo ""

# Check frontend
echo -e "${YELLOW}3. Checking frontend...${NC}"
if curl -sf "http://localhost:3002" > /dev/null 2>&1; then
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:3002")
    if [ "$STATUS" = "200" ] || [ "$STATUS" = "307" ]; then
        echo -e "  ${GREEN}✅${NC} Frontend (port 3002) - HTTP $STATUS"
    else
        echo -e "  ${YELLOW}⚠️${NC}  Frontend (port 3002) - HTTP $STATUS"
    fi
else
    echo -e "  ${RED}❌${NC} Frontend (port 3002) - not responding"
    FAILED=1
fi
echo ""

# Check public pages
echo -e "${YELLOW}4. Checking public pages...${NC}"
PUBLIC_PAGES=(
    "/"
    "/login"
    "/register"
    "/how-it-works"
    "/pricing"
    "/contact"
    "/terms"
    "/privacy"
)

for page in "${PUBLIC_PAGES[@]}"; do
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:3002${page}" || echo "000")
    if [ "$STATUS" = "200" ] || [ "$STATUS" = "307" ]; then
        echo -e "  ${GREEN}✅${NC} $page - HTTP $STATUS"
    else
        echo -e "  ${YELLOW}⚠️${NC}  $page - HTTP $STATUS"
    fi
done
echo ""

# Check database connectivity (via identity service)
echo -e "${YELLOW}5. Checking database connectivity...${NC}"
if curl -sf "http://localhost:8001/health" | grep -q "database" || curl -sf "http://localhost:8001/ready" > /dev/null; then
    echo -e "  ${GREEN}✅${NC} Database appears accessible (via identity-service)"
else
    echo -e "  ${YELLOW}⚠️${NC}  Database connectivity check unavailable"
fi
echo ""

# Summary
echo "=========================================="
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ Smoke test passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Smoke test failed!${NC}"
    echo ""
    echo "Troubleshooting:"
    echo "  1. Check Docker containers: docker compose -f docker-compose.all.yml ps"
    echo "  2. Check service logs: docker compose -f docker-compose.all.yml logs [service-name]"
    echo "  3. Restart services: ./start.sh"
    exit 1
fi
