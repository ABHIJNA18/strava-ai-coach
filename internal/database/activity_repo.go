package database

import (
	"database/sql"
	"time"
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

func GetActivityStats(db *sql.DB, athleteID int64) (*ActivityStats, error) {
	query := `
		SELECT 
			COUNT(*) AS total_activities,

			COUNT (*) FILTER (
				WHERE sport_type = 'Run'
				)AS total_runs,

			COUNT (*) FILTER (
				WHERE sport_type = 'Hike'
				)AS total_hikes,

			COUNT(*) FILTER (
				WHERE sport_type = 'WeightTraining'
				) AS total_weight_training,

			COALESCE(SUM(distance), 0),

			COALESCE (
				SUM (distance)
				FILTER (WHERE sport_type = 'Run'),
				0
			),

			COALESCE (
				SUM(distance)
				FILTER (WHERE sport_type = 'Hike'),
				0
			),

			COALESCE(
				AVG(average_heartrate)
				FILTER (WHERE average_heartrate > 0),
				0
			) AS average_heart_rate,

			COALESCE(
				AVG(average_heartrate)
				FILTER (
					WHERE sport_type = 'Run'
					AND average_heartrate > 0
				),
				0
			) AS average_run_heart_rate,


			COALESCE(
				MAX(distance)
				FILTER (WHERE sport_type='Run'),
				0
			),

			COALESCE(
				MAX(distance)
				FILTER (WHERE sport_type='Hike'),
				0
			)

		FROM activities
		WHERE athlete_id = $1
	
	`

	var stats ActivityStats

	err := db.QueryRow(
		query,
		athleteID,
	).Scan(
		&stats.TotalActivities,
		&stats.TotalRuns,
		&stats.TotalHikes,
		&stats.TotalWeightTraining,
		&stats.TotalDistance,
		&stats.TotalRunDistance,
		&stats.TotalHikeDistance,
		&stats.AverageHeartRate,
		&stats.AverageRunHeartRate,
		&stats.LongestRunDistance,
		&stats.LongestHikeDistance,
	)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func GetActivitiesByAthleteID(
	db *sql.DB,
	athleteID int64,
) ([]Activity, error) {

	query := `
	SELECT
		id,
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
		start_date_local,
		created_at
	FROM activities
	WHERE athlete_id = $1
	ORDER BY start_date DESC
	`

	rows, err := db.Query(
		query,
		athleteID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var activities []Activity

	for rows.Next() {

		var activity Activity

		err := rows.Scan(
			&activity.ID,
			&activity.StravaActivityID,
			&activity.AthleteID,
			&activity.Name,
			&activity.Type,
			&activity.SportType,
			&activity.Distance,
			&activity.MovingTime,
			&activity.ElapsedTime,
			&activity.TotalElevationGain,
			&activity.AverageSpeed,
			&activity.MaxSpeed,
			&activity.AverageHeartrate,
			&activity.MaxHeartrate,
			&activity.AverageCadence,
			&activity.AverageWatts,
			&activity.MaxWatts,
			&activity.WeightedAverageWatts,
			&activity.Kilojoules,
			&activity.SufferScore,
			&activity.DeviceName,
			&activity.StartDate,
			&activity.StartDateLocal,
			&activity.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		activities = append(
			activities,
			activity,
		)
	}

	return activities, nil
}

func GetActivitiesByType(

	db *sql.DB,
	athleteID int64,
	sportType string,

) ([]Activity, error) {

	query := `
	SELECT
		id,
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
		start_date_local,
		created_at
	FROM activities
	WHERE athlete_id = $1
	AND sport_type = $2
	ORDER BY start_date DESC
	`

	rows, err := db.Query(query, athleteID, sportType)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//iterate through rows and fill in each activity and append to activities

	var activities []Activity

	for rows.Next() {
		var activity Activity

		err := rows.Scan(
			&activity.ID,
			&activity.StravaActivityID,
			&activity.AthleteID,
			&activity.Name,
			&activity.Type,
			&activity.SportType,
			&activity.Distance,
			&activity.MovingTime,
			&activity.ElapsedTime,
			&activity.TotalElevationGain,
			&activity.AverageSpeed,
			&activity.MaxSpeed,
			&activity.AverageHeartrate,
			&activity.MaxHeartrate,
			&activity.AverageCadence,
			&activity.AverageWatts,
			&activity.MaxWatts,
			&activity.WeightedAverageWatts,
			&activity.Kilojoules,
			&activity.SufferScore,
			&activity.DeviceName,
			&activity.StartDate,
			&activity.StartDateLocal,
			&activity.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}
	return activities, nil
}

func GetRecentActivitiesByType(
	db *sql.DB,
	athleteID int64,
	sportType string,
	limit int,
) ([]Activity, error) {

	query := `
	SELECT
		id,
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
		start_date_local,
		created_at
	FROM activities
	WHERE athlete_id = $1
	AND sport_type = $2
	ORDER BY start_date DESC
	LIMIT $3
	`
	rows, err := db.Query(query, athleteID, sportType, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//iterate through rows and fill in each activity and append to activities

	var activities []Activity

	for rows.Next() {
		var activity Activity

		err := rows.Scan(
			&activity.ID,
			&activity.StravaActivityID,
			&activity.AthleteID,
			&activity.Name,
			&activity.Type,
			&activity.SportType,
			&activity.Distance,
			&activity.MovingTime,
			&activity.ElapsedTime,
			&activity.TotalElevationGain,
			&activity.AverageSpeed,
			&activity.MaxSpeed,
			&activity.AverageHeartrate,
			&activity.MaxHeartrate,
			&activity.AverageCadence,
			&activity.AverageWatts,
			&activity.MaxWatts,
			&activity.WeightedAverageWatts,
			&activity.Kilojoules,
			&activity.SufferScore,
			&activity.DeviceName,
			&activity.StartDate,
			&activity.StartDateLocal,
			&activity.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}
	return activities, nil
}

func GetActivitiesByTypeSince(
	db *sql.DB,
	athleteID int64,
	sportType string,
	since time.Time,
) ([]Activity, error) {

	query := `
		SELECT
			id,
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
			start_date_local,
			created_at
		FROM activities
		WHERE athlete_id = $1
		AND sport_type = $2
		AND start_date >= $3
		ORDER BY start_date ASC
		`
	rows, err := db.Query(query, athleteID, sportType, since)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []Activity

	for rows.Next() {
		var activity Activity

		err := rows.Scan(
			&activity.ID,
			&activity.StravaActivityID,
			&activity.AthleteID,
			&activity.Name,
			&activity.Type,
			&activity.SportType,
			&activity.Distance,
			&activity.MovingTime,
			&activity.ElapsedTime,
			&activity.TotalElevationGain,
			&activity.AverageSpeed,
			&activity.MaxSpeed,
			&activity.AverageHeartrate,
			&activity.MaxHeartrate,
			&activity.AverageCadence,
			&activity.AverageWatts,
			&activity.MaxWatts,
			&activity.WeightedAverageWatts,
			&activity.Kilojoules,
			&activity.SufferScore,
			&activity.DeviceName,
			&activity.StartDate,
			&activity.StartDateLocal,
			&activity.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}
	return activities, nil
}

func GetRecentActivities(
	db *sql.DB,
	athleteID int64,
	limit int,
) ([]Activity, error) {

	query := `
	SELECT
		id,
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
		start_date_local,
		created_at
	FROM activities
	WHERE athlete_id = $1
	ORDER BY start_date DESC
	LIMIT $2
	`
	rows, err := db.Query(query, athleteID, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//iterate through rows and fill in each activity and append to activities

	var activities []Activity

	for rows.Next() {
		var activity Activity

		err := rows.Scan(
			&activity.ID,
			&activity.StravaActivityID,
			&activity.AthleteID,
			&activity.Name,
			&activity.Type,
			&activity.SportType,
			&activity.Distance,
			&activity.MovingTime,
			&activity.ElapsedTime,
			&activity.TotalElevationGain,
			&activity.AverageSpeed,
			&activity.MaxSpeed,
			&activity.AverageHeartrate,
			&activity.MaxHeartrate,
			&activity.AverageCadence,
			&activity.AverageWatts,
			&activity.MaxWatts,
			&activity.WeightedAverageWatts,
			&activity.Kilojoules,
			&activity.SufferScore,
			&activity.DeviceName,
			&activity.StartDate,
			&activity.StartDateLocal,
			&activity.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}
	return activities, nil
}

func GetTopSportSince(db *sql.DB, athleteID int64, since time.Time) ([]TopSport, error) {

	query := `
		WITH sport_counts AS (
			SELECT
				sport_type,
				COUNT(*)::int AS activity_count
			FROM activities
			WHERE athlete_id = $1
			AND start_date>= $2
			GROUP BY sport_type
		)
		SELECT
			sport_type,
			activity_count
		FROM sport_counts
		WHERE activity_count = (
			SELECT MAX(activity_count)
			FROM sport_counts
		)ORDER BY sport_type

	`

	rows, err := db.Query(
		query,
		athleteID,
		since,
	)

	if err != nil {

		return nil, err
	}

	defer rows.Close()

	var topSports []TopSport

	for rows.Next() {
		var sport TopSport

		err := rows.Scan(
			&sport.Sport,
			&sport.Count,
		)

		if err != nil {

			return nil, err
		}

		topSports = append(topSports, sport)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topSports, nil
}
