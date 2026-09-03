// This file converts Strava API activities into database activity models.

package strava

import (
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
)

func MapActivityToDatabase(
	activity Activity,
	athleteID int64,
) (database.Activity, error) {
	startDate, err := time.Parse(
		time.RFC3339,
		activity.StartDate,
	)
	if err != nil {
		return database.Activity{}, err
	}

	startDateLocal, err := time.Parse(
		time.RFC3339,
		activity.StartDateLocal,
	)
	if err != nil {
		return database.Activity{}, err
	}

	return database.Activity{
		StravaActivityID:     activity.ID,
		AthleteID:            athleteID,
		Name:                 activity.Name,
		Type:                 activity.Type,
		SportType:            activity.SportType,
		Distance:             activity.Distance,
		MovingTime:           activity.MovingTime,
		ElapsedTime:          activity.ElapsedTime,
		TotalElevationGain:   activity.TotalElevationGain,
		AverageSpeed:         activity.AverageSpeed,
		MaxSpeed:             activity.MaxSpeed,
		AverageHeartrate:     activity.AverageHeartrate,
		MaxHeartrate:         activity.MaxHeartrate,
		AverageCadence:       activity.AverageCadence,
		AverageWatts:         activity.AverageWatts,
		MaxWatts:             activity.MaxWatts,
		WeightedAverageWatts: activity.WeightedAverageWatts,
		Kilojoules:           activity.Kilojoules,
		SufferScore:          activity.SufferScore,
		DeviceName:           activity.DeviceName,
		StartDate:            startDate,
		StartDateLocal:       startDateLocal,
	}, nil
}