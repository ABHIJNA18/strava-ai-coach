// This file fetches all paginated Strava activities after a given timestamp.
//this preserves multi-page synchronization but limits results to the requested time window.

package strava

import (
	"fmt"
	"time"
)

func GetActivitiesAfter(
	accessToken string,
	after time.Time,
) ([]Activity, error) {
	var allActivities []Activity

	page := 1
	perPage := 100

	for {
		activities, err := GetActivitiesPage(
			accessToken,
			page,
			perPage,
			after,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to fetch activity page %d: %w",
				page,
				err,
			)
		}

		if len(activities) == 0 {
			break
		}

		allActivities = append(
			allActivities,
			activities...,
		)

		fmt.Printf(
			"Fetched page %d (%d activities)\n",
			page,
			len(activities),
		)

		page++
	}

	return allActivities, nil
}