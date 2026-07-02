package database

import (
	"database/sql"
)

func SaveActivities(db *sql.DB, activities []Activity) error {
	query := `
		INSERT INTO activities (

		strava_activity_id,
		athlete_id,
		name,
		type,
		sport_type,
		distance,
		moving_time,
		elapsed_time,
		total_elevation_gain,
		average_speed,
		max_speed,
		average_heartrate,
		max_heartrate,
		average_cadence,
		average_watts,
		max_watts,
		weighted_average_watts,
		kilojoules,
		suffer_score,
		device_name,
		start_date,
		start_date_local
	)
	VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,
		$12,$13,$14,$15,$16,$17,$18,$19,$20,
		$21,$22
	)
	ON CONFLICT (strava_activity_id)

	DO UPDATE SET
		name = EXCLUDED.name,
		type = EXCLUDED.type,
		sport_type = EXCLUDED.sport_type,
		distance = EXCLUDED.distance,
		moving_time = EXCLUDED.moving_time,
		elapsed_time = EXCLUDED.elapsed_time,
		total_elevation_gain = EXCLUDED.total_elevation_gain,
		average_speed = EXCLUDED.average_speed,
		max_speed = EXCLUDED.max_speed,
		average_heartrate = EXCLUDED.average_heartrate,
		max_heartrate = EXCLUDED.max_heartrate,
		average_cadence = EXCLUDED.average_cadence,
		average_watts = EXCLUDED.average_watts,
		max_watts = EXCLUDED.max_watts,
		weighted_average_watts = EXCLUDED.weighted_average_watts,
		kilojoules = EXCLUDED.kilojoules,
		suffer_score = EXCLUDED.suffer_score,
		device_name = EXCLUDED.device_name,
		start_date = EXCLUDED.start_date,
		start_date_local = EXCLUDED.start_date_local
	`

	//iterate iver aqctivities and insert into db

	for _, activity := range activities {
		_, err := db.Exec(
			query,
			activity.StravaActivityID,
			activity.AthleteID,
			activity.Name,
			activity.Type,
			activity.SportType,
			activity.Distance,
			activity.MovingTime,
			activity.ElapsedTime,
			activity.TotalElevationGain,
			activity.AverageSpeed,
			activity.MaxSpeed,
			activity.AverageHeartrate,
			activity.MaxHeartrate,
			activity.AverageCadence,
			activity.AverageWatts,
			activity.MaxWatts,
			activity.WeightedAverageWatts,
			activity.Kilojoules,
			activity.SufferScore,
			activity.DeviceName,
			activity.StartDate,
			activity.StartDateLocal,
		)
		if err != nil {
			return err
		}
	}
	return nil

}
