package database

import (
	"database/sql"
)

func SaveOauthToken(db *sql.DB, token OAuthToken) error {

	query := ` 
	INSERT INTO oauth_tokens (
	athlete_id,
	access_token,
	refresh_token,
	expires_at
	)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (athlete_id)
	DO UPDATE SET 
		access_token = EXCLUDED.access_token,
		refresh_token = EXCLUDED.refresh_token,
		expires_at = EXCLUDED.expires_at
	`

	_, err := db.Exec(
		query,
		token.AthleteID,
		token.AccessToken,
		token.RefreshToken,
		token.ExpiresAt,
	)

    return err
}
