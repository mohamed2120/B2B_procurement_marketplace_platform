# B2B Procurement Marketplace Platform

A comprehensive enterprise-grade B2B procurement marketplace platform built with microservices architecture.

## Architecture

- **Backend**: Golang microservices
- **Frontend**: React (Next.js)
- **Database**: PostgreSQL 15 (one database, multiple schemas)
- **Cache**: Redis
- **Search**: OpenSearch
- **File Storage**: MinIO (local), S3 (production)
- **Auth**: JWT-based with RBAC
- **Events**: Event-driven architecture with Redis pub/sub

## Services

1. **identity-service** (Port 8001) - User management, roles, permissions, JWT
2. **company-service** (Port 8002) - Company profiles, verification, subdomain management
3. **catalog-service** (Port 8003) - Global catalog, manufacturers, parts library
4. **equipment-service** (Port 8004) - Equipment, BOM, compatibility
5. **marketplace-service** (Port 8005) - Stores, listings, stock, pricing
6. **procurement-service** (Port 8006) - PR, RFQ, quotes, PO, orders
7. **logistics-service** (Port 8007) - Shipments, ETA, tracking, alerts
8. **collaboration-service** (Port 8008) - Chat, disputes, ratings
9. **notification-service** (Port 8009) - Templates, preferences, notifications
10. **billing-service** (Port 8010) - Plans, subscriptions, entitlements
11. **virtual-warehouse-service** (Port 8011) - Shared inventory, transfers
12. **search-indexer-service** (Port 8012) - Event consumer, OpenSearch updates

## RUN EVERYTHING LOCALLY

### Prerequisites

- **Docker Desktop** (must be running)
- Make (optional, but recommended)
- Go 1.23+ (for running migrations/seeds locally)

### Quick Start Scripts

We provide convenient start scripts for different platforms:

**For Unix/Mac/Linux:**
```bash
./start.sh
```

**For Windows (PowerShell):**
```powershell
.\start.ps1
```

These scripts will:
- Start all infrastructure, backend services, and frontend
- Run database migrations
- Seed the database with demo data
- Check service health
- Display access information

### Manual Start (Alternative)

If you prefer to start manually or don't have the scripts:

Start all infrastructure, backend services, and frontend with a single command:

```bash
make up-all
```

This will:
- Start PostgreSQL, Redis, OpenSearch, and MinIO
- Build and start all 12 backend microservices
- Build and start the Next.js frontend
- Expose all services on their respective ports

Then run migrations and seeds:
```bash
make migrate-all
make seed-all
```

### Service URLs

| Service | Local URL | Port |
|---------|-----------|------|
| Frontend | http://localhost:3000 | 3000 |
| identity-service | http://localhost:8001 | 8001 |
| company-service | http://localhost:8002 | 8002 |
| catalog-service | http://localhost:8003 | 8003 |
| equipment-service | http://localhost:8004 | 8004 |
| marketplace-service | http://localhost:8005 | 8005 |
| procurement-service | http://localhost:8006 | 8006 |
| logistics-service | http://localhost:8007 | 8007 |
| collaboration-service | http://localhost:8008 | 8008 |
| notification-service | http://localhost:8009 | 8009 |
| billing-service | http://localhost:8010 | 8010 |
| virtual-warehouse-service | http://localhost:8011 | 8011 |
| search-indexer-service | http://localhost:8012 | 8012 |
| PostgreSQL | localhost:5432 | 5432 |
| Redis | localhost:6379 | 6379 |
| OpenSearch | http://localhost:9200 | 9200 |
| MinIO Console | http://localhost:9001 | 9001 |

### Check Container Status

```bash
docker compose -f docker-compose.all.yml ps
```

### Health Check

Verify all services are running and healthy:

```bash
make health-check
```

This will check all `/health` endpoints and report OK/FAIL for each service.

### Stop Everything

```bash
make down-all
```

### View Logs

```bash
make logs-all
```

Or view logs for a specific service:

```bash
docker compose -f docker-compose.all.yml logs -f identity-service
```

### Reset Everything

Stop, remove volumes, and restart fresh:

```bash
make reset-all
```

**Note:** After starting services, you still need to run migrations and seed data:

```bash
make migrate-all
make seed-all
```

---

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+
- Node.js 18+ (for frontend)
- Make

### Local Development Setup

1. **Clone and setup environment**:
   ```bash
   cp env.example .env
   # Edit .env if needed
   ```

2. **Start infrastructure services**:
   ```bash
   make dev-up
   ```

   This starts:
   - PostgreSQL on port 5432
   - Redis on port 6379
   - OpenSearch on port 9200
   - MinIO on ports 9000 (API) and 9001 (Console)

3. **Run migrations**:
   ```bash
   make migrate-all
   ```

