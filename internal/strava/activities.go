package strava

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// contains code to fetch activities of an athlete
func GetActivities(accessToken string) ([]Activity, error) {
	req, err := http.NewRequest("GET", "https://www.strava.com/api/v3/athlete/activities?per_page=10", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := stravaHTTPClient

	resp, err := client.Do(req)

	//error handling
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get activities. Status code : %s", resp.Status)
	}

	defer resp.Body.Close()

	var activities []Activity
	err = json.NewDecoder(resp.Body).Decode(&activities)

	if err != nil {
		return nil, err
	}

	fmt.Println("========== ACTIVITIES ==========")
	//fmt.Printf("Fetched activities %+v ", activities)
	for _, activity := range activities {
		fmt.Printf(" ID : %d | Name : %s | Type : %s | Distance : %.2f km | Start Date : %s \n",
			activity.ID, activity.Name, activity.Type, activity.Distance/1000, activity.StartDate)
	}
	return activities, nil

}
