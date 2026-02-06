-- Create notification schema
CREATE SCHEMA IF NOT EXISTS notification;

-- Notification templates table
CREATE TABLE IF NOT EXISTS notification.templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    subject VARCHAR(255),
    body TEXT NOT NULL,
    body_html TEXT,
    channel VARCHAR(50) NOT NULL,
    event_type VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_templates_code ON notification.templates(code);
CREATE INDEX idx_templates_event_type ON notification.templates(event_type);
CREATE INDEX idx_templates_channel ON notification.templates(channel);

-- Notification preferences table
CREATE TABLE IF NOT EXISTS notification.preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    channel VARCHAR(50) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    is_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_channel_event UNIQUE (user_id, channel, event_type)
);

CREATE INDEX idx_preferences_tenant_id ON notification.preferences(tenant_id);
CREATE INDEX idx_preferences_user_id ON notification.preferences(user_id);
CREATE INDEX idx_preferences_channel ON notification.preferences(channel);
CREATE INDEX idx_preferences_event_type ON notification.preferences(event_type);

-- Notifications table
CREATE TABLE IF NOT EXISTS notification.notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    template_id UUID REFERENCES notification.templates(id),
    channel VARCHAR(50) NOT NULL,
    type VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    status VARCHAR(50) DEFAULT 'pending',
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    sent_at TIMESTAMP,
    failed_at TIMESTAMP,
    error TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_notifications_tenant_id ON notification.notifications(tenant_id);
CREATE INDEX idx_notifications_user_id ON notification.notifications(user_id);
CREATE INDEX idx_notifications_template_id ON notification.notifications(template_id);
CREATE INDEX idx_notifications_channel ON notification.notifications(channel);
CREATE INDEX idx_notifications_type ON notification.notifications(type);
CREATE INDEX idx_notifications_status ON notification.notifications(status);
CREATE INDEX idx_notifications_is_read ON notification.notifications(is_read);
CREATE INDEX idx_notifications_deleted_at ON notification.notifications(deleted_at);
