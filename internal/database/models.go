package database

import "time"

//tags are needed as APIs usually return: api formats like "total_distance" instead of TotalDistance

type Athlete struct {
	ID               int64
	StravaAthleteID int64
	Firstname       string
	Lastname        string
	CreatedAt       time.Time
}

type OAuthToken struct {
	ID           int64
	AthleteID    int64
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
	CreatedAt    time.Time
}

type Activity struct {
	ID               int64     `json:"id"`
	StravaActivityID int64     `json:"strava_activity_id"`
	AthleteID        int64     `json:"athlete_id"`

	Name             string    `json:"name"`
	Type             string    `json:"type"`
	SportType        string    `json:"sport_type"`

	Distance         float64   `json:"distance"`
	MovingTime       int       `json:"moving_time"`
	ElapsedTime      int       `json:"elapsed_time"`

	TotalElevationGain float64 `json:"total_elevation_gain"`

	AverageSpeed       float64 `json:"average_speed"`
	MaxSpeed           float64 `json:"max_speed"`

	AverageHeartrate   float64 `json:"average_heartrate"`
	MaxHeartrate       float64 `json:"max_heartrate"`

	AverageCadence     float64 `json:"average_cadence"`

	AverageWatts       float64 `json:"average_watts"`
	MaxWatts           float64 `json:"max_watts"`

	WeightedAverageWatts float64 `json:"weighted_average_watts"`

	Kilojoules         float64 `json:"kilojoules"`
	SufferScore        float64 `json:"suffer_score"`

	DeviceName         string `json:"device_name"`

	StartDate          time.Time `json:"start_date"`
	StartDateLocal     time.Time `json:"start_date_local"`

	CreatedAt          time.Time `json:"created_at"`
}

type ActivityStats struct {
	TotalActivities     int     `json:"total_activities"`
	TotalRuns           int     `json:"total_runs"`
	TotalHikes          int     `json:"total_hikes"`
	TotalWeightTraining int     `json:"total_weight_training"`

	TotalDistance       float64 `json:"total_distance"`
	TotalRunDistance    float64 `json:"total_run_distance"`
	TotalHikeDistance   float64 `json:"total_hike_distance"`

	AverageHeartRate    float64 `json:"average_heart_rate"`
	AverageRunHeartRate float64 `json:"average_run_heart_rate"`

	LongestRunDistance  float64 `json:"longest_run_distance"`
	LongestHikeDistance float64 `json:"longest_hike_distance"`
}