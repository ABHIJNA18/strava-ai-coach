package database

import (
	"database/sql"
	"fmt"
	"time"
)

// save athlete to DB and return the database ID for oauth to be stored in DB later on
func SaveAthlete(db *sql.DB, athlete Athlete) (int64, error) {

	//define the query
	query := `

		INSERT INTO athletes (
			strava_athlete_id,
			firstname,
			lastname
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (strava_athlete_id)
		DO UPDATE SET
			firstname = EXCLUDED.firstname,
			lastname = EXCLUDED.lastname
		RETURNING id
	`
	var athleteID int64

	err := db.QueryRow(
		query,
		athlete.StravaAthleteID,
		athlete.Firstname,
		athlete.Lastname,
	).Scan(&athleteID)

	if err != nil {
		return 0, err
	}

	fmt.Println("Athlete saved successfully to database")
	return athleteID, nil

}

func GetAthleteByStravaAthleteID ( db *sql.DB, stravaAthleteID int64 )(int64, error){
	query := `
		SELECT id
		FROM athletes
		WHERE strava_athlete_id = $1
		`
		var  athleteID int64

		err := db.QueryRow(
			query,
			stravaAthleteID,
			).Scan(&athleteID)
		
		if err != nil {
			return 0, err
		}
		return athleteID, nil
}

func GetStravaAthleteIDByID(
	db *sql.DB,
	athleteID int64,
) (int64, error) {
	query := `
		SELECT strava_athlete_id
		FROM athletes
		WHERE id = $1
	`

	var stravaAthleteID int64

	err := db.QueryRow(
		query,
		athleteID,
	).Scan(&stravaAthleteID)

	if err != nil {
		return 0, err
	}

	return stravaAthleteID, nil
}

func GetInitialSyncCompletedAt(
	db *sql.DB,
	athleteID int64,
) (*time.Time, error) {
	query := `
		SELECT initial_sync_completed_at
		FROM athletes
		WHERE id = $1
	`

	var completedAt sql.NullTime

	err := db.QueryRow(
		query,
		athleteID,
	).Scan(&completedAt)

	if err != nil {
		return nil, err
	}

	if !completedAt.Valid {
		return nil, nil
	}

	return &completedAt.Time, nil
}

func MarkInitialSyncCompleted(
	db *sql.DB,
	athleteID int64,
) error {
	query := `
		UPDATE athletes
		SET initial_sync_completed_at = NOW()
		WHERE id = $1
	`

	_, err := db.Exec(
		query,
		athleteID,
	)

	return err
}