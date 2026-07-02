package database

import "time"

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
	ID                     int64

	StravaActivityID       int64

	AthleteID              int64

	Name                   string
	Type                   string
	SportType              string

	Distance               float64

	MovingTime             int
	ElapsedTime            int

	TotalElevationGain     float64

	AverageSpeed           float64
	MaxSpeed               float64

	AverageHeartrate       float64
	MaxHeartrate           float64

	AverageCadence         float64

	AverageWatts           float64
	MaxWatts               float64
	WeightedAverageWatts   float64

	Kilojoules             float64

	SufferScore            float64

	DeviceName             string

	StartDate              time.Time
	StartDateLocal         time.Time

	CreatedAt              time.Time
}