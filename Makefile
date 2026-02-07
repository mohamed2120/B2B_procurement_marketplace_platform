.PHONY: help dev-up dev-up-safe dev-up-search dev-down dev-down-safe migrate-service migrate-all seed-all test test-integration run-identity run-company run-catalog run-equipment run-marketplace run-procurement run-logistics run-collaboration run-notification run-billing run-virtual-warehouse run-search-indexer clean up-all down-all logs-all reset-all health-check verify verify-logs

help:
	@echo "Available commands:"
	@echo "  make up-all              - Start ALL services (infra + backend + frontend) with docker-compose"
	@echo "  make down-all            - Stop all services"
	@echo "  make logs-all            - View logs from all services"
	@echo "  make reset-all           - Stop, remove volumes, and restart all services"
	@echo "  make health-check        - Check health of all backend services"
	@echo "  make verify              - Full verification: build, start, test, and verify everything"
	@echo "  make verify-logs         - Dump logs for failing services"
	@echo "  make dev-up              - Start all services WITHOUT search (OpenSearch optional)"
	@echo "  make dev-up-safe         - Start all services with automatic Docker recovery (macOS)"
	@echo "  make dev-up-search       - Start all services WITH search profile (OpenSearch enabled)"
	@echo "  make dev-down            - Stop all services"
	@echo "  make dev-down-safe       - Stop all services with automatic Docker recovery (macOS)"
	@echo "  make docker-recover      - Restart Docker Desktop on macOS when daemon is dead"
	@echo "  make migrate-all         - Run all database migrations"
	@echo "  make seed-all            - Seed all databases with test data"
	@echo "  make test                - Run all unit tests"
	@echo "  make test-integration    - Run integration tests (requires services running)"
	@echo "  make run-<service>       - Run a specific service locally"
	@echo "  make clean               - Clean up generated files"

up-all:
	@echo "Checking Docker readiness..."
	@bash scripts/check-docker.sh > /dev/null 2>&1 || (echo "❌ Docker not ready. Run: bash scripts/check-docker.sh" && exit 1)
	@echo "✅ Docker is ready"
	@echo ""
	@echo "Starting ALL services (infrastructure + backend + frontend) WITHOUT search..."
	@docker compose -f docker-compose.all.yml up -d --build
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Services started (without OpenSearch). Run 'make migrate-all' to set up databases, then 'make seed-all' for test data."
	@echo "Frontend: http://localhost:3002"
	@echo "Note: Search features will show 'temporarily unavailable'. Use 'make dev-up-search' to enable search."
	@echo "Check status: docker compose -f docker-compose.all.yml ps"

down-all:
	@echo "Stopping all services..."
	@docker compose -f docker-compose.all.yml down

logs-all:
	@docker compose -f docker-compose.all.yml logs -f

reset-all:
	@echo "Checking Docker readiness..."
	@bash scripts/check-docker.sh > /dev/null 2>&1 || (echo "❌ Docker not ready. Run: bash scripts/check-docker.sh" && exit 1)
	@echo "✅ Docker is ready"
	@echo ""
	@echo "Stopping all services and removing volumes..."
	@docker compose -f docker-compose.all.yml down -v
	@echo "Starting fresh..."
	@docker compose -f docker-compose.all.yml up -d --build
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Services restarted. Run 'make migrate-all' to set up databases."

smoke:
	@echo "Running smoke tests..."
	@bash scripts/smoke.sh

health-check:
	@echo "Checking health of all backend services..."
	@echo ""
	@curl -f -s -X GET http://localhost:8001/health > /dev/null 2>&1 && echo "✅ identity-service (port 8001): OK" || echo "❌ identity-service (port 8001): FAIL"
	@curl -f -s -X GET http://localhost:8002/health > /dev/null 2>&1 && echo "✅ company-service (port 8002): OK" || echo "❌ company-service (port 8002): FAIL"
	@curl -f -s -X GET http://localhost:8003/health > /dev/null 2>&1 && echo "✅ catalog-service (port 8003): OK" || echo "❌ catalog-service (port 8003): FAIL"
	@curl -f -s -X GET http://localhost:8004/health > /dev/null 2>&1 && echo "✅ equipment-service (port 8004): OK" || echo "❌ equipment-service (port 8004): FAIL"
	@curl -f -s -X GET http://localhost:8005/health > /dev/null 2>&1 && echo "✅ marketplace-service (port 8005): OK" || echo "❌ marketplace-service (port 8005): FAIL"
	@curl -f -s -X GET http://localhost:8006/health > /dev/null 2>&1 && echo "✅ procurement-service (port 8006): OK" || echo "❌ procurement-service (port 8006): FAIL"
	@curl -f -s -X GET http://localhost:8007/health > /dev/null 2>&1 && echo "✅ logistics-service (port 8007): OK" || echo "❌ logistics-service (port 8007): FAIL"
	@curl -f -s -X GET http://localhost:8008/health > /dev/null 2>&1 && echo "✅ collaboration-service (port 8008): OK" || echo "❌ collaboration-service (port 8008): FAIL"
	@curl -f -s -X GET http://localhost:8009/health > /dev/null 2>&1 && echo "✅ notification-service (port 8009): OK" || echo "❌ notification-service (port 8009): FAIL"
	@curl -f -s -X GET http://localhost:8010/health > /dev/null 2>&1 && echo "✅ billing-service (port 8010): OK" || echo "❌ billing-service (port 8010): FAIL"
	@curl -f -s -X GET http://localhost:8011/health > /dev/null 2>&1 && echo "✅ virtual-warehouse-service (port 8011): OK" || echo "❌ virtual-warehouse-service (port 8011): FAIL"
	@curl -f -s -X GET http://localhost:8012/health > /dev/null 2>&1 && echo "✅ search-indexer-service (port 8012): OK" || echo "❌ search-indexer-service (port 8012): FAIL"
	@curl -f -s -X GET http://localhost:8013/health > /dev/null 2>&1 && echo "✅ diagnostics-service (port 8013): OK" || echo "❌ diagnostics-service (port 8013): FAIL"
	@echo ""
	@echo "Frontend check:"
	@curl -f -s http://localhost:3002 > /dev/null 2>&1 && echo "✅ frontend (port 3002): OK" || echo "❌ frontend (port 3002): FAIL"

