CREATE TABLE deployments (
    id UUID PRIMARY KEY,
    service_name TEXT NOT NULL,
    version TEXT NOT NULL,
    environment TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_deployments_service_name ON deployments (service_name);
CREATE INDEX idx_deployments_environment ON deployments (environment);
