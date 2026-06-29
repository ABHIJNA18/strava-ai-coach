package database

import (
	"database/sql"
	"fmt"
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
