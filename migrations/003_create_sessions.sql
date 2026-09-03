-- This migration creates persistent application authentication credentials.

CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,

    athlete_id BIGINT NOT NULL,

    token_hash TEXT UNIQUE NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    revoked_at TIMESTAMPTZ,

    CONSTRAINT fk_session_athlete
        FOREIGN KEY (athlete_id)
        REFERENCES athletes(id)
        ON DELETE CASCADE
);

-- UNIQUE(token_hash) already creates an index for token lookup.
CREATE INDEX idx_sessions_athlete_id
    ON sessions(athlete_id);