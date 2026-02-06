-- Add payment_mode and payment_status to purchase_orders
ALTER TABLE procurement.purchase_orders 
ADD COLUMN IF NOT EXISTS payment_mode VARCHAR(50) DEFAULT 'DIRECT';

ALTER TABLE procurement.purchase_orders 
ADD COLUMN IF NOT EXISTS payment_status VARCHAR(50) DEFAULT 'pending';

ALTER TABLE procurement.purchase_orders 
ADD COLUMN IF NOT EXISTS payment_id UUID;

CREATE INDEX IF NOT EXISTS idx_po_payment_mode ON procurement.purchase_orders(payment_mode);
CREATE INDEX IF NOT EXISTS idx_po_payment_status ON procurement.purchase_orders(payment_status);
CREATE INDEX IF NOT EXISTS idx_po_payment_id ON procurement.purchase_orders(payment_id);