package database

import "time"


type Athlete struct {
	ID              int64
	StravaAthleteID int64
	Firstname       string
	Lastname        string
	CreatedAt       time.Time
}