-- Create identity schema
CREATE SCHEMA IF NOT EXISTS identity;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE IF NOT EXISTS identity.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT idx_tenant_email UNIQUE (tenant_id, email)
);

CREATE INDEX idx_users_tenant_id ON identity.users(tenant_id);
CREATE INDEX idx_users_email ON identity.users(email);
CREATE INDEX idx_users_deleted_at ON identity.users(deleted_at);

-- Roles table
CREATE TABLE IF NOT EXISTS identity.roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_system BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_roles_name ON identity.roles(name);

-- Permissions table
CREATE TABLE IF NOT EXISTS identity.permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_permissions_resource ON identity.permissions(resource);
CREATE INDEX idx_permissions_resource_action ON identity.permissions(resource, action);

-- Role permissions table
CREATE TABLE IF NOT EXISTS identity.role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES identity.roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES identity.permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_role_permission UNIQUE (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role_id ON identity.role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON identity.role_permissions(permission_id);

-- User roles table
CREATE TABLE IF NOT EXISTS identity.user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES identity.roles(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_role_tenant UNIQUE (user_id, role_id, tenant_id)
);

CREATE INDEX idx_user_roles_user_id ON identity.user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON identity.user_roles(role_id);
CREATE INDEX idx_user_roles_tenant_id ON identity.user_roles(tenant_id);

-- User invitations table
CREATE TABLE IF NOT EXISTS identity.user_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    invited_by UUID NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    role_ids UUID[],
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_invitations_tenant_id ON identity.user_invitations(tenant_id);
CREATE INDEX idx_user_invitations_token ON identity.user_invitations(token);
CREATE INDEX idx_user_invitations_email ON identity.user_invitations(email);
