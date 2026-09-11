-- Ensures each athlete has one OAuth token record.

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM oauth_tokens
        GROUP BY athlete_id
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION
            'Cannot add unique constraint: duplicate oauth_tokens.athlete_id values exist';
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'oauth_tokens_athlete_id_unique'
    ) THEN
        ALTER TABLE oauth_tokens
        ADD CONSTRAINT oauth_tokens_athlete_id_unique
        UNIQUE (athlete_id);
    END IF;
END
$$;