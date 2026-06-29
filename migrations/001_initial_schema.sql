CREATE TABLE athletes (
    id BIGSERIAL PRIMARY KEY,
    strava_athlete_id BIGINT UNIQUE NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE oauth_tokens (
    id BIGSERIAL PRIMARY KEY,
    athlete_id BIGINT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_athlete
        FOREIGN KEY (athlete_id)
        REFERENCES athletes(id)
        ON DELETE CASCADE
);