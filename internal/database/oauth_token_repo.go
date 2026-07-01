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

func GetOAuthTokenByAthleteID (db *sql.DB, athleteID int64)(*OAuthToken, error){
	query := `

	SELECT
		id,
		athlete_id,
		access_token,
		refresh_token,
		expires_at,
		created_at
	FROM oauth_tokens
	WHERE athlete_id = $1
	`
	var token OAuthToken

	err := db.QueryRow(
		query,
		athleteID,
	).Scan(
		&token.ID,
		&token.AthleteID,
		&token.AccessToken,
		&token.RefreshToken,
		&token.ExpiresAt,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func UpdateOAuthToken (db *sql.DB, token *OAuthToken ) error {
	query := `
	UPDATE oauth_tokens
	SET
		access_token = $1,
		refresh_token = $2,
		expires_at = $3
	WHERE athlete_id = $4
	`

	_, err := db.Exec(
		query,
		token.AccessToken,
		token.RefreshToken,
		token.ExpiresAt,
		token.AthleteID,
		)

	return err

}
