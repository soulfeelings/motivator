CREATE TABLE member_achievements (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    membership_id  UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    completed_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(membership_id, achievement_id)
);

CREATE INDEX idx_member_achievements_membership_id ON member_achievements(membership_id);
CREATE INDEX idx_member_achievements_achievement_id ON member_achievements(achievement_id);
