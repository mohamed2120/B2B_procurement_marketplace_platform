-- Create catalog schema
CREATE SCHEMA IF NOT EXISTS catalog;

-- Manufacturers table
CREATE TABLE IF NOT EXISTS catalog.manufacturers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    code VARCHAR(100) UNIQUE,
    website VARCHAR(255),
    country VARCHAR(100),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_manufacturers_name ON catalog.manufacturers(name);
CREATE INDEX idx_manufacturers_code ON catalog.manufacturers(code);

-- Categories table
CREATE TABLE IF NOT EXISTS catalog.categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    code VARCHAR(100) UNIQUE,
    description TEXT,
    parent_id UUID,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES catalog.categories(id) ON DELETE SET NULL
);

CREATE INDEX idx_categories_name ON catalog.categories(name);
CREATE INDEX idx_categories_code ON catalog.categories(code);
CREATE INDEX idx_categories_parent_id ON catalog.categories(parent_id);

-- Attributes table
CREATE TABLE IF NOT EXISTS catalog.attributes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    code VARCHAR(100) UNIQUE,
    data_type VARCHAR(50) NOT NULL,
    unit VARCHAR(50),
    is_required BOOLEAN DEFAULT false,
    is_searchable BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attributes_name ON catalog.attributes(name);
CREATE INDEX idx_attributes_code ON catalog.attributes(code);

-- Library parts table
CREATE TABLE IF NOT EXISTS catalog.library_parts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    part_number VARCHAR(255) NOT NULL,
    manufacturer_id UUID NOT NULL REFERENCES catalog.manufacturers(id) ON DELETE RESTRICT,
    category_id UUID NOT NULL REFERENCES catalog.categories(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    approved_at TIMESTAMP,
    approved_by UUID,
    rejected_reason TEXT,
    is_duplicate BOOLEAN DEFAULT false,
    duplicate_of UUID,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT idx_part_number_mfr UNIQUE (part_number, manufacturer_id)
);

CREATE INDEX idx_library_parts_part_number ON catalog.library_parts(part_number);
CREATE INDEX idx_library_parts_manufacturer_id ON catalog.library_parts(manufacturer_id);
CREATE INDEX idx_library_parts_category_id ON catalog.library_parts(category_id);
CREATE INDEX idx_library_parts_status ON catalog.library_parts(status);
CREATE INDEX idx_library_parts_is_duplicate ON catalog.library_parts(is_duplicate);
CREATE INDEX idx_library_parts_duplicate_of ON catalog.library_parts(duplicate_of);
CREATE INDEX idx_library_parts_deleted_at ON catalog.library_parts(deleted_at);

-- Part attributes table
CREATE TABLE IF NOT EXISTS catalog.part_attributes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    part_id UUID NOT NULL REFERENCES catalog.library_parts(id) ON DELETE CASCADE,
    attribute_id UUID NOT NULL REFERENCES catalog.attributes(id) ON DELETE CASCADE,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_part_attribute UNIQUE (part_id, attribute_id)
);

CREATE INDEX idx_part_attributes_part_id ON catalog.part_attributes(part_id);
CREATE INDEX idx_part_attributes_attribute_id ON catalog.part_attributes(attribute_id);
