-- Create collaboration schema
CREATE SCHEMA IF NOT EXISTS collaboration;

-- Chat threads table
CREATE TABLE IF NOT EXISTS collaboration.chat_threads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    title VARCHAR(255),
    thread_type VARCHAR(50) NOT NULL,
    reference_id UUID,
    is_archived BOOLEAN DEFAULT false,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_chat_threads_tenant_id ON collaboration.chat_threads(tenant_id);
CREATE INDEX idx_chat_threads_thread_type ON collaboration.chat_threads(thread_type);
CREATE INDEX idx_chat_threads_reference_id ON collaboration.chat_threads(reference_id);
CREATE INDEX idx_chat_threads_deleted_at ON collaboration.chat_threads(deleted_at);

-- Thread participants table
CREATE TABLE IF NOT EXISTS collaboration.thread_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    thread_id UUID NOT NULL REFERENCES collaboration.chat_threads(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    role VARCHAR(50),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_thread_user UNIQUE (thread_id, user_id)
);

CREATE INDEX idx_thread_participants_thread_id ON collaboration.thread_participants(thread_id);
CREATE INDEX idx_thread_participants_user_id ON collaboration.thread_participants(user_id);
CREATE INDEX idx_thread_participants_tenant_id ON collaboration.thread_participants(tenant_id);

-- Chat messages table
CREATE TABLE IF NOT EXISTS collaboration.chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    thread_id UUID NOT NULL REFERENCES collaboration.chat_threads(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL,
    message TEXT NOT NULL,
    message_type VARCHAR(50) DEFAULT 'text',
    is_read BOOLEAN DEFAULT false,
    is_edited BOOLEAN DEFAULT false,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_chat_messages_thread_id ON collaboration.chat_messages(thread_id);
CREATE INDEX idx_chat_messages_sender_id ON collaboration.chat_messages(sender_id);
CREATE INDEX idx_chat_messages_created_at ON collaboration.chat_messages(created_at);
CREATE INDEX idx_chat_messages_deleted_at ON collaboration.chat_messages(deleted_at);

-- Message files table
CREATE TABLE IF NOT EXISTS collaboration.message_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL REFERENCES collaboration.chat_messages(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_url TEXT NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_message_files_message_id ON collaboration.message_files(message_id);

-- Disputes table
CREATE TABLE IF NOT EXISTS collaboration.disputes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    order_id UUID NOT NULL,
    dispute_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'open',
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    raised_by UUID NOT NULL,
    resolved_by UUID,
    resolution TEXT,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_disputes_tenant_id ON collaboration.disputes(tenant_id);
CREATE INDEX idx_disputes_order_id ON collaboration.disputes(order_id);
CREATE INDEX idx_disputes_status ON collaboration.disputes(status);
CREATE INDEX idx_disputes_deleted_at ON collaboration.disputes(deleted_at);

-- Ratings table
CREATE TABLE IF NOT EXISTS collaboration.ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    order_id UUID NOT NULL,
    rated_by UUID NOT NULL,
    rated_entity_type VARCHAR(50) NOT NULL,
    rated_entity_id UUID NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    is_verified BOOLEAN DEFAULT false,
    is_moderated BOOLEAN DEFAULT false,
    moderated_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_order_rated_by_entity UNIQUE (order_id, rated_by, rated_entity_type, rated_entity_id)
);

CREATE INDEX idx_ratings_tenant_id ON collaboration.ratings(tenant_id);
CREATE INDEX idx_ratings_order_id ON collaboration.ratings(order_id);
CREATE INDEX idx_ratings_rated_entity ON collaboration.ratings(rated_entity_type, rated_entity_id);
CREATE INDEX idx_ratings_rating ON collaboration.ratings(rating);