verify:
	@echo "=========================================="
	@echo "VERIFICATION GATE - Definition of Done"
	@echo "=========================================="
	@echo ""
	@mkdir -p reports/logs
	@echo "Step 0/9: Checking Docker readiness..."
	@bash scripts/check-docker.sh || (echo "❌ FAIL: Docker not ready" && exit 1)
	@echo "✅ Docker is ready"
	@echo ""
	@echo "Step 1/9: Starting all services..."
	@docker compose -f docker-compose.all.yml up -d --build || (echo "❌ FAIL: Service startup failed" && exit 1)
	@echo "✅ Services started"
	@echo ""
	@echo "Step 2/9: Waiting for services to be ready..."
	@python3 scripts/wait_for_ready.py || (echo "❌ FAIL: Services not ready within timeout" && exit 1)
	@echo "✅ All services ready"
	@echo ""
	@echo "Step 3/9: Running migrations..."
	@make migrate-all || (echo "❌ FAIL: Migrations failed" && exit 1)
	@echo "✅ Migrations completed"
	@echo ""
	@echo "Step 4/9: Seeding databases..."
	@make seed-all || (echo "❌ FAIL: Seeding failed" && exit 1)
	@echo "✅ Seeding completed"
	@echo ""
	@echo "Step 5/9: Running unit tests..."
	@make test || (echo "❌ FAIL: Unit tests failed" && exit 1)
	@echo "✅ Unit tests passed"
	@echo ""
	@echo "Step 6/9: Running integration tests..."
	@make test-integration || (echo "❌ FAIL: Integration tests failed" && exit 1)
	@echo "✅ Integration tests passed"
	@echo ""
	@echo "Step 7/9: Checking frontend build..."
	@docker compose -f docker-compose.all.yml exec -T frontend npm run lint || (echo "❌ FAIL: Frontend lint failed" && exit 1)
	@docker compose -f docker-compose.all.yml exec -T frontend npm run build || (echo "❌ FAIL: Frontend build failed" && exit 1)
	@echo "✅ Frontend build passed"
	@echo ""
	@echo "Step 8/9: Running smoke tests..."
	@cd frontend && npm run test:e2e || (echo "❌ FAIL: Smoke tests failed" && exit 1)
	@echo "✅ Smoke tests passed"
	@echo ""
	@echo "Step 9/9: Generating verification report..."
	@python3 scripts/generate_verify_report.py
	@echo ""
	@echo "=========================================="
	@echo "✅ VERIFICATION PASSED - All checks OK"
	@echo "=========================================="
	@echo "Report: reports/verify_report.md"

verify-logs:
	@echo "Collecting logs for failing services..."
	@mkdir -p reports/logs
	@python3 scripts/collect_failure_logs.py

