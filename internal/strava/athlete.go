package strava

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//contains code which gets details of the athelete using the access token

func GetAthlete(accessToken string) (*Athlete, error) {

	//using http.newRequest and usinng a client sending a request 
	req, err := http.NewRequest("GET", "https://www.strava.com/api/v3/athlete", nil)
	if err != nil {
		fmt.Println("Error creating new request", err.Error())
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	//get reusable client 
	client := stravaHTTPClient
	resp, err := client.Do(req)
	
	//error handling
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("failed to get athlete data. Status code : %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var athlete Athlete
	err = json.NewDecoder(resp.Body).Decode(&athlete)
	if err != nil {
		return nil, err
	}

	fmt.Println("========== ATHLETE DATA ==========")
	fmt.Printf("%+v\n", athlete)

	return &athlete, nil

}
