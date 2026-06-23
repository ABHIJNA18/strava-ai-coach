package strava

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// reusable client
var stravaHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// contains code which exchanges authorisation code for an access token
func ExchangeTokenForCode(clientID string, clientSecret string, code string) (*TokenResponse, error) {
	payload := TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Code:         code,
		GrantType:    "authorization_code",
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// one method of post request using http.post 
	//another way is to use the resuable client and use http.NewRequest to create a new request and then use client.Do like GetAthlete()

	resp, err := http.Post("https://www.strava.com/oauth/token", "application/json", bytes.NewBuffer(requestBody))

	//Error Handling
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("failed to exchange token. Status code : %d", resp.StatusCode)
	}

	//always close the response body
	defer resp.Body.Close() 

	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse) //decode expects a memory address to write the decoded data into, so we pass a pointer to tokenresponse
	if err != nil {
		return nil, err
	}
	fmt.Println("Token response received")
	return &tokenResponse, nil

}

/*
defined a variable at the top instead of this fucntion. simpler

func returnClient() *http.Client {
	httpClient := &http.Client{}
	return httpClient
}
*/
