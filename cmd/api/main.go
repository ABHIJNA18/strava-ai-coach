package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	} else {
		fmt.Println(".env file loaded successfully")
	}

	client_id := os.Getenv("STRAVA_CLIENT_ID")
	if client_id == "" {
		panic("STRAVA_CLIENT_ID environment variable is not set")
	} else {
		fmt.Println("STRAVA_CLIENT_ID fetched")
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Strava AI Coach running")
	})

	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