dev-up:
	@echo "Checking Docker readiness..."
	@bash scripts/check-docker.sh > /dev/null 2>&1 || (echo "❌ Docker not ready. Run: bash scripts/check-docker.sh" && exit 1)
	@echo "✅ Docker is ready"
	@echo ""
	@echo "Checking service versions..."
	@bash scripts/check-versions.sh > /tmp/version-check.log 2>&1; \
	VERSION_CHECK_EXIT=$$?; \
	cat /tmp/version-check.log; \
	echo ""; \
	if [ $$VERSION_CHECK_EXIT -eq 0 ]; then \
		echo "Services need rebuild. Getting list..."; \
		SERVICES_TO_BUILD=$$(bash scripts/check-versions.sh --rebuild-list 2>/dev/null); \
		if [ -n "$$SERVICES_TO_BUILD" ]; then \
			echo "Building changed services:"; \
			echo "$$SERVICES_TO_BUILD" | while read service; do \
				if [ -n "$$service" ] && [ "$$service" != "" ]; then \
					echo "  - $$service"; \
				fi; \
			done; \
			echo ""; \
			echo "Building services (this may take a while and use significant memory)..."; \
			echo "If you get 'cannot allocate memory' errors, increase Docker Desktop memory to at least 8GB"; \
			echo "$$SERVICES_TO_BUILD" | while read service; do \
				if [ -n "$$service" ] && [ "$$service" != "" ]; then \
					echo "Building $$service..."; \
					VERSION=$$(bash scripts/get-service-version.sh $$service 2>/dev/null || echo "unknown"); \
					docker compose -f docker-compose.all.yml build --build-arg SERVICE_VERSION=$$VERSION $$service || true; \
				fi; \
			done; \
			SERVICES_LIST=$$(echo "$$SERVICES_TO_BUILD" | tr '\n' ' '); \
			docker compose -f docker-compose.all.yml up -d $$SERVICES_LIST; \
			bash scripts/update-container-versions.sh; \
			echo "Waiting for services to be ready..."; \
			sleep 10; \
			echo "Services started (without OpenSearch). Run 'make migrate-all' to set up databases, then 'make seed-all' for test data."; \
			echo "Frontend: http://localhost:3002"; \
			echo "Note: Search features will show 'temporarily unavailable'"; \
			echo "Check status: docker compose -f docker-compose.all.yml ps"; \
		fi; \
	else \
		running=$$(docker compose -f docker-compose.all.yml ps --format json 2>/dev/null | grep -c '"State":"running"' || echo "0"); \
		total=$$(docker compose -f docker-compose.all.yml config --services 2>/dev/null | wc -l | tr -d ' ' || echo "0"); \
		if [ "$$running" -gt 0 ] && [ "$$running" -ge $$((total / 2)) ]; then \
			echo "✅ Containers are already running ($$running services)"; \
			echo "Skipping build/start. Using existing containers."; \
			echo "To rebuild, run: make dev-down && make dev-up"; \
		else \
			echo "Starting all services WITHOUT search (OpenSearch disabled)..."; \
			echo "Building services (this may take a while and use significant memory)..."; \
			echo "If you get 'cannot allocate memory' errors, increase Docker Desktop memory to at least 8GB"; \
			docker compose -f docker-compose.all.yml build; \
			docker compose -f docker-compose.all.yml up -d; \
			bash scripts/update-container-versions.sh; \
			echo "Waiting for services to be ready..."; \
			sleep 10; \
			echo "Services started (without OpenSearch). Run 'make migrate-all' to set up databases, then 'make seed-all' for test data."; \
			echo "Frontend: http://localhost:3002"; \
			echo "Note: Search features will show 'temporarily unavailable'"; \
			echo "Check status: docker compose -f docker-compose.all.yml ps"; \
		fi; \
	fi

dev-up-search:
	@echo "Checking Docker readiness..."
	@bash scripts/check-docker.sh > /dev/null 2>&1 || (echo "❌ Docker not ready. Run: bash scripts/check-docker.sh" && exit 1)
	@echo "✅ Docker is ready"
	@echo ""
	@echo "Starting all services WITH search profile (OpenSearch enabled)..."
	@docker compose -f docker-compose.all.yml --profile search up -d --build
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Services started (with OpenSearch). Run 'make migrate-all' to set up databases, then 'make seed-all' for test data."
	@echo "Frontend: http://localhost:3002"
	@echo "Check status: docker compose -f docker-compose.all.yml ps"

dev-down:
	@echo "Stopping all services..."
	@docker compose -f docker-compose.all.yml down

dev-up-safe:
	@echo "Starting services with automatic Docker recovery (macOS)..."
	@bash scripts/dev-up-safe.sh

dev-down-safe:
	@echo "Stopping services with automatic Docker recovery (macOS)..."
	@bash scripts/dev-down-safe.sh

migrate-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: SERVICE variable is required. Usage: make migrate-service SERVICE=<service-name>"; \
		exit 1; \
	fi
	@echo "Running migrations for service: $(SERVICE)..."
	@cd tools/migrate && go run main.go --service $(SERVICE) || exit 1
	@echo "✅ Migrations completed for $(SERVICE)"

migrate-all:
	@echo "Running migrations for all services..."
	@for service in identity company catalog equipment marketplace procurement logistics collaboration notification billing virtual-warehouse diagnostics; do \
		echo ""; \
		echo "=== Migrating $$service-service ==="; \
		$(MAKE) migrate-service SERVICE=$$service-service || exit 1; \
	done
	@echo ""
	@echo "✅ All migrations completed successfully!"

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
	@cd services/search-indexer-service && go run cmd/seed/main.go
	@cd services/diagnostics-service && go run cmd/seed/main.go
	@echo "All seeds completed."

test:
	@echo "Running unit tests for all services..."
	@for service in identity company catalog equipment marketplace procurement logistics collaboration notification billing virtual-warehouse search-indexer diagnostics; do \
		echo "Testing $$service-service..."; \
		cd services/$$service-service && go test ./... -v || exit 1; \
	done
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

clean:
	@echo "Cleaning up..."
	@find . -type f -name "*.test" -delete
	@find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	@echo "Cleanup completed."
