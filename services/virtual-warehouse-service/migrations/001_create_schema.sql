-- Create virtual_warehouse schema
CREATE SCHEMA IF NOT EXISTS virtual_warehouse;

-- Shared inventory table
CREATE TABLE IF NOT EXISTS virtual_warehouse.shared_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    part_id UUID NOT NULL,
    equipment_id UUID,
    quantity DECIMAL(15,2) NOT NULL,
    location VARCHAR(255),
    is_available BOOLEAN DEFAULT true,
    reserved_qty DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_shared_inventory_tenant_id ON virtual_warehouse.shared_inventory(tenant_id);
CREATE INDEX idx_shared_inventory_part_id ON virtual_warehouse.shared_inventory(part_id);
CREATE INDEX idx_shared_inventory_equipment_id ON virtual_warehouse.shared_inventory(equipment_id);
CREATE INDEX idx_shared_inventory_is_available ON virtual_warehouse.shared_inventory(is_available);
CREATE INDEX idx_shared_inventory_deleted_at ON virtual_warehouse.shared_inventory(deleted_at);

-- Equipment groups table
CREATE TABLE IF NOT EXISTS virtual_warehouse.equipment_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_equipment_groups_tenant_id ON virtual_warehouse.equipment_groups(tenant_id);
CREATE INDEX idx_equipment_groups_deleted_at ON virtual_warehouse.equipment_groups(deleted_at);

-- Equipment group members table
CREATE TABLE IF NOT EXISTS virtual_warehouse.equipment_group_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES virtual_warehouse.equipment_groups(id) ON DELETE CASCADE,
    equipment_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_group_equipment UNIQUE (group_id, equipment_id)
);

CREATE INDEX idx_equipment_group_members_group_id ON virtual_warehouse.equipment_group_members(group_id);
CREATE INDEX idx_equipment_group_members_equipment_id ON virtual_warehouse.equipment_group_members(equipment_id);

-- Inter-company transfers table
CREATE TABLE IF NOT EXISTS virtual_warehouse.inter_company_transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_tenant_id UUID NOT NULL,
    to_tenant_id UUID NOT NULL,
    part_id UUID NOT NULL,
    quantity DECIMAL(15,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    requested_by UUID NOT NULL,
    approved_by UUID,
    rejection_reason TEXT,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_transfers_from_tenant_id ON virtual_warehouse.inter_company_transfers(from_tenant_id);
CREATE INDEX idx_transfers_to_tenant_id ON virtual_warehouse.inter_company_transfers(to_tenant_id);
CREATE INDEX idx_transfers_part_id ON virtual_warehouse.inter_company_transfers(part_id);
CREATE INDEX idx_transfers_status ON virtual_warehouse.inter_company_transfers(status);
CREATE INDEX idx_transfers_deleted_at ON virtual_warehouse.inter_company_transfers(deleted_at);

-- Emergency sourcing table
CREATE TABLE IF NOT EXISTS virtual_warehouse.emergency_sourcing (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    part_id UUID NOT NULL,
    equipment_id UUID,
    quantity DECIMAL(15,2) NOT NULL,
    priority VARCHAR(50) DEFAULT 'high',
    status VARCHAR(50) DEFAULT 'open',
    requested_by UUID NOT NULL,
    fulfilled_by UUID,
    fulfilled_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_emergency_sourcing_tenant_id ON virtual_warehouse.emergency_sourcing(tenant_id);
CREATE INDEX idx_emergency_sourcing_part_id ON virtual_warehouse.emergency_sourcing(part_id);
CREATE INDEX idx_emergency_sourcing_status ON virtual_warehouse.emergency_sourcing(status);
CREATE INDEX idx_emergency_sourcing_priority ON virtual_warehouse.emergency_sourcing(priority);
CREATE INDEX idx_emergency_sourcing_deleted_at ON virtual_warehouse.emergency_sourcing(deleted_at);
