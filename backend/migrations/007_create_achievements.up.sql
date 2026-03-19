CREATE TYPE metric_operator AS ENUM ('gte', 'lte', 'eq', 'gt', 'lt');

CREATE TABLE achievements (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id  UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    metric      VARCHAR(100) NOT NULL,
    operator    metric_operator NOT NULL DEFAULT 'gte',
    threshold   INTEGER NOT NULL,
    badge_id    UUID REFERENCES badges(id) ON DELETE SET NULL,
    xp_reward   INTEGER NOT NULL DEFAULT 0,
    coin_reward INTEGER NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(company_id, name)
);

CREATE INDEX idx_achievements_company_id ON achievements(company_id);
CREATE INDEX idx_achievements_metric ON achievements(metric);