4. **Seed data**:
   ```bash
   make seed-all
   ```

5. **Start services** (in separate terminals or use docker-compose):
   ```bash
   make run-identity
   make run-company
   make run-procurement
   # ... etc
   ```

   Or use docker-compose:
   ```bash
   docker-compose -f deployments/docker-compose.yml up
   ```

### Frontend

The frontend is built with Next.js (React) and provides a modern web interface.

**Setup:**
```bash
cd frontend
npm install
cp .env.local.example .env.local
```

**Development:**
```bash
# From project root
make run-frontend

# Or from frontend directory
cd frontend && npm run dev
```

The frontend will be available at [http://localhost:3000](http://localhost:3000)

**Features:**
- Tenant resolution via subdomain
- JWT-based authentication
- Role-based routing and guards
- Customer screens: PR list/create, RFQ list, Quote compare/award, Orders/Shipments
- Supplier screens: RFQ inbox, Quote submit, Listings CRUD
- Admin screens: Company verification, Catalog approvals
- Shared: Notifications center, Chat panel

**E2E Tests:**
```bash
cd frontend
npm run test:e2e
```

See `frontend/README.md` for more details.

### Testing

#### Unit Tests

Run all unit tests:
```bash
make test
```

This runs unit tests for all services.

#### Integration Tests

Integration tests verify the system end-to-end against the running local stack (PostgreSQL, Redis, OpenSearch, MinIO).

**Prerequisites:**
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

4. Ensure all services are running (either via docker-compose or locally)

**Run Integration Tests:**
```bash
make test-integration
```

**What Integration Tests Cover:**

1. **Tenant Isolation Test** (`tenant_isolation_test.go`)
   - Verifies that tenants cannot access each other's data
   - Tests company and PR isolation across tenants

2. **RBAC Enforcement Test** (`rbac_test.go`)
   - Verifies role-based access control
   - Tests that buyers can create PRs but not approve them
   - Tests that approvers can approve PRs
   - Tests that catalog admins can manage catalog
   - Tests that suppliers cannot create PRs

3. **End-to-End Flow Test** (`e2e_flow_test.go`)
   - Tests complete procurement flow:
     - Create Purchase Request (PR)
     - Approve PR
     - Create RFQ
     - Submit Quote
     - Create Purchase Order (PO)
     - Create Shipment
     - Trigger late shipment alert
     - Verify notification creation

4. **Search Indexing Test** (`search_indexing_test.go`)
   - Verifies that catalog parts are indexed in OpenSearch after approval
   - Verifies that companies are indexed after approval
   - Tests search functionality

**Test Structure:**
- Tests are located in `/tests/integration/`
- Tests use unique IDs to ensure repeatability
- Tests clean up after themselves or use unique identifiers
- Tests wait for services to be ready before running

**Troubleshooting:**
- If tests fail, ensure all services are running and healthy
- Check service health endpoints: `curl http://localhost:8001/health`
- Verify database migrations completed successfully
- Ensure seed data was created (check for demo users)

### Demo Accounts

After seeding, you can use these accounts (all use password: `demo123456`):

- **Platform Admin**: `admin@demo.com` - Full platform access
- **Requester (Buyer Company)**: `buyer.requester@demo.com` - Can create PRs and RFQs
- **Procurement Manager (Buyer Company)**: `buyer.procurement@demo.com` - Can approve PRs/POs and award quotes
- **Supplier**: `supplier@demo.com` - Can manage listings and respond to RFQs

All accounts use tenant ID: `00000000-0000-0000-0000-000000000001`

## API Endpoints

### Identity Service

- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/register` - Register
- `GET /api/v1/auth/validate` - Validate token (protected)

### Company Service

- `GET /api/v1/companies` - List companies
- `GET /api/v1/companies/:id` - Get company
- `POST /api/v1/companies` - Create company
- `POST /api/v1/companies/:id/approve` - Approve company
- `POST /api/v1/companies/:id/subdomain-request` - Request subdomain

### Procurement Service

- `GET /api/v1/purchase-requests` - List PRs
- `POST /api/v1/purchase-requests` - Create PR
- `POST /api/v1/purchase-requests/:id/approve` - Approve PR
- `POST /api/v1/rfqs` - Create RFQ
- `POST /api/v1/quotes` - Submit quote
- `POST /api/v1/purchase-orders` - Create PO

## AWS Deployment

### Terraform Infrastructure

The project includes Terraform configurations for deploying to AWS staging environment.

**Location**: `terraform/`

**Infrastructure Components**:
- VPC with 2 AZs (public/private subnets)
- ECS Fargate cluster for microservices
- Application Load Balancer with path-based routing
- RDS PostgreSQL (private subnets)
- ElastiCache Redis
- OpenSearch domain
- S3 buckets (docs-private, media)
- Cognito user pool
- EventBridge bus + SQS queues
- IAM roles and security groups

**Quick Start**:
```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
terraform init
terraform plan  # Review changes
# terraform apply  # Uncomment to deploy (will create AWS resources)
```

**⚠️ Important**: 
- Do not run `terraform apply` without reviewing the plan
- Configure S3 backend for state management
- Use AWS Secrets Manager for production credentials
- Review cost implications before deploying

See `terraform/README.md` for detailed deployment instructions.

## End-to-End Flow Testing

### 1. Create Purchase Request

```bash
curl -X POST http://localhost:8006/api/v1/purchase-requests \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Office Supplies",
    "description": "Need office supplies",
    "items": [
      {
        "description": "Printer Paper",
        "quantity": 100,
        "unit": "ream"
      }
    ]
  }'
