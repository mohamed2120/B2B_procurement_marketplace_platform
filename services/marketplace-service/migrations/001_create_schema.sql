-- Create marketplace schema
CREATE SCHEMA IF NOT EXISTS marketplace;

-- Stores table
CREATE TABLE IF NOT EXISTS marketplace.stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    company_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active',
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stores_tenant_id ON marketplace.stores(tenant_id);
CREATE INDEX IF NOT EXISTS idx_stores_company_id ON marketplace.stores(company_id);
CREATE INDEX IF NOT EXISTS idx_stores_status ON marketplace.stores(status);
CREATE INDEX IF NOT EXISTS idx_stores_deleted_at ON marketplace.stores(deleted_at);

-- Listings table
CREATE TABLE IF NOT EXISTS marketplace.listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    store_id UUID NOT NULL REFERENCES marketplace.stores(id) ON DELETE CASCADE,
    listing_type VARCHAR(50) NOT NULL,
    part_id UUID,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    sku VARCHAR(100),
    status VARCHAR(50) DEFAULT 'draft',
    price DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    stock_quantity DECIMAL(15,2) DEFAULT 0,
    min_order_quantity DECIMAL(15,2) DEFAULT 1,
    lead_time_days INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_listings_tenant_id ON marketplace.listings(tenant_id);
CREATE INDEX IF NOT EXISTS idx_listings_store_id ON marketplace.listings(store_id);
CREATE INDEX IF NOT EXISTS idx_listings_listing_type ON marketplace.listings(listing_type);
CREATE INDEX IF NOT EXISTS idx_listings_part_id ON marketplace.listings(part_id);
CREATE INDEX IF NOT EXISTS idx_listings_status ON marketplace.listings(status);
CREATE INDEX IF NOT EXISTS idx_listings_is_active ON marketplace.listings(is_active);
CREATE INDEX IF NOT EXISTS idx_listings_deleted_at ON marketplace.listings(deleted_at);

-- Listing media table
CREATE TABLE IF NOT EXISTS marketplace.listing_media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id UUID NOT NULL REFERENCES marketplace.listings(id) ON DELETE CASCADE,
    media_type VARCHAR(50) NOT NULL,
    url TEXT NOT NULL,
    thumbnail_url TEXT,
    file_name VARCHAR(255),
    file_size BIGINT,
    mime_type VARCHAR(100),
    is_primary BOOLEAN DEFAULT false,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_listing_media_listing_id ON marketplace.listing_media(listing_id);
CREATE INDEX IF NOT EXISTS idx_listing_media_is_primary ON marketplace.listing_media(is_primary);
