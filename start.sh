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
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running. Please start Docker Desktop first.${NC}"
    exit 1
fi

# Check if Make is available
if ! command -v make &> /dev/null; then
    echo -e "${YELLOW}Warning: Make is not installed. Some commands may fail.${NC}"
fi

echo -e "${GREEN}Step 1: Starting all services (infrastructure + backend + frontend)...${NC}"
make up-all

echo ""
echo -e "${GREEN}Step 2: Waiting for services to be ready...${NC}"
sleep 15

echo ""
echo -e "${GREEN}Step 3: Running database migrations...${NC}"
if make migrate-all 2>&1 | grep -q "already exists"; then
    echo -e "${YELLOW}Some migrations already applied (this is normal).${NC}"
else
    echo -e "${GREEN}Migrations completed.${NC}"
fi

echo ""
echo -e "${GREEN}Step 4: Seeding database with demo data...${NC}"
cd services/identity-service && go run cmd/seed/main.go && cd ../..
echo -e "${GREEN}Identity service seeded.${NC}"

echo ""
echo -e "${GREEN}Step 5: Checking service health...${NC}"
make health-check

echo ""
echo "=========================================="
echo -e "${GREEN}System is ready!${NC}"
echo "=========================================="
echo ""
echo "Access the application:"
echo "  Frontend:  http://localhost:3000"
echo "  API Docs:  http://localhost:8001/health"
echo ""
echo "Demo Accounts (password: demo123456):"
echo "  - Platform Admin: admin@demo.com"
echo "  - Requester: buyer.requester@demo.com"
echo "  - Procurement: buyer.procurement@demo.com"
echo "  - Supplier: supplier@demo.com"
echo ""
echo "Useful commands:"
echo "  make logs-all      - View all service logs"
echo "  make health-check  - Check service health"
echo "  make down-all      - Stop all services"
echo ""
