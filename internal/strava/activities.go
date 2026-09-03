// This file fetches one paginated activity page from Strava.
//the Strava request must include the after timestamp while preserving pagination

package strava

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GetActivitiesPage(
	accessToken string,
	page int,
	perPage int,
	after time.Time,
) ([]Activity, error) {
	endpoint, err := url.Parse(
		"https://www.strava.com/api/v3/athlete/activities",
	)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("page", strconv.Itoa(page))
	query.Set("per_page", strconv.Itoa(perPage))
	query.Set("after", strconv.FormatInt(after.Unix(), 10))
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequest(
		http.MethodGet,
		endpoint.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+accessToken,
	)

	resp, err := stravaHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"failed to get activities: %s",
			resp.Status,
		)
	}

	var activities []Activity

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&activities); err != nil {
		return nil, err
	}

	return activities, nil
}