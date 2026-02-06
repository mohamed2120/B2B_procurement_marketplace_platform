# Event System Documentation

## Overview

The platform uses an event-driven architecture with Redis pub/sub for local development. Events follow a consistent naming pattern: `<domain>.<entity>.<action>.v1`

## Event Definitions

### 1. `core.company.approved.v1`
- **Source**: company-service
- **Published**: ✅ Yes
- **Consumed By**: 
  - search-indexer-service ✅
  - notification-service ✅ (notifies company admins)
- **Payload**:
  - `company_id` (string)
  - `name` (string)
  - `subdomain` (string)
- **When**: Company approval workflow completes

### 2. `catalog.lib_part.approved.v1`
- **Source**: catalog-service
- **Published**: ✅ Yes
- **Consumed By**:
  - search-indexer-service ✅
  - notification-service ✅ (notifies part creator)
- **Payload**:
  - `part_id` (string)
  - `part_number` (string)
  - `name` (string)
  - `manufacturer_id` (string)
- **When**: Catalog part is approved by admin

### 3. `procurement.pr.approved.v1`
- **Source**: procurement-service
- **Published**: ✅ Yes
- **Consumed By**:
  - notification-service ✅ (notifies requester and approvers)
- **Payload**:
  - `pr_id` (string)
  - `pr_number` (string)
- **When**: Purchase Request is approved

### 4. `procurement.rfq.created.v1`
- **Source**: procurement-service
- **Published**: ✅ Yes
- **Consumed By**:
  - search-indexer-service ✅ (indexes RFQ)
  - notification-service ✅ (notifies suppliers)
- **Payload**:
  - `rfq_id` (string)
  - `rfq_number` (string)
  - `pr_id` (string)
- **When**: RFQ is created and sent to suppliers

### 5. `procurement.quote.submitted.v1`
- **Source**: procurement-service
- **Published**: ✅ Yes
- **Consumed By**:
  - search-indexer-service ✅ (indexes quote)
  - notification-service ✅ (notifies buyers)
- **Payload**:
  - `quote_id` (string)
  - `quote_number` (string)
  - `rfq_id` (string)
  - `supplier_id` (string)
- **When**: Supplier submits a quote

### 6. `procurement.order.placed.v1`
- **Source**: procurement-service
- **Published**: ✅ Yes
- **Consumed By**:
  - search-indexer-service ✅
  - notification-service ✅ (notifies buyer and supplier)
- **Payload**:
  - `po_id` (string)
  - `po_number` (string)
  - `pr_id` (string)
  - `quote_id` (string)
- **When**: Purchase Order is created

### 7. `logistics.shipment.late.v1`
- **Source**: logistics-service
- **Published**: ✅ Yes
- **Consumed By**:
  - search-indexer-service ✅ (indexes shipment)
  - notification-service ✅ (notifies buyer and logistics team)
- **Payload**:
  - `shipment_id` (string)
  - `tracking_number` (string)
  - `eta` (timestamp)
- **When**: Shipment ETA is exceeded

### 8. `collab.chat.message_sent.v1`
- **Source**: collaboration-service
- **Published**: ✅ Yes
- **Consumed By**:
  - notification-service ✅ (notifies thread participants)
- **Payload**:
  - `message_id` (string)
  - `thread_id` (string)
  - `sender_id` (string)
- **When**: Chat message is sent

### 9. `billing.subscription.started.v1`
- **Source**: billing-service
- **Published**: ✅ Yes
- **Consumed By**:
  - notification-service ✅ (notifies company admins)
- **Payload**:
  - `subscription_id` (string)
  - `tenant_id` (string)
  - `plan_id` (string)
- **When**: Subscription is activated

### 10. `billing.payment.succeeded.v1`
- **Source**: billing-service
- **Published**: ✅ Yes
- **Consumed By**: None (future: notification-service)
- **Payload**:
  - `payment_id` (string)
  - `payment_intent_id` (string)
  - `order_id` (string)
  - `amount` (float64)
  - `payment_mode` (string) - DIRECT or ESCROW
- **When**: Payment is successfully processed

### 11. `billing.payment.failed.v1`
- **Source**: billing-service
- **Published**: ✅ Yes
- **Consumed By**: None (future: notification-service)
- **Payload**:
  - `payment_id` (string)
  - `payment_intent_id` (string)
  - `order_id` (string)
  - `failed_reason` (string)
- **When**: Payment processing fails

### 12. `billing.escrow.held.v1`
- **Source**: billing-service
- **Published**: ✅ Yes
- **Consumed By**: None (future: notification-service)
- **Payload**:
  - `escrow_hold_id` (string)
  - `payment_id` (string)
  - `order_id` (string)
  - `amount` (float64)
- **When**: Funds are held in escrow after payment succeeds

### 13. `billing.settlement.released.v1`
- **Source**: billing-service
- **Published**: ✅ Yes
- **Consumed By**: None (future: notification-service)
- **Payload**:
  - `settlement_id` (string)
  - `escrow_hold_id` (string)
  - `order_id` (string)
  - `supplier_id` (string)
  - `amount` (float64)
- **When**: Escrow funds are released to supplier

### 14. `billing.refund.issued.v1`
- **Source**: billing-service
- **Published**: ✅ Yes
- **Consumed By**: None (future: notification-service)
- **Payload**:
  - `refund_id` (string)
  - `refund_number` (string)
  - `payment_id` (string)
  - `order_id` (string)
  - `amount` (float64)
- **When**: Refund is issued for a payment

## Event Status Summary

| Event | Published | Search Indexer | Notification Service |
|-------|-----------|----------------|---------------------|
| core.company.approved.v1 | ✅ | ✅ | ✅ |
| catalog.lib_part.approved.v1 | ✅ | ✅ | ✅ |
| procurement.pr.approved.v1 | ✅ | ❌ | ✅ |
| procurement.rfq.created.v1 | ✅ | ✅ | ✅ |
| procurement.quote.submitted.v1 | ✅ | ✅ | ✅ |
| procurement.order.placed.v1 | ✅ | ✅ | ✅ |
| logistics.shipment.late.v1 | ✅ | ✅ | ✅ |
| collab.chat.message_sent.v1 | ✅ | ❌ | ✅ |
| billing.subscription.started.v1 | ✅ | ❌ | ✅ |
| billing.payment.succeeded.v1 | ✅ | ❌ | ❌ |
| billing.payment.failed.v1 | ✅ | ❌ | ❌ |
| billing.escrow.held.v1 | ✅ | ❌ | ❌ |
| billing.settlement.released.v1 | ✅ | ❌ | ❌ |
| billing.refund.issued.v1 | ✅ | ❌ | ❌ |

## Implementation Notes

- All events use Redis pub/sub for local development
- Events include tenant_id for multi-tenancy
- Events are serialized as JSON
- Event bus interface allows swapping implementations (local Redis → AWS EventBridge)
- **notification-service** consumes all events and creates notifications based on templates and user preferences
- **search-indexer-service** consumes events for indexing: parts, companies, orders, RFQs, quotes, and shipments

## Event Consumers

### notification-service
- Consumes all 9 events
- Creates notifications based on templates
- Respects user preferences (channel, event type)
- Supports in-app and email notifications (email mocked locally)

### search-indexer-service
- Consumes events for search indexing:
  - `catalog.lib_part.approved.v1` → indexes parts
  - `core.company.approved.v1` → indexes companies
  - `procurement.order.placed.v1` → indexes orders
  - `procurement.rfq.created.v1` → indexes RFQs
  - `procurement.quote.submitted.v1` → indexes quotes
  - `logistics.shipment.late.v1` → indexes shipments
