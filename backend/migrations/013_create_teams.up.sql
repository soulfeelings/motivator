CREATE TABLE teams (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id  UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    color       VARCHAR(7) DEFAULT '#8b5cf6',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(company_id, name)
);

CREATE TABLE team_members (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id       UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    membership_id UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    joined_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(team_id, membership_id)
);

CREATE TYPE team_battle_status AS ENUM ('pending', 'active', 'completed');

CREATE TABLE team_battles (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id   UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    team_a_id    UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    team_b_id    UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    metric       VARCHAR(100) NOT NULL,
    target       INTEGER NOT NULL DEFAULT 0,
    team_a_score INTEGER NOT NULL DEFAULT 0,
    team_b_score INTEGER NOT NULL DEFAULT 0,
    status       team_battle_status NOT NULL DEFAULT 'pending',
    winner_id    UUID REFERENCES teams(id),
    xp_reward    INTEGER NOT NULL DEFAULT 100,
    coin_reward  INTEGER NOT NULL DEFAULT 50,
    deadline     TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '7 days',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_teams_company_id ON teams(company_id);
CREATE INDEX idx_team_members_team_id ON team_members(team_id);
CREATE INDEX idx_team_members_membership_id ON team_members(membership_id);
CREATE INDEX idx_team_battles_company_id ON team_battles(company_id);
