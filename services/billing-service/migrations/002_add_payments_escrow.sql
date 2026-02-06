-- Add payment_mode to purchase_orders (in procurement schema)
ALTER TABLE procurement.purchase_orders 
ADD COLUMN IF NOT EXISTS payment_mode VARCHAR(50) DEFAULT 'DIRECT';

CREATE INDEX IF NOT EXISTS idx_po_payment_mode ON procurement.purchase_orders(payment_mode);

-- Payments table
CREATE TABLE IF NOT EXISTS billing.payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    order_id UUID NOT NULL,
    payment_intent_id VARCHAR(255) NOT NULL UNIQUE,
    provider VARCHAR(50) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    status VARCHAR(50) DEFAULT 'pending',
    payment_mode VARCHAR(50) NOT NULL,
    metadata JSONB,
    failed_reason TEXT,
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payments_tenant_id ON billing.payments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON billing.payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON billing.payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_payment_mode ON billing.payments(payment_mode);
CREATE INDEX IF NOT EXISTS idx_payments_payment_intent_id ON billing.payments(payment_intent_id);

-- Escrow holds table
CREATE TABLE IF NOT EXISTS billing.escrow_holds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    payment_id UUID NOT NULL UNIQUE REFERENCES billing.payments(id) ON DELETE RESTRICT,
    order_id UUID NOT NULL,
    supplier_id UUID NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    status VARCHAR(50) DEFAULT 'held',
    auto_release_days INTEGER DEFAULT 30,
    auto_release_date TIMESTAMP,
    released_at TIMESTAMP,
    released_by UUID,
    release_reason TEXT,
    blocked_by_dispute BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_escrow_holds_tenant_id ON billing.escrow_holds(tenant_id);
CREATE INDEX IF NOT EXISTS idx_escrow_holds_payment_id ON billing.escrow_holds(payment_id);
CREATE INDEX IF NOT EXISTS idx_escrow_holds_order_id ON billing.escrow_holds(order_id);
CREATE INDEX IF NOT EXISTS idx_escrow_holds_supplier_id ON billing.escrow_holds(supplier_id);
CREATE INDEX IF NOT EXISTS idx_escrow_holds_status ON billing.escrow_holds(status);
CREATE INDEX IF NOT EXISTS idx_escrow_holds_blocked_by_dispute ON billing.escrow_holds(blocked_by_dispute);
CREATE INDEX IF NOT EXISTS idx_escrow_holds_auto_release_date ON billing.escrow_holds(auto_release_date);

-- Settlements table
CREATE TABLE IF NOT EXISTS billing.settlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    escrow_hold_id UUID NOT NULL REFERENCES billing.escrow_holds(id) ON DELETE RESTRICT,
    supplier_id UUID NOT NULL,
    payout_account_id UUID,
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    status VARCHAR(50) DEFAULT 'pending',
    provider_payout_id VARCHAR(255),
    failed_reason TEXT,
    settled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_settlements_tenant_id ON billing.settlements(tenant_id);
CREATE INDEX IF NOT EXISTS idx_settlements_escrow_hold_id ON billing.settlements(escrow_hold_id);
CREATE INDEX IF NOT EXISTS idx_settlements_supplier_id ON billing.settlements(supplier_id);
CREATE INDEX IF NOT EXISTS idx_settlements_payout_account_id ON billing.settlements(payout_account_id);
CREATE INDEX IF NOT EXISTS idx_settlements_status ON billing.settlements(status);

-- Refunds table
CREATE TABLE IF NOT EXISTS billing.refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    payment_id UUID NOT NULL REFERENCES billing.payments(id) ON DELETE RESTRICT,
    order_id UUID NOT NULL,
    refund_number VARCHAR(50) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    reason TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    provider_refund_id VARCHAR(255),
    failed_reason TEXT,
    refunded_at TIMESTAMP,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_tenant_refund_number UNIQUE (tenant_id, refund_number)
);

CREATE INDEX IF NOT EXISTS idx_refunds_tenant_id ON billing.refunds(tenant_id);
CREATE INDEX IF NOT EXISTS idx_refunds_payment_id ON billing.refunds(payment_id);
CREATE INDEX IF NOT EXISTS idx_refunds_order_id ON billing.refunds(order_id);
CREATE INDEX IF NOT EXISTS idx_refunds_status ON billing.refunds(status);
CREATE INDEX IF NOT EXISTS idx_refunds_refund_number ON billing.refunds(refund_number);

-- Payout accounts table
CREATE TABLE IF NOT EXISTS billing.payout_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    supplier_id UUID NOT NULL,
    account_type VARCHAR(50) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    account_number VARCHAR(255),
    routing_number VARCHAR(50),
    account_holder_name VARCHAR(255),
    bank_name VARCHAR(255),
    provider_account_id VARCHAR(255),
    is_default BOOLEAN DEFAULT false,
    is_verified BOOLEAN DEFAULT false,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payout_accounts_tenant_id ON billing.payout_accounts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payout_accounts_supplier_id ON billing.payout_accounts(supplier_id);
CREATE INDEX IF NOT EXISTS idx_payout_accounts_provider_account_id ON billing.payout_accounts(provider_account_id);
CREATE INDEX IF NOT EXISTS idx_payout_accounts_is_default ON billing.payout_accounts(is_default);
