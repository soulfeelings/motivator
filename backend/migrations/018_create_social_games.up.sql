CREATE TABLE social_games (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id    UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    game_type     VARCHAR(20) NOT NULL CHECK (game_type IN ('trivia', 'photo_challenge', 'two_truths')),
    name          VARCHAR(255) NOT NULL,
    description   TEXT,
    status        VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'voting', 'completed')),
    config        JSONB NOT NULL DEFAULT '{}',
    xp_reward     INTEGER NOT NULL DEFAULT 25,
    coin_reward   INTEGER NOT NULL DEFAULT 10,
    starts_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    ends_at       TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '24 hours',
    created_by    UUID NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at  TIMESTAMPTZ
);

CREATE TABLE social_game_questions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id       UUID NOT NULL REFERENCES social_games(id) ON DELETE CASCADE,
    question      TEXT NOT NULL,
    options       JSONB NOT NULL,
    correct_index INTEGER NOT NULL,
    sort_order    INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE social_game_answers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id         UUID NOT NULL REFERENCES social_games(id) ON DELETE CASCADE,
    question_id     UUID NOT NULL REFERENCES social_game_questions(id) ON DELETE CASCADE,
    member_id       UUID NOT NULL,
    selected_index  INTEGER NOT NULL,
    is_correct      BOOLEAN NOT NULL,
    answered_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(question_id, member_id)
);

CREATE TABLE social_game_submissions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id       UUID NOT NULL REFERENCES social_games(id) ON DELETE CASCADE,
    member_id     UUID NOT NULL,
    content       TEXT,
    statements    JSONB,
    lie_index     INTEGER,
    submitted_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(game_id, member_id)
);

CREATE TABLE social_game_votes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id         UUID NOT NULL REFERENCES social_games(id) ON DELETE CASCADE,
    submission_id   UUID NOT NULL REFERENCES social_game_submissions(id) ON DELETE CASCADE,
    voter_id        UUID NOT NULL,
    vote_value      INTEGER,
    voted_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(game_id, submission_id, voter_id)
);

CREATE INDEX idx_social_games_company ON social_games(company_id);
CREATE INDEX idx_social_game_questions_game ON social_game_questions(game_id);
CREATE INDEX idx_social_game_answers_game ON social_game_answers(game_id);
CREATE INDEX idx_social_game_submissions_game ON social_game_submissions(game_id);
CREATE INDEX idx_social_game_votes_submission ON social_game_votes(submission_id);