```

### 2. Approve PR

```bash
curl -X POST http://localhost:8006/api/v1/purchase-requests/<pr_id>/approve \
  -H "Authorization: Bearer <token>"
```

### 3. Create RFQ

```bash
curl -X POST http://localhost:8006/api/v1/rfqs \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pr_id": "<pr_id>",
    "title": "RFQ for Office Supplies",
    "due_date": "2024-12-31T00:00:00Z"
  }'
```

### 4. Submit Quote (as supplier)

```bash
curl -X POST http://localhost:8006/api/v1/quotes \
  -H "Authorization: Bearer <supplier_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "rfq_id": "<rfq_id>",
    "supplier_id": "<supplier_id>",
    "items": [
      {
        "pr_item_id": "<pr_item_id>",
        "description": "Printer Paper",
        "quantity": 100,
        "unit_price": 25.00
      }
    ],
    "total_amount": 2500.00
  }'
```

### 5. Create Purchase Order

```bash
curl -X POST http://localhost:8006/api/v1/purchase-orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pr_id": "<pr_id>",
    "rfq_id": "<rfq_id>",
    "quote_id": "<quote_id>",
    "supplier_id": "<supplier_id>"
  }'
```

## Database Schemas

Each service has its own PostgreSQL schema:

- `identity` - Users, roles, permissions
- `company` - Companies, documents, subdomain requests
- `catalog` - Manufacturers, categories, parts
- `equipment` - Equipment, BOM, compatibility
- `marketplace` - Stores, listings, stock
- `procurement` - PR, RFQ, quotes, PO
- `logistics` - Shipments, tracking
- `collaboration` - Chat, disputes, ratings
- `notification` - Templates, preferences
- `billing` - Plans, subscriptions
- `virtual_warehouse` - Shared inventory
- `audit` - Audit logs

## Event System

Events are published via Redis pub/sub. Key events:

- `core.company.approved.v1`
- `catalog.lib_part.approved.v1`
- `procurement.pr.approved.v1`
- `procurement.rfq.created.v1`
- `procurement.quote.submitted.v1`
- `procurement.order.placed.v1`
- `logistics.shipment.late.v1`
- `collab.chat.message_sent.v1`
- `billing.subscription.started.v1`

## Development

### Adding a New Service

1. Create service directory: `services/<service-name>/`
2. Add `go.mod`, models, migrations, handlers, services
3. Add to `docker-compose.yml`
4. Add migration and seed commands to Makefile
5. Update this README

### Running Individual Services

```bash
cd services/<service-name>
go run cmd/server/main.go
```

### Database Migrations

Migrations are SQL files in `services/<service>/migrations/`. Run with:

```bash
cd services/<service>
go run cmd/migrate/main.go
```

## Production Deployment

### AWS Deployment

Terraform configurations are in `infrastructure/terraform/` (skeleton provided).

Key components:
- ECS/Fargate for services
- RDS PostgreSQL
- ElastiCache Redis
- OpenSearch Service
- S3 for file storage
- EventBridge + SQS for events
- Cognito for authentication

### Environment Variables

See `env.example` for all required environment variables.

## Testing

### Unit Tests

Each service includes unit tests:

```bash
cd services/<service>
go test ./...
```

### Integration Tests

Integration tests require a running database:

```bash
make test
```

### End-to-End Tests

See "End-to-End Flow Testing" section above.

## Troubleshooting

### Services won't start

1. Check database is running: `docker ps`
2. Check environment variables: `cat .env`
3. Check logs: `docker-compose logs <service-name>`

### Database connection errors

1. Verify PostgreSQL is running: `docker ps | grep postgres`
2. Check connection string in `.env`
3. Verify migrations ran: `make migrate-all`

### Redis connection errors

1. Verify Redis is running: `docker ps | grep redis`
2. Check Redis connection in `.env`

## License

[Your License Here]

## Contributing

[Contributing Guidelines]
