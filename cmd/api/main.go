package main

import (
	"fmt"
	"net/http"
	"os"

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

	//get env variables
	clientID := os.Getenv("STRAVA_CLIENT_ID")
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")

	if clientID == "" {
		panic("STRAVA_CLIENT_ID environment variable is not set")
	} else {
		fmt.Println("STRAVA_CLIENT_ID fetched")
	}
	if clientSecret == "" {
		panic("STRAVA_CLIENT_SECRET environment variable is not set")
	} else {
		fmt.Println("STRAVA_CLIENT_SECRET fetched")
	}

	//Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Strava AI Coach running")
	})

	//login endpoint
	http.HandleFunc("/login", strava.LoginHandler(clientID))

	//callback endpoint
	http.HandleFunc("/oauth/callback", strava.CallbackHandler(clientID, clientSecret))

	//verify server is running
	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
