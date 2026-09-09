// This file provides centralized configuration for Strava API endpoints.

package strava

import (
	"os"
	"strings"
)

const defaultStravaAPIBaseURL = "https://www.strava.com/api/v3"

func StravaAPIBaseURL() string {
	configuredURL := strings.TrimRight(
		os.Getenv("STRAVA_API_BASE_URL"),
		"/",
	)

	if configuredURL == "" {
		return defaultStravaAPIBaseURL
	}

	return configuredURL
}