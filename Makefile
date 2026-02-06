.PHONY: help dev-up dev-down migrate-all seed-all test test-integration run-identity run-company run-catalog run-equipment run-marketplace run-procurement run-logistics run-collaboration run-notification run-billing run-virtual-warehouse run-search-indexer clean up-all down-all logs-all reset-all health-check

help:
	@echo "Available commands:"
	@echo "  make up-all              - Start ALL services (infra + backend + frontend) with docker-compose"
	@echo "  make down-all            - Stop all services"
	@echo "  make logs-all            - View logs from all services"
	@echo "  make reset-all           - Stop, remove volumes, and restart all services"
	@echo "  make health-check        - Check health of all backend services"
	@echo "  make dev-up              - Start infrastructure services only"
	@echo "  make dev-down            - Stop infrastructure services"
	@echo "  make migrate-all         - Run all database migrations"
	@echo "  make seed-all            - Seed all databases with test data"
	@echo "  make test                - Run all unit tests"
	@echo "  make test-integration    - Run integration tests (requires services running)"
	@echo "  make run-<service>       - Run a specific service locally"
	@echo "  make clean               - Clean up generated files"

up-all:
	@echo "Starting ALL services (infrastructure + backend + frontend)..."
	@docker compose -f docker-compose.all.yml up -d --build
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Services started. Run 'make migrate-all' to set up databases, then 'make seed-all' for test data."
	@echo "Frontend: http://localhost:3000"
	@echo "Check status: docker compose -f docker-compose.all.yml ps"

down-all:
	@echo "Stopping all services..."
	@docker compose -f docker-compose.all.yml down

logs-all:
	@docker compose -f docker-compose.all.yml logs -f

reset-all:
	@echo "Stopping all services and removing volumes..."
	@docker compose -f docker-compose.all.yml down -v
	@echo "Starting fresh..."
	@docker compose -f docker-compose.all.yml up -d --build
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Services restarted. Run 'make migrate-all' to set up databases."

health-check:
	@echo "Checking health of all backend services..."
	@echo ""
	@curl -f -s http://localhost:8001/health > /dev/null 2>&1 && echo "✅ identity-service (port 8001): OK" || echo "❌ identity-service (port 8001): FAIL"
	@curl -f -s http://localhost:8002/health > /dev/null 2>&1 && echo "✅ company-service (port 8002): OK" || echo "❌ company-service (port 8002): FAIL"
	@curl -f -s http://localhost:8003/health > /dev/null 2>&1 && echo "✅ catalog-service (port 8003): OK" || echo "❌ catalog-service (port 8003): FAIL"
	@curl -f -s http://localhost:8004/health > /dev/null 2>&1 && echo "✅ equipment-service (port 8004): OK" || echo "❌ equipment-service (port 8004): FAIL"
	@curl -f -s http://localhost:8005/health > /dev/null 2>&1 && echo "✅ marketplace-service (port 8005): OK" || echo "❌ marketplace-service (port 8005): FAIL"
	@curl -f -s http://localhost:8006/health > /dev/null 2>&1 && echo "✅ procurement-service (port 8006): OK" || echo "❌ procurement-service (port 8006): FAIL"
	@curl -f -s http://localhost:8007/health > /dev/null 2>&1 && echo "✅ logistics-service (port 8007): OK" || echo "❌ logistics-service (port 8007): FAIL"
	@curl -f -s http://localhost:8008/health > /dev/null 2>&1 && echo "✅ collaboration-service (port 8008): OK" || echo "❌ collaboration-service (port 8008): FAIL"
	@curl -f -s http://localhost:8009/health > /dev/null 2>&1 && echo "✅ notification-service (port 8009): OK" || echo "❌ notification-service (port 8009): FAIL"
	@curl -f -s http://localhost:8010/health > /dev/null 2>&1 && echo "✅ billing-service (port 8010): OK" || echo "❌ billing-service (port 8010): FAIL"
	@curl -f -s http://localhost:8011/health > /dev/null 2>&1 && echo "✅ virtual-warehouse-service (port 8011): OK" || echo "❌ virtual-warehouse-service (port 8011): FAIL"
	@curl -f -s http://localhost:8012/health > /dev/null 2>&1 && echo "✅ search-indexer-service (port 8012): OK" || echo "❌ search-indexer-service (port 8012): FAIL"
	@echo ""
	@echo "Frontend check:"
	@curl -f -s http://localhost:3000 > /dev/null 2>&1 && echo "✅ frontend (port 3000): OK" || echo "❌ frontend (port 3000): FAIL"

dev-up:
	@echo "Starting all services..."
	docker-compose -f deployments/docker-compose.yml up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Services started. Run 'make migrate-all' to set up databases."

dev-down:
	@echo "Stopping all services..."
	docker-compose -f deployments/docker-compose.yml down

migrate-all:
	@echo "Running migrations for all services..."
	@cd services/identity-service && go run cmd/migrate/main.go
	@cd services/company-service && go run cmd/migrate/main.go
	@cd services/catalog-service && go run cmd/migrate/main.go
	@cd services/equipment-service && go run cmd/migrate/main.go
	@cd services/marketplace-service && go run cmd/migrate/main.go
	@cd services/procurement-service && go run cmd/migrate/main.go
	@cd services/logistics-service && go run cmd/migrate/main.go
	@cd services/collaboration-service && go run cmd/migrate/main.go
	@cd services/notification-service && go run cmd/migrate/main.go
	@cd services/billing-service && go run cmd/migrate/main.go
	@cd services/virtual-warehouse-service && go run cmd/migrate/main.go
	@echo "All migrations completed."

seed-all:
	@echo "Seeding all databases..."
	@cd services/identity-service && go run cmd/seed/main.go
	@cd services/company-service && go run cmd/seed/main.go
	@cd services/catalog-service && go run cmd/seed/main.go
	@cd services/equipment-service && go run cmd/seed/main.go
	@cd services/marketplace-service && go run cmd/seed/main.go
	@cd services/procurement-service && go run cmd/seed/main.go
	@cd services/logistics-service && go run cmd/seed/main.go
	@cd services/collaboration-service && go run cmd/seed/main.go
	@cd services/notification-service && go run cmd/seed/main.go
	@cd services/billing-service && go run cmd/seed/main.go
	@cd services/virtual-warehouse-service && go run cmd/seed/main.go
	@echo "All seeds completed."

test:
	@echo "Running all unit tests..."
	@cd services/identity-service && go test ./... -v
	@cd services/company-service && go test ./... -v
	@cd services/catalog-service && go test ./... -v
	@cd services/equipment-service && go test ./... -v
	@cd services/marketplace-service && go test ./... -v
	@cd services/procurement-service && go test ./... -v
	@cd services/logistics-service && go test ./... -v
	@cd services/collaboration-service && go test ./... -v
	@cd services/notification-service && go test ./... -v
	@cd services/billing-service && go test ./... -v
	@cd services/virtual-warehouse-service && go test ./... -v
	@cd services/search-indexer-service && go test ./... -v
	@echo "All unit tests completed."

test-integration:
	@echo "Running integration tests..."
	@echo "Make sure services are running: make dev-up && make migrate-all && make seed-all"
	@cd tests/integration && go test ./... -v
	@echo "Integration tests completed."

run-identity:
	@cd services/identity-service && go run cmd/server/main.go

run-company:
	@cd services/company-service && go run cmd/server/main.go

run-catalog:
	@cd services/catalog-service && go run cmd/server/main.go

run-equipment:
	@cd services/equipment-service && go run cmd/server/main.go

run-marketplace:
	@cd services/marketplace-service && go run cmd/server/main.go

run-procurement:
	@cd services/procurement-service && go run cmd/server/main.go

run-logistics:
	@cd services/logistics-service && go run cmd/server/main.go

run-collaboration:
	@cd services/collaboration-service && go run cmd/server/main.go

run-notification:
	@cd services/notification-service && go run cmd/server/main.go

run-billing:
	@cd services/billing-service && go run cmd/server/main.go

run-virtual-warehouse:
	@cd services/virtual-warehouse-service && go run cmd/server/main.go

run-search-indexer:
	@cd services/search-indexer-service && go run cmd/server/main.go

run-frontend:
	@cd frontend && npm run dev

build-frontend:
	@cd frontend && npm run build

test-frontend-e2e:
	@cd frontend && npm run test:e2e

clean:
	@echo "Cleaning up..."
	@find . -type f -name "*.test" -delete
	@find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	@echo "Cleanup completed."
