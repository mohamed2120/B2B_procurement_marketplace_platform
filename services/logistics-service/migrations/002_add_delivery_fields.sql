-- Add delivery confirmation fields to shipments
ALTER TABLE logistics.shipments 
ADD COLUMN IF NOT EXISTS delivered_at TIMESTAMP;

ALTER TABLE logistics.shipments 
ADD COLUMN IF NOT EXISTS proof_of_delivery TEXT;

ALTER TABLE logistics.shipments 
ADD COLUMN IF NOT EXISTS delivery_status VARCHAR(50) DEFAULT 'in_transit';

CREATE INDEX IF NOT EXISTS idx_shipments_delivery_status ON logistics.shipments(delivery_status);
CREATE INDEX IF NOT EXISTS idx_shipments_delivered_at ON logistics.shipments(delivered_at);
