// This file synchronizes recent Strava activities into PostgreSQL.
//this centralizes token retrieval, pagination, mapping, and database saving for future OAuth and webhook use

package strava

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
)

type SyncService struct {
	db           *sql.DB
	clientID     string
	clientSecret string
}

func NewSyncService(
	db *sql.DB,
	clientID string,
	clientSecret string,
) *SyncService {
	return &SyncService{
		db:           db,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// this is the function used when a user logs in for the first time, we fetch 90 days of activities from Strava
// Store 90 days of activities in database
func (s *SyncService) SyncActivities(athleteID int64, after time.Time) error {

	stravaAthleteID, err := database.GetStravaAthleteIDByID(
		s.db,
		athleteID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to find Strava athlete ID: %w",
			err,
		)
	}

	accessToken, err := GetValidAccessToken(
		s.db,
		s.clientID,
		s.clientSecret,
		stravaAthleteID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to obtain valid Strava access token: %w",
			err,
		)
	}

	//calls strava api and gets the activities
	activities, err := GetActivitiesAfter(
		accessToken,
		after,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to fetch Strava activities: %w",
			err,
		)
	}

	dbActivities := make(
		[]database.Activity,
		0,
		len(activities),
	)

	for _, activity := range activities {
		dbActivity, err := MapActivityToDatabase(
			activity,
			athleteID,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to map Strava activity %d: %w",
				activity.ID,
				err,
			)
		}

		dbActivities = append(
			dbActivities,
			dbActivity,
		)
	}

	if len(dbActivities) == 0 {
		return nil
	}

	if err := database.SaveActivities(
		s.db,
		dbActivities,
	); err != nil {
		return fmt.Errorf(
			"failed to save activities: %w",
			err,
		)
	}

	return nil
}

//this is the fucntion used when a webhook triggeres syncing an activity, ex: update, delete or create

func (s *SyncService) SyncActivity(athleteID int64, stravaActivityID int64) error {

	stravaAthleteID, err :=
		database.GetStravaAthleteIDByID(
			s.db,
			athleteID,
		)
	if err != nil {
		return fmt.Errorf(
			"failed to find Strava athlete ID: %w",
			err,
		)
	}

	accessToken, err := GetValidAccessToken(
		s.db,
		s.clientID,
		s.clientSecret,
		stravaAthleteID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to obtain valid Strava access token: %w",
			err,
		)
	}

	//call Strava API to get the activity details using the activity ID
	activity, err := GetActivity(
		accessToken,
		stravaActivityID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to fetch activity: %w",
			err,
		)
	}

	dbActivity, err := MapActivityToDatabase(
		activity,
		athleteID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to map activity: %w",
			err,
		)
	}

	if err := database.SaveActivities(
		s.db,
		[]database.Activity{dbActivity},
	); err != nil {
		return fmt.Errorf(
			"failed to save activity: %w",
			err,
		)
	}

	return nil
}
