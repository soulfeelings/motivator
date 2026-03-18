CREATE TABLE memberships (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL,
    company_id   UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    role         user_role NOT NULL DEFAULT 'employee',
    display_name VARCHAR(255),
    job_title    VARCHAR(255),
    is_active    BOOLEAN NOT NULL DEFAULT true,
    joined_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, company_id)
);

CREATE INDEX idx_memberships_user_id ON memberships(user_id);
CREATE INDEX idx_memberships_company_id ON memberships(company_id);
