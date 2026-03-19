CREATE TABLE device_tokens (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    membership_id UUID NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    token         TEXT NOT NULL,
    platform      VARCHAR(20) NOT NULL DEFAULT 'android',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(membership_id, token)
);

CREATE INDEX idx_device_tokens_membership_id ON device_tokens(membership_id);
