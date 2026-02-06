-- Create company schema
CREATE SCHEMA IF NOT EXISTS company;

-- Companies table
CREATE TABLE IF NOT EXISTS company.companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255),
    tax_id VARCHAR(100),
    subdomain VARCHAR(100) UNIQUE,
    status VARCHAR(50) DEFAULT 'pending',
    verification_status VARCHAR(50) DEFAULT 'pending',
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    phone VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),
    industry VARCHAR(100),
    company_type VARCHAR(50),
    approved_at TIMESTAMP,
    approved_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_companies_status ON company.companies(status);
CREATE INDEX IF NOT EXISTS idx_companies_subdomain ON company.companies(subdomain);
CREATE INDEX IF NOT EXISTS idx_companies_deleted_at ON company.companies(deleted_at);

-- Company documents table
CREATE TABLE IF NOT EXISTS company.documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES company.companies(id) ON DELETE CASCADE,
    document_type VARCHAR(100) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_url TEXT NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(100),
    uploaded_by UUID NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    verified_at TIMESTAMP,
    verified_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_documents_company_id ON company.documents(company_id);
CREATE INDEX IF NOT EXISTS idx_documents_status ON company.documents(status);

-- Subdomain requests table
CREATE TABLE IF NOT EXISTS company.subdomain_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES company.companies(id) ON DELETE CASCADE,
    subdomain VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(50) DEFAULT 'pending',
    requested_by UUID NOT NULL,
    reviewed_by UUID,
    reviewed_at TIMESTAMP,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_subdomain_requests_company_id ON company.subdomain_requests(company_id);
CREATE INDEX IF NOT EXISTS idx_subdomain_requests_status ON company.subdomain_requests(status);
CREATE INDEX IF NOT EXISTS idx_subdomain_requests_subdomain ON company.subdomain_requests(subdomain);
