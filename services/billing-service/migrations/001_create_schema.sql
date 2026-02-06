-- Create billing schema
CREATE SCHEMA IF NOT EXISTS billing;

-- Plans table
CREATE TABLE IF NOT EXISTS billing.plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    code VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    price DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    billing_cycle VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_plans_code ON billing.plans(code);
CREATE INDEX IF NOT EXISTS idx_plans_is_active ON billing.plans(is_active);

-- Entitlements table
CREATE TABLE IF NOT EXISTS billing.entitlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES billing.plans(id) ON DELETE CASCADE,
    feature VARCHAR(100) NOT NULL,
    limit INTEGER,
    unit VARCHAR(50),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_entitlements_plan_id ON billing.entitlements(plan_id);
CREATE INDEX IF NOT EXISTS idx_entitlements_feature ON billing.entitlements(feature);

-- Subscriptions table
CREATE TABLE IF NOT EXISTS billing.subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    plan_id UUID NOT NULL REFERENCES billing.plans(id) ON DELETE RESTRICT,
    status VARCHAR(50) DEFAULT 'active',
    started_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    auto_renew BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT unique_tenant_active_subscription UNIQUE (tenant_id, status) WHERE status = 'active'
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_tenant_id ON billing.subscriptions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_id ON billing.subscriptions(plan_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON billing.subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_deleted_at ON billing.subscriptions(deleted_at);
