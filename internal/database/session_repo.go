package database

import (
	"database/sql"
	"fmt"
)

func CreateSession(db *sql.DB, session Session) (int64, error) {

	query := `
		INSERT INTO sessions (
			athlete_id,
			token_hash
		)
		VALUES ($1, $2)
		RETURNING id
	`

	var sessionID int64

	err := db.QueryRow(
		query,
		session.AthleteID,
		session.TokenHash,
	).Scan(&sessionID)

	if err != nil {
		fmt.Println("Error creating session")
		return 0, err
	}

	return sessionID, nil

}

func GetSessionByTokenHash(
	db *sql.DB,
	tokenHash string,
) (*Session, error) {
	query := `
		SELECT
			id,
			athlete_id,
			token_hash,
			created_at,
			last_seen_at,
			revoked_at
		FROM sessions
		WHERE token_hash = $1
	`

	var session Session
	var revokedAt sql.NullTime //to read the database value, and tells whether it exists and what the time is.

	err := db.QueryRow(
		query,
		tokenHash,
	).Scan(
		&session.ID,
		&session.AthleteID,
		&session.TokenHash,
		&session.CreatedAt,
		&session.LastSeenAt,
		&revokedAt,
	)

	if err != nil {
		return nil, err
	}

	// If this session has actually been revoked, store the revocation time in my Go Session object.

	if revokedAt.Valid {
		session.RevokedAt = &revokedAt.Time
	}

	return &session, nil
}

//simply updates the last_seen_at 
// keep it as useful metadata without using it to expire the credential.

func TouchSession(db *sql.DB, sessionID int64) error {

	query := `

		UPDATE sessions
		SET last_seen_at = NOW()
		WHERE id = $1
		AND revoked_at IS NULL 
	`
	_, err := db.Exec(
		query,
		sessionID,
	)

	return err

}

func RevokeSession(db *sql.DB, sessionID int64) error {

	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE id = $1
		AND revoked_at IS NULL
	`
	_, err := db.Exec(
		query,
		sessionID,
	)

	return err
}

func RevokeAllSessionsForAthlete(
	db *sql.DB,
	athleteID int64,
) error {
	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE athlete_id = $1
		AND revoked_at IS NULL
	`

	_, err := db.Exec(
		query,
		athleteID,
	)

	return err
}
