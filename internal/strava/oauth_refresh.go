package strava

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
)

func IsTokenExpired(expiresAt int64) bool {
	return time.Now().Unix() >= expiresAt
}

func RefreshAccessToken(clientID string, clientSecret string, refreshToken string) (*TokenResponse, error) {

	payload := RefreshTokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	}

	requestBody, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://www.strava.com/oauth/token",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := stravaHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"refresh failed. status code: %d",
			resp.StatusCode,
		)
	}

	defer resp.Body.Close()

	var tokenResponse TokenResponse

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, err
	}

	fmt.Println("Token refreshed successfully")
	return &tokenResponse, nil

}

// main refresh logic
func GetValidAccessToken(db *sql.DB, clientID string, clientSecret string, stravaAthleteID int64) (string, error) {

	//get athlete databse id
	athleteID, err := database.GetAthleteByStravaAthleteID(db, stravaAthleteID)
	if err != nil {
		return "", fmt.Errorf("Error fetching athlete id : %w", err)
	}

	//get oauth token details of the athlete using athlete id
	token, err := database.GetOAuthTokenByAthleteID(db, athleteID)
	if err != nil {
		fmt.Println("Error fetching token from athlete id ")
		return "", err
	}

	//check if the token is expired
	if !IsTokenExpired(token.ExpiresAt) {
		fmt.Println("Token has not expired")
		return token.AccessToken, nil
	}

	//if expired
	fmt.Println("Token has expired")
	refreshedToken, err := RefreshAccessToken(clientID, clientSecret, token.RefreshToken)

	if err != nil {
		fmt.Println("Error fetching refreshed token")
		return "", err
	}

	//save the refreshed token in db
	//edit the exisitng token

	token.AccessToken = refreshedToken.AccessToken
	token.RefreshToken = refreshedToken.RefreshToken
	token.ExpiresAt = refreshedToken.ExpiresAt

	err = database.UpdateOAuthToken(db, token)
	if err != nil {
		fmt.Println("Error updating Oauth token after refresh")
		return "", err
	}

	return token.AccessToken, nil
}
