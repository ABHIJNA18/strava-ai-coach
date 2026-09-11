package strava
// This file tests conversion from Strava activity models to database models.

import (
	"testing"
	"time"
)

func TestMapActivityToDatabase(t *testing.T) {
	input := Activity{
		ID:                   123,
		Name:                 "Morning Run",
		Type:                 "Run",
		SportType:            "Run",
		Distance:             5000.5,
		MovingTime:           1800,
		ElapsedTime:          1850,
		TotalElevationGain:   42.5,
		AverageSpeed:         2.78,
		MaxSpeed:             3.5,
		AverageHeartrate:     150,
		MaxHeartrate:         170,
		AverageCadence:       165,
		AverageWatts:         210,
		MaxWatts:             300,
		WeightedAverageWatts: 220,
		Kilojoules:           400,
		SufferScore:          80,
		DeviceName:           "Garmin",
		StartDate:            "2026-09-03T06:00:00Z",
		StartDateLocal:       "2026-09-03T07:00:00Z",
	}

	result, err := MapActivityToDatabase(
		input,
		42,
	)
	if err != nil {
		t.Fatal(err)
	}

	if result.StravaActivityID != 123 {
		t.Fatalf("unexpected Strava activity ID")
	}

	if result.AthleteID != 42 {
		t.Fatalf("expected athlete ID 42, got %d", result.AthleteID)
	}

	if result.Distance != input.Distance {
		t.Fatalf("distance was not copied")
	}

	if result.AverageCadence != input.AverageCadence {
		t.Fatalf("cadence was not copied")
	}

	expectedStart, _ := time.Parse(
		time.RFC3339,
		input.StartDate,
	)

	if !result.StartDate.Equal(expectedStart) {
		t.Fatalf("StartDate was not parsed correctly")
	}
}