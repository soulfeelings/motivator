CREATE TYPE quest_status AS ENUM ('draft', 'active', 'voting', 'revealed', 'completed');

CREATE TABLE quests (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id    UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL DEFAULT 'Secret Motivator',
    description   TEXT,
    status        quest_status NOT NULL DEFAULT 'draft',
    xp_reward     INTEGER NOT NULL DEFAULT 25,
    coin_reward   INTEGER NOT NULL DEFAULT 10,
    bonus_xp      INTEGER NOT NULL DEFAULT 50,
    bonus_coins   INTEGER NOT NULL DEFAULT 25,
    deadline      TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '3 days',
    reveal_at     TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at  TIMESTAMPTZ
);

CREATE TABLE quest_pairs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quest_id      UUID NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    sender_id     UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    receiver_id   UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    message       TEXT,
    sent_at       TIMESTAMPTZ,
    UNIQUE(quest_id, sender_id)
);

CREATE TABLE quest_votes (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quest_id      UUID NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    pair_id       UUID NOT NULL REFERENCES quest_pairs(id) ON DELETE CASCADE,
    voter_id      UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    UNIQUE(quest_id, voter_id)
);

CREATE INDEX idx_quests_company_id ON quests(company_id);
CREATE INDEX idx_quest_pairs_quest_id ON quest_pairs(quest_id);
CREATE INDEX idx_quest_pairs_receiver_id ON quest_pairs(receiver_id);
CREATE INDEX idx_quest_votes_pair_id ON quest_votes(pair_id);
