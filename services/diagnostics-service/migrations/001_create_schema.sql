-- Diagnostics schema for platform-wide debugging and monitoring

CREATE SCHEMA IF NOT EXISTS diagnostics;

-- Service heartbeats: track service health and last seen
CREATE TABLE IF NOT EXISTS diagnostics.service_heartbeats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(100) NOT NULL,
    instance_id VARCHAR(255) NOT NULL,
    last_seen_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'healthy',
    version VARCHAR(50),
    env VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(service_name, instance_id)
);

CREATE INDEX IF NOT EXISTS idx_heartbeats_service ON diagnostics.service_heartbeats(service_name);
CREATE INDEX IF NOT EXISTS idx_heartbeats_last_seen ON diagnostics.service_heartbeats(last_seen_at);
CREATE INDEX IF NOT EXISTS idx_heartbeats_status ON diagnostics.service_heartbeats(status);

-- Incidents: track errors and issues
CREATE TABLE IF NOT EXISTS diagnostics.incidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('INFO', 'WARN', 'ERROR', 'CRITICAL')),
    service_name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('DB', 'AUTH', 'EVENT', 'API', 'FILE', 'SEARCH', 'BILLING', 'LOGISTICS', 'OTHER')),
    error_code VARCHAR(50),
    title VARCHAR(255) NOT NULL,
    details_json JSONB,
    request_id VARCHAR(255),
    tenant_id UUID,
    user_id UUID,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolution_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_incidents_severity ON diagnostics.incidents(severity);
CREATE INDEX IF NOT EXISTS idx_incidents_service ON diagnostics.incidents(service_name);
CREATE INDEX IF NOT EXISTS idx_incidents_category ON diagnostics.incidents(category);
CREATE INDEX IF NOT EXISTS idx_incidents_occurred_at ON diagnostics.incidents(occurred_at);
CREATE INDEX IF NOT EXISTS idx_incidents_resolved_at ON diagnostics.incidents(resolved_at);
CREATE INDEX IF NOT EXISTS idx_incidents_tenant_id ON diagnostics.incidents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_incidents_request_id ON diagnostics.incidents(request_id);

-- Event failures: track event publish/consume failures
CREATE TABLE IF NOT EXISTS diagnostics.event_failures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_name VARCHAR(100) NOT NULL,
    direction VARCHAR(20) NOT NULL CHECK (direction IN ('PUBLISH', 'CONSUME')),
    service_name VARCHAR(100) NOT NULL,
    payload_json JSONB,
    error_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    last_attempt_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'FAILED' CHECK (status IN ('FAILED', 'RETRYING', 'RESOLVED')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_event_failures_event_name ON diagnostics.event_failures(event_name);
CREATE INDEX IF NOT EXISTS idx_event_failures_service ON diagnostics.event_failures(service_name);
CREATE INDEX IF NOT EXISTS idx_event_failures_status ON diagnostics.event_failures(status);
CREATE INDEX IF NOT EXISTS idx_event_failures_last_attempt ON diagnostics.event_failures(last_attempt_at);

-- Jobs: track background job status
CREATE TABLE IF NOT EXISTS diagnostics.jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_name VARCHAR(100) NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('RUNNING', 'FAILED', 'SUCCESS')),
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    metadata_json JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_jobs_name ON diagnostics.jobs(job_name);
CREATE INDEX IF NOT EXISTS idx_jobs_service ON diagnostics.jobs(service_name);
CREATE INDEX IF NOT EXISTS idx_jobs_status ON diagnostics.jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_started_at ON diagnostics.jobs(started_at);

-- API metrics: simple per-minute metrics
CREATE TABLE IF NOT EXISTS diagnostics.api_metrics_minute (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    minute_ts TIMESTAMP WITH TIME ZONE NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    route VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    count_total INTEGER NOT NULL DEFAULT 0,
    count_2xx INTEGER NOT NULL DEFAULT 0,
    count_4xx INTEGER NOT NULL DEFAULT 0,
    count_5xx INTEGER NOT NULL DEFAULT 0,
    p95_ms INTEGER,
    avg_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(minute_ts, service_name, route, method)
);

CREATE INDEX IF NOT EXISTS idx_metrics_minute_ts ON diagnostics.api_metrics_minute(minute_ts);
CREATE INDEX IF NOT EXISTS idx_metrics_service ON diagnostics.api_metrics_minute(service_name);
CREATE INDEX IF NOT EXISTS idx_metrics_route ON diagnostics.api_metrics_minute(route);
