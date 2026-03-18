CREATE TABLE invites (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id  UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    email       VARCHAR(255) NOT NULL,
    role        user_role NOT NULL DEFAULT 'employee',
    status      invite_status NOT NULL DEFAULT 'pending',
    invited_by  UUID NOT NULL,
    token       VARCHAR(64) UNIQUE NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '7 days',
    accepted_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(company_id, email)
);

CREATE INDEX idx_invites_token ON invites(token);
CREATE INDEX idx_invites_company_id ON invites(company_id);
