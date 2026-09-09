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
		StravaAPIBaseURL() + "/athlete/activities",
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

//Fetches activity details from Strava API from the given activity ID
// Currently used when processing Strava Webhookn events

func GetActivity(
	accessToken string,
	stravaActivityID int64,
) (Activity, error) {
	url := fmt.Sprintf(
		"%s/activities/%d",
		StravaAPIBaseURL(),
		stravaActivityID,
	)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return Activity{}, err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+accessToken,
	)

	resp, err := stravaHTTPClient.Do(req)
	if err != nil {
		return Activity{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Activity{}, fmt.Errorf(
			"failed to fetch activity %d: %s",
			stravaActivityID,
			resp.Status,
		)
	}

	var activity Activity

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&activity); err != nil {
		return Activity{}, err
	}

	return activity, nil
}
