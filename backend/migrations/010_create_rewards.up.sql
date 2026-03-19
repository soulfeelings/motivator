CREATE TYPE redemption_status AS ENUM ('pending', 'approved', 'fulfilled', 'rejected');

CREATE TABLE rewards (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id  UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    cost_coins  INTEGER NOT NULL,
    stock       INTEGER,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE redemptions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    membership_id UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    reward_id     UUID NOT NULL REFERENCES rewards(id) ON DELETE CASCADE,
    coins_spent   INTEGER NOT NULL,
    status        redemption_status NOT NULL DEFAULT 'pending',
    redeemed_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    fulfilled_at  TIMESTAMPTZ
);

CREATE INDEX idx_rewards_company_id ON rewards(company_id);
CREATE INDEX idx_redemptions_membership_id ON redemptions(membership_id);
CREATE INDEX idx_redemptions_reward_id ON redemptions(reward_id);
