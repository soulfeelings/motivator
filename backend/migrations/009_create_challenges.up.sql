CREATE TYPE challenge_status AS ENUM ('pending', 'active', 'completed', 'declined', 'cancelled');

CREATE TABLE challenges (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id     UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    challenger_id  UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    opponent_id    UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    metric         VARCHAR(100) NOT NULL,
    target         INTEGER NOT NULL,
    wager          INTEGER NOT NULL DEFAULT 0,
    status         challenge_status NOT NULL DEFAULT 'pending',
    challenger_score INTEGER NOT NULL DEFAULT 0,
    opponent_score   INTEGER NOT NULL DEFAULT 0,
    winner_id      UUID REFERENCES memberships(id),
    xp_reward      INTEGER NOT NULL DEFAULT 50,
    deadline       TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '7 days',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at   TIMESTAMPTZ
);

CREATE INDEX idx_challenges_company_id ON challenges(company_id);
CREATE INDEX idx_challenges_challenger_id ON challenges(challenger_id);
CREATE INDEX idx_challenges_opponent_id ON challenges(opponent_id);
CREATE INDEX idx_challenges_status ON challenges(status);
