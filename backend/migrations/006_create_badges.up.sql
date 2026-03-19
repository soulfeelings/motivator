CREATE TABLE badges (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id  UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    icon_url    TEXT,
    xp_reward   INTEGER NOT NULL DEFAULT 0,
    coin_reward INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(company_id, name)
);

CREATE TABLE member_badges (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    membership_id UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    badge_id      UUID NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    awarded_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(membership_id, badge_id)
);

CREATE INDEX idx_badges_company_id ON badges(company_id);
CREATE INDEX idx_member_badges_membership_id ON member_badges(membership_id);
CREATE INDEX idx_member_badges_badge_id ON member_badges(badge_id);
