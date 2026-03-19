CREATE TYPE tournament_status AS ENUM ('draft', 'registration', 'active', 'completed', 'cancelled');

CREATE TABLE tournaments (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id    UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    description   TEXT,
    season        VARCHAR(50),
    metric        VARCHAR(100) NOT NULL,
    prize_pool    INTEGER NOT NULL DEFAULT 0,
    xp_prizes     JSONB NOT NULL DEFAULT '[100, 50, 25]',
    coin_prizes   JSONB NOT NULL DEFAULT '[50, 25, 10]',
    status        tournament_status NOT NULL DEFAULT 'draft',
    max_participants INTEGER,
    starts_at     TIMESTAMPTZ NOT NULL,
    ends_at       TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at  TIMESTAMPTZ
);

CREATE TABLE tournament_participants (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tournament_id  UUID NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    membership_id  UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    score          INTEGER NOT NULL DEFAULT 0,
    rank           INTEGER,
    joined_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tournament_id, membership_id)
);

CREATE INDEX idx_tournaments_company_id ON tournaments(company_id);
CREATE INDEX idx_tournaments_status ON tournaments(status);
CREATE INDEX idx_tournament_participants_tournament_id ON tournament_participants(tournament_id);
CREATE INDEX idx_tournament_participants_membership_id ON tournament_participants(membership_id);
