CREATE TABLE game_plans (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id  UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    flow_data   JSONB NOT NULL DEFAULT '{"nodes":[],"edges":[]}',
    is_active   BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(company_id, name)
);

CREATE INDEX idx_game_plans_company_id ON game_plans(company_id);
