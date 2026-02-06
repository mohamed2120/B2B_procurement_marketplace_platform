-- Create equipment schema
CREATE SCHEMA IF NOT EXISTS equipment;

-- Equipment table
CREATE TABLE IF NOT EXISTS equipment.equipment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    equipment_number VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100),
    manufacturer VARCHAR(255),
    model VARCHAR(255),
    serial_number VARCHAR(255),
    year INTEGER,
    status VARCHAR(50) DEFAULT 'active',
    location VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT idx_tenant_eq_num UNIQUE (tenant_id, equipment_number)
);

CREATE INDEX idx_equipment_tenant_id ON equipment.equipment(tenant_id);
CREATE INDEX idx_equipment_status ON equipment.equipment(status);
CREATE INDEX idx_equipment_type ON equipment.equipment(type);
CREATE INDEX idx_equipment_deleted_at ON equipment.equipment(deleted_at);

-- BOM nodes table
CREATE TABLE IF NOT EXISTS equipment.bom_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    equipment_id UUID NOT NULL REFERENCES equipment.equipment(id) ON DELETE CASCADE,
    part_id UUID,
    part_number VARCHAR(255),
    part_name VARCHAR(255) NOT NULL,
    description TEXT,
    quantity DECIMAL(15,2) NOT NULL DEFAULT 1,
    unit VARCHAR(50),
    position VARCHAR(100),
    parent_node_id UUID,
    level INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_node_id) REFERENCES equipment.bom_nodes(id) ON DELETE CASCADE
);

CREATE INDEX idx_bom_nodes_tenant_id ON equipment.bom_nodes(tenant_id);
CREATE INDEX idx_bom_nodes_equipment_id ON equipment.bom_nodes(equipment_id);
CREATE INDEX idx_bom_nodes_part_id ON equipment.bom_nodes(part_id);
CREATE INDEX idx_bom_nodes_parent_node_id ON equipment.bom_nodes(parent_node_id);

-- Compatibility mappings table
CREATE TABLE IF NOT EXISTS equipment.compatibility_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    equipment_id UUID NOT NULL REFERENCES equipment.equipment(id) ON DELETE CASCADE,
    part_id UUID NOT NULL,
    is_compatible BOOLEAN DEFAULT true,
    notes TEXT,
    verified_by UUID,
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_equipment_part UNIQUE (equipment_id, part_id)
);

CREATE INDEX idx_compatibility_tenant_id ON equipment.compatibility_mappings(tenant_id);
CREATE INDEX idx_compatibility_equipment_id ON equipment.compatibility_mappings(equipment_id);
CREATE INDEX idx_compatibility_part_id ON equipment.compatibility_mappings(part_id);
