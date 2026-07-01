CREATE TABLE activities (

    id BIGSERIAL PRIMARY KEY,
    strava_activity_id BIGINT UNIQUE NOT NULL,
    athlete_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    sport_type TEXT,
    distance DOUBLE PRECISION,

    moving_time INTEGER,
    elapsed_time INTEGER,

    total_elevation_gain DOUBLE PRECISION,

    average_speed DOUBLE PRECISION,
    max_speed DOUBLE PRECISION,

    average_heartrate DOUBLE PRECISION,
    max_heartrate DOUBLE PRECISION,

    average_cadence DOUBLE PRECISION,

    average_watts DOUBLE PRECISION,
    max_watts DOUBLE PRECISION,
    weighted_average_watts DOUBLE PRECISION,

    kilojoules DOUBLE PRECISION,

    suffer_score DOUBLE PRECISION,

    device_name TEXT,

    start_date TIMESTAMPTZ,
    start_date_local TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_activity_athlete
        FOREIGN KEY (athlete_id)
        REFERENCES athletes(id)
        ON DELETE CASCADE
);
