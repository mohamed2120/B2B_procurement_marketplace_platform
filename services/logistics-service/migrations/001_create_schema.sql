-- Create logistics schema
CREATE SCHEMA IF NOT EXISTS logistics;

-- Shipments table
CREATE TABLE IF NOT EXISTS logistics.shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    po_id UUID NOT NULL,
    tracking_number VARCHAR(100) UNIQUE,
    status VARCHAR(50) DEFAULT 'pending',
    carrier VARCHAR(100),
    eta TIMESTAMP NOT NULL,
    actual_delivery_date TIMESTAMP,
    origin VARCHAR(255),
    destination VARCHAR(255),
    is_late BOOLEAN DEFAULT false,
    late_alert_sent BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_shipments_tenant_id ON logistics.shipments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_shipments_po_id ON logistics.shipments(po_id);
CREATE INDEX IF NOT EXISTS idx_shipments_status ON logistics.shipments(status);
CREATE INDEX IF NOT EXISTS idx_shipments_eta ON logistics.shipments(eta);
CREATE INDEX IF NOT EXISTS idx_shipments_is_late ON logistics.shipments(is_late);

-- Tracking events table
CREATE TABLE IF NOT EXISTS logistics.tracking_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_id UUID NOT NULL REFERENCES logistics.shipments(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    location VARCHAR(255),
    description TEXT,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tracking_events_shipment_id ON logistics.tracking_events(shipment_id);

-- Proof of delivery table
CREATE TABLE IF NOT EXISTS logistics.proof_of_delivery (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_id UUID NOT NULL UNIQUE REFERENCES logistics.shipments(id) ON DELETE CASCADE,
    signed_by VARCHAR(255),
    signature_url TEXT,
    delivered_at TIMESTAMP NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pod_shipment_id ON logistics.proof_of_delivery(shipment_id);
