# Payments/Escrow MVP Verification Checklist

## A1) Database Migrations ✅/❌

- [ ] procurement.purchase_orders has payment_mode (DIRECT|ESCROW)
- [ ] procurement.purchase_orders has payment_status (pending|processing|succeeded|failed)
- [ ] procurement.purchase_orders has payment_id (nullable UUID, optional)
- [ ] billing.payments table exists
- [ ] billing.escrow_holds table exists
- [ ] billing.settlements table exists
- [ ] billing.refunds table exists
- [ ] billing.payout_accounts table exists
- [ ] All migrations apply successfully with `make migrate-all`

## A2) Billing Service Endpoints ✅/❌

- [ ] POST /api/billing/v1/payments/intent exists and works
- [ ] POST /api/billing/v1/payments/webhook exists (mock provider supported)
- [ ] POST /api/billing/v1/escrow/release exists and works
- [ ] POST /api/billing/v1/refunds exists and works
- [ ] GET /api/billing/v1/payout-accounts exists (list)
- [ ] POST /api/billing/v1/payout-accounts exists (create)
- [ ] GET /api/billing/v1/payout-accounts/:id exists (get)
- [ ] PUT /api/billing/v1/payout-accounts/:id exists (update)
- [ ] DELETE /api/billing/v1/payout-accounts/:id exists (delete)
- [ ] GET /health returns 200

## A3) Procurement Service Integration ✅/❌

- [ ] Order creation accepts payment_mode field
- [ ] Order creation accepts payment_status field
- [ ] If payment_mode=ESCROW, procurement calls billing-service /payments/intent
- [ ] order.payment_status updated to "processing" after intent creation
- [ ] order.payment_status updated to "failed" if intent creation fails
- [ ] BILLING_SERVICE_URL env var used for billing service base URL
- [ ] PUT /api/v1/purchase-orders/:id/payment-status exists for billing callbacks

## A4) Escrow Rules Enforcement ✅/❌

- [ ] Payment success webhook creates escrow_hold with status="held"
- [ ] Escrow release requires delivery confirmation OR auto-release days passed
- [ ] Escrow release blocked if dispute.status="open"
- [ ] Auto-release function exists (ProcessAutoRelease)
- [ ] Dispute check function exists (CheckDisputeStatus)

## A5) Events Implementation ✅/❌

- [ ] billing.payment.succeeded.v1 published on payment success
- [ ] billing.payment.failed.v1 published on payment failure
- [ ] billing.escrow.held.v1 published when escrow hold created
- [ ] billing.settlement.released.v1 published when escrow released
- [ ] billing.refund.issued.v1 published when refund created
- [ ] docs/events.md includes all 5 billing events
- [ ] notification-service consumes billing events
- [ ] notification-service creates in-app notifications for billing events

## A6) Seed Data ✅/❌

- [ ] Demo customer company + user seeded
- [ ] Demo supplier company + user seeded
- [ ] Supplier payout account seeded
- [ ] One DIRECT order seeded with payment
- [ ] One ESCROW order seeded with payment and escrow_hold
- [ ] One shipment linked to escrow order with ETA
- [ ] Shipment has delivered flag or proof_of_delivery field

## A7) Integration Tests ✅/❌

- [ ] Test 1: Escrow happy path (create order → intent → webhook → hold → delivery → release → settlement)
- [ ] Test 2: Dispute blocks release (open dispute → attempt release → verify blocked)
- [ ] Test 3: Auto-release (set auto-release days → verify release without manual action)
- [ ] All tests use unique IDs and are repeatable
- [ ] Tests run with `make test-integration`

## Verification Commands ✅/❌

- [ ] `make verify-escrow` exists and runs all checks
- [ ] `make test-integration` exists and runs integration tests
- [ ] All health endpoints return 200 after `make dev-up`
