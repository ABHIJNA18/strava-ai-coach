package database

import "time"

type OAuthToken struct {
	ID           int64
	AthleteID    int64
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
	CreatedAt    time.Time
}