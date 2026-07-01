package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
	"github.com/ABHIJNA18/strava-ai-coach/internal/strava"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	} else {
		fmt.Println(".env file loaded successfully")
	}

	//create connect to databse
	db, err := database.NewPostgresConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//get strava env variables
	clientID := os.Getenv("STRAVA_CLIENT_ID")
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")

	if clientID == "" {
		panic("STRAVA_CLIENT_ID environment variable is not set")
	} else {
		fmt.Println("Strava configuration loaded")
	}
	if clientSecret == "" {
		panic("STRAVA_CLIENT_SECRET environment variable is not set")
	} else {
		fmt.Println("Strava configuration loaded")
	}

	//Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Strava AI Coach running")
	})

	//test endpoint for tokem refresh logic
	http.HandleFunc("/test-token", func(w http.ResponseWriter, r *http.Request) {

		accessToken, err := strava.GetValidAccessToken(
			db,
			clientID,
			clientSecret,
			155503972,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("Access Token Retrieved Successfully\n%s", accessToken)
	})

	//login endpoint
	http.HandleFunc("/login", strava.LoginHandler(clientID))

	//callback endpoint
	http.HandleFunc("/oauth/callback", strava.CallbackHandler(clientID, clientSecret, db))

	//verify server is running
	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
