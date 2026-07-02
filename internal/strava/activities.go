package strava

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// contains code to fetch activities of an athlete using pagination
func GetActivitiesPage(accessToken string, page int, perPage int) ([]Activity, error) {

	//url needs to have page and per page parameters

	url := fmt.Sprintf("https://www.strava.com/api/v3/athlete/activities?page=%d&per_page=%d", page, perPage )

	req, err := http.NewRequest("GET", url, nil)
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
	
	/*for _, activity := range activities {

		/*fmt.Printf(" ID : %d | Name : %s | Type : %s | Distance : %.2f km | Start Date : %s | HR: %.1f | MaxHR: %.1f | Speed: %.2f | Suffer: %.1f \n",
			activity.ID, activity.Name, activity.Type, activity.Distance/1000, activity.StartDate, activity.AverageHeartrate, activity.MaxHeartrate, activity.AverageSpeed, activity.SufferScore)
	
	}*/
	fmt.Printf("\n Fetched %d activities :", len(activities))
	return activities, nil

}
