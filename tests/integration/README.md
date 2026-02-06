# Integration Tests

This directory contains integration tests for the B2B Procurement Marketplace Platform.

## Overview

Integration tests verify the system end-to-end against the running local stack:
- PostgreSQL (database)
- Redis (cache and events)
- OpenSearch (search)
- MinIO (file storage)

## Test Files

### `integration_test.go`
Main test file that ensures services are ready before running tests.

### `helpers.go`
Utility functions for:
- HTTP client with authentication
- Login functionality
- Service health checks
- Response parsing

### `tenant_isolation_test.go`
Tests tenant isolation:
- Verifies tenants cannot access each other's data
- Tests company and PR isolation across tenants
- Ensures multi-tenancy is properly enforced

### `rbac_test.go`
Tests role-based access control:
- Buyer can create PRs but cannot approve
- Approver can approve PRs
- Catalog admin can manage catalog
- Supplier cannot create PRs
- Admin has full access

### `e2e_flow_test.go`
Tests complete end-to-end procurement flow:
1. Create Purchase Request (PR) as buyer
2. Approve PR as admin/approver
3. Create RFQ
4. Submit Quote as supplier
5. Create Purchase Order (PO)
6. Create Shipment
7. Trigger late shipment alert
8. Verify notification creation

### `search_indexing_test.go`
Tests search indexing:
- Catalog parts indexed after approval
- Companies indexed after approval
- Verifies OpenSearch integration

## Running Tests

### Prerequisites

1. Start all services:
   ```bash
   make dev-up
   ```

2. Run migrations:
   ```bash
   make migrate-all
   ```

3. Seed test data:
   ```bash
   make seed-all
   ```

4. Ensure all services are running and healthy

### Run Tests

```bash
make test-integration
```

Or directly:
```bash
cd tests/integration
go test ./... -v
```

### Run Specific Test

```bash
cd tests/integration
go test -v -run TestTenantIsolation
go test -v -run TestRBACEnforcement
go test -v -run TestEndToEndFlow
go test -v -run TestSearchIndexing
```

## Test Data

Tests use seeded demo accounts:
- **Admin**: `admin@demo.com` / `demo123456`
- **Buyer**: `buyer@demo.com` / `demo123456`
- **Supplier**: `supplier@demo.com` / `demo123456`
- **Procurement Manager**: `procurement@demo.com` / `demo123456`

All accounts use tenant ID: `00000000-0000-0000-0000-000000000001`

## Repeatability

Tests are designed to be repeatable:
- Use unique IDs for all created resources
- Tests can be run multiple times without conflicts
- No manual cleanup required (unique identifiers prevent collisions)

## Service URLs

Tests connect to services running on:
- Identity Service: `http://localhost:8001`
- Company Service: `http://localhost:8002`
- Catalog Service: `http://localhost:8003`
- Procurement Service: `http://localhost:8006`
- Logistics Service: `http://localhost:8007`
- Notification Service: `http://localhost:8009`
- OpenSearch: `http://localhost:9200`

## Troubleshooting

### Tests fail with connection errors
- Ensure all services are running: `docker-compose ps`
- Check service health: `curl http://localhost:8001/health`
- Verify services are listening on correct ports

### Tests fail with authentication errors
- Ensure seed data was created: `make seed-all`
- Verify demo accounts exist in database
- Check JWT secret matches across services

### Tests fail with database errors
- Ensure migrations ran: `make migrate-all`
- Check database connection settings
- Verify PostgreSQL is running: `docker-compose ps postgres`

### OpenSearch indexing tests fail
- Wait longer for indexing (tests wait 5 seconds)
- Check OpenSearch is running: `curl http://localhost:9200/_cluster/health`
- Verify search-indexer-service is running and consuming events

## Adding New Tests

1. Create a new test file: `new_feature_test.go`
2. Import required packages:
   ```go
   package integration
   
   import (
       "testing"
       "github.com/stretchr/testify/require"
       "github.com/stretchr/testify/assert"
   )
   ```
3. Use helper functions from `helpers.go`
4. Use unique IDs for test data
5. Add test to appropriate test file or create new one

## Best Practices

- Always use unique IDs for test data
- Wait for services to be ready before running tests
- Clean up test data or use unique identifiers
- Verify both success and failure cases
- Test error handling and edge cases
- Use descriptive test names
- Add logging for debugging
