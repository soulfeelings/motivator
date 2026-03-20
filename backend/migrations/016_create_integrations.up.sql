CREATE TABLE integrations (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id   UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    provider     VARCHAR(50) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    config       JSONB NOT NULL DEFAULT '{}',
    webhook_secret VARCHAR(64) NOT NULL,
    is_active    BOOLEAN NOT NULL DEFAULT true,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE integration_mappings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integration_id  UUID NOT NULL REFERENCES integrations(id) ON DELETE CASCADE,
    external_event  VARCHAR(255) NOT NULL,
    metric          VARCHAR(100) NOT NULL,
    user_field      VARCHAR(100) NOT NULL DEFAULT 'email',
    transform       JSONB NOT NULL DEFAULT '{"value": 1}',
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(integration_id, external_event)
);

CREATE TABLE integration_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integration_id  UUID NOT NULL REFERENCES integrations(id) ON DELETE CASCADE,
    external_event  VARCHAR(255) NOT NULL,
    raw_data        JSONB,
    metric          VARCHAR(100),
    user_email      VARCHAR(255),
    value           INTEGER NOT NULL DEFAULT 0,
    processed       BOOLEAN NOT NULL DEFAULT false,
    error           TEXT,
    received_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_integrations_company_id ON integrations(company_id);
CREATE INDEX idx_integration_mappings_integration_id ON integration_mappings(integration_id);
CREATE INDEX idx_integration_events_integration_id ON integration_events(integration_id);
CREATE INDEX idx_integration_events_received_at ON integration_events(received_at DESC);
