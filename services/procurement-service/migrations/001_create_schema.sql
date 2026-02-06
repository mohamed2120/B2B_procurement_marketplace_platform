-- Create procurement schema
CREATE SCHEMA IF NOT EXISTS procurement;

-- Purchase Requests table
CREATE TABLE IF NOT EXISTS procurement.purchase_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    pr_number VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'draft',
    priority VARCHAR(50) DEFAULT 'normal',
    requested_by UUID NOT NULL,
    department VARCHAR(100),
    budget DECIMAL(15,2),
    currency VARCHAR(10) DEFAULT 'USD',
    required_date TIMESTAMP,
    approved_at TIMESTAMP,
    approved_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT idx_tenant_pr UNIQUE (tenant_id, pr_number)
);

CREATE INDEX IF NOT EXISTS idx_pr_tenant_id ON procurement.purchase_requests(tenant_id);
CREATE INDEX IF NOT EXISTS idx_pr_status ON procurement.purchase_requests(status);
CREATE INDEX IF NOT EXISTS idx_pr_deleted_at ON procurement.purchase_requests(deleted_at);

-- PR Items table
CREATE TABLE IF NOT EXISTS procurement.pr_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pr_id UUID NOT NULL REFERENCES procurement.purchase_requests(id) ON DELETE CASCADE,
    part_id UUID,
    description TEXT NOT NULL,
    quantity DECIMAL(15,2) NOT NULL,
    unit VARCHAR(50),
    unit_price DECIMAL(15,2),
    total_price DECIMAL(15,2),
    specifications TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pr_items_pr_id ON procurement.pr_items(pr_id);
CREATE INDEX IF NOT EXISTS idx_pr_items_part_id ON procurement.pr_items(part_id);

-- PR Approvals table
CREATE TABLE IF NOT EXISTS procurement.pr_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pr_id UUID NOT NULL REFERENCES procurement.purchase_requests(id) ON DELETE CASCADE,
    approver_id UUID NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    comments TEXT,
    approved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pr_approvals_pr_id ON procurement.pr_approvals(pr_id);
CREATE INDEX IF NOT EXISTS idx_pr_approvals_status ON procurement.pr_approvals(status);

-- RFQs table
CREATE TABLE IF NOT EXISTS procurement.rfqs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    pr_id UUID NOT NULL REFERENCES procurement.purchase_requests(id) ON DELETE CASCADE,
    rfq_number VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'draft',
    due_date TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT idx_tenant_rfq UNIQUE (tenant_id, rfq_number)
);

CREATE INDEX IF NOT EXISTS idx_rfqs_tenant_id ON procurement.rfqs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_rfqs_pr_id ON procurement.rfqs(pr_id);
CREATE INDEX IF NOT EXISTS idx_rfqs_status ON procurement.rfqs(status);

-- Quotes table
CREATE TABLE IF NOT EXISTS procurement.quotes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    rfq_id UUID NOT NULL REFERENCES procurement.rfqs(id) ON DELETE CASCADE,
    supplier_id UUID NOT NULL,
    quote_number VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'submitted',
    total_amount DECIMAL(15,2),
    currency VARCHAR(10) DEFAULT 'USD',
    valid_until TIMESTAMP,
    notes TEXT,
    submitted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT idx_tenant_quote UNIQUE (tenant_id, quote_number)
);

CREATE INDEX IF NOT EXISTS idx_quotes_tenant_id ON procurement.quotes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_quotes_rfq_id ON procurement.quotes(rfq_id);
CREATE INDEX IF NOT EXISTS idx_quotes_supplier_id ON procurement.quotes(supplier_id);
CREATE INDEX IF NOT EXISTS idx_quotes_status ON procurement.quotes(status);

-- Quote Items table
CREATE TABLE IF NOT EXISTS procurement.quote_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quote_id UUID NOT NULL REFERENCES procurement.quotes(id) ON DELETE CASCADE,
    pr_item_id UUID NOT NULL,
    description TEXT NOT NULL,
    quantity DECIMAL(15,2) NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2),
    lead_time INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_quote_items_quote_id ON procurement.quote_items(quote_id);
CREATE INDEX IF NOT EXISTS idx_quote_items_pr_item_id ON procurement.quote_items(pr_item_id);

-- Purchase Orders table
CREATE TABLE IF NOT EXISTS procurement.purchase_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    pr_id UUID NOT NULL REFERENCES procurement.purchase_requests(id) ON DELETE CASCADE,
    rfq_id UUID REFERENCES procurement.rfqs(id),
    quote_id UUID NOT NULL REFERENCES procurement.quotes(id),
    po_number VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    total_amount DECIMAL(15,2),
    currency VARCHAR(10) DEFAULT 'USD',
    supplier_id UUID NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT idx_tenant_po UNIQUE (tenant_id, po_number)
);

CREATE INDEX IF NOT EXISTS idx_pos_tenant_id ON procurement.purchase_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_pos_pr_id ON procurement.purchase_orders(pr_id);
CREATE INDEX IF NOT EXISTS idx_pos_quote_id ON procurement.purchase_orders(quote_id);
CREATE INDEX IF NOT EXISTS idx_pos_status ON procurement.purchase_orders(status);

-- PO Items table
CREATE TABLE IF NOT EXISTS procurement.po_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    po_id UUID NOT NULL REFERENCES procurement.purchase_orders(id) ON DELETE CASCADE,
    pr_item_id UUID NOT NULL,
    description TEXT NOT NULL,
    quantity DECIMAL(15,2) NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_po_items_po_id ON procurement.po_items(po_id);
CREATE INDEX IF NOT EXISTS idx_po_items_pr_item_id ON procurement.po_items(pr_item_id);
