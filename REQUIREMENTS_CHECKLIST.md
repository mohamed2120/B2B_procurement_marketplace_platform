# Requirements Checklist

## Infrastructure & Setup
- [x] docker-compose.yml with postgres, redis, opensearch, minio
- [x] .env.local.example
- [x] Makefile with dev-up, dev-down, migrate-all, seed-all, test, run-<service>
- [x] README with setup, testing, verification, deployment instructions

## Shared Libraries
- [x] Common auth middleware (JWT, tenant resolution)
- [x] RBAC middleware
- [x] Event system (local pub/sub)
- [x] Database utilities
- [x] Common models/types

## Microservices (Each must have: main.go, router, auth/RBAC middleware, DB layer, migrations, seeds, tests, /health, OpenAPI, Dockerfile)

### 1. identity-service
- [x] Users, roles, permissions tables
- [x] JWT generation/validation
- [x] Tenant resolution
- [x] User invitations
- [x] Migrations
- [x] Seeds
- [x] Tests

### 2. company-service
- [x] Company profiles
- [x] Verification workflow
- [x] Documents
- [x] Subdomain request & approval
- [x] Migrations
- [x] Seeds
- [x] Tests

### 3. catalog-service (GLOBAL)
- [x] Manufacturers
- [x] Categories
- [x] Attributes
- [x] Spare parts library (anti-duplicate)
- [x] Admin approval
- [x] Events on approval
- [x] Migrations
- [x] Seeds
- [x] Tests

### 4. equipment-service
- [x] Equipment
- [x] BOM
- [x] Compatibility mapping
- [x] Compatibility verification
- [x] Migrations
- [x] Seeds
- [x] Tests

### 5. marketplace-service
- [x] Stores
- [x] Listings (product/service/surplus)
- [x] Services
- [x] Stock & pricing
- [x] Lead time
- [x] Listing media
- [x] Migrations
- [x] Seeds
- [x] Tests

### 6. procurement-service
- [x] PR (Purchase Requests)
- [x] Approvals
- [x] RFQ
- [x] Quotes
- [x] Award
- [x] PO
- [x] Order
- [x] Payment mode support (DIRECT | ESCROW)
- [x] Payment status tracking
- [x] Integration with billing-service for ESCROW payments
- [x] Migrations
- [x] Seeds
- [x] Integration tests (escrow flow)
- [x] Unit tests

### 7. logistics-service
- [x] Shipments
- [x] ETA
- [x] Tracking
- [x] Late alerts
- [x] Proof of delivery
- [x] Migrations
- [x] Seeds
- [x] Tests

### 8. collaboration-service
- [x] Chat (threads/messages/files)
- [x] Disputes
- [x] Ratings
- [x] Moderation hooks
- [x] Migrations
- [x] Seeds
- [x] Tests

### 9. notification-service
- [x] Templates
- [x] Preferences
- [x] In-app notifications
- [x] Email (mock locally)
- [x] Migrations
- [x] Seeds
- [x] Tests

### 10. billing-service
- [x] Plans, entitlements, subscriptions
- [x] Payments (DIRECT and ESCROW modes)
- [x] Escrow holds and releases
- [x] Settlements
- [x] Refunds
- [x] Payout accounts (CRUD)
- [x] Mock payment provider for local dev
- [x] Webhook handler
- [x] Escrow policy (hold on payment, release on delivery/auto-release, block on dispute)
- [x] Migrations
- [x] Seeds (escrow order, direct order, payout account)
- [x] Tests
- [x] Plans
- [x] Entitlements
- [x] Subscriptions
- [x] Events on subscription start
- [x] Migrations
- [x] Seeds
- [x] Tests

### 11. virtual-warehouse-service
- [x] Shared inventory
- [x] Equipment groups
- [x] Inter-company transfers
- [x] Emergency sourcing
- [x] Migrations
- [x] Seeds
- [x] Tests

### 12. search-indexer-service
- [x] Event consumer
- [x] OpenSearch updates
- [x] Document indexing
- [x] Tests

## Database
- [x] All schemas created (identity, company, catalog, equipment, marketplace, procurement, logistics, collaboration, notification, billing, virtual_warehouse, audit)
- [x] All migrations for all services (12/12)
- [x] UUID primary keys
- [x] tenant_id on all tenant-owned tables
- [x] Foreign keys enforced
- [x] Indexes on tenant_id, status, dates

## Event System
- [x] Common event envelope
- [x] core.company.approved.v1 (published ✅, consumed by search-indexer ✅, notification ✅)
- [x] catalog.lib_part.approved.v1 (published ✅, consumed by search-indexer ✅, notification ✅)
- [x] procurement.pr.approved.v1 (published ✅, consumed by notification ✅)
- [x] procurement.rfq.created.v1 (published ✅, consumed by search-indexer ✅, notification ✅)
- [x] procurement.quote.submitted.v1 (published ✅, consumed by search-indexer ✅, notification ✅)
- [x] procurement.order.placed.v1 (published ✅, consumed by search-indexer ✅, notification ✅)
- [x] logistics.shipment.late.v1 (published ✅, consumed by search-indexer ✅, notification ✅)
- [x] collab.chat.message_sent.v1 (published ✅, consumed by notification ✅)
- [x] billing.subscription.started.v1 (published ✅, consumed by notification ✅)
- [x] billing.payment.succeeded.v1 (published ✅)
- [x] billing.payment.failed.v1 (published ✅)
- [x] billing.escrow.held.v1 (published ✅)
- [x] billing.settlement.released.v1 (published ✅)
- [x] billing.refund.issued.v1 (published ✅)
- [x] Local pub/sub implementation
- [x] notification-service event consumer (consumes all 9 events)
- [x] search-indexer-service event consumer (consumes 6 events)
- [x] Event tests (publishing and consumption)

## Frontend (Next.js)
- [ ] Role-based routes (/customer, /supplier, /admin)
- [ ] Tenant resolution from subdomain
- [ ] JWT attachment to API calls
- [ ] Route permission enforcement
- [ ] Seeded demo accounts
- [ ] Visual validation of all flows

## Testing
- [x] Unit tests per service
- [x] Integration tests with Postgres
- [x] End-to-end flow tests:
  - [x] Create PR
  - [x] Approve PR
  - [x] Create RFQ
  - [x] Supplier submits quote
  - [x] Award quote
  - [x] Create order
  - [x] Create shipment
  - [x] Trigger late ETA alert
  - [x] Send notification
  - [x] Chat messages exchanged
- [x] make test command
- [ ] Test coverage summary

## AWS Deployment
- [ ] Terraform skeleton
- [ ] No AWS calls in local runtime
