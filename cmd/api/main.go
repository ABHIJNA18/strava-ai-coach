package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ABHIJNA18/strava-ai-coach/internal/coach"
	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
	"github.com/ABHIJNA18/strava-ai-coach/internal/handlers"
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

	// Connect to Python Coach Service
	coachClient, coachConn, err := coach.NewClient("localhost:50051")
	if err != nil {
		panic(err)
	}
	defer coachConn.Close()

	coachService := coach.NewService(db, coachClient)
	coachHandler := handlers.NewCoachHandler(coachService)

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

	//auth and activity handlers
	activityHandler := handlers.NewActivityHandler(db)
	authHandler := handlers.NewAuthHandler(db, clientID, clientSecret)

	//Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Strava AI Coach running")
	})

	//==== AUTH ENDPOINTS  =====

	http.HandleFunc("/login", authHandler.Login)
	//callback endpoint
	http.HandleFunc("/oauth/callback", authHandler.Callback)

	//==== ACTIVITY ENDPOINTS =====

	http.HandleFunc("/stats", activityHandler.GetStats)
	http.HandleFunc("/activities", activityHandler.GetActivities)
	http.HandleFunc("/activities/runs", activityHandler.GetRuns)
	http.HandleFunc("/activities/hikes", activityHandler.GetHikes)
	http.HandleFunc("/activities/weight-training", activityHandler.GetWeightTraining)
	http.HandleFunc("/activities/recent", activityHandler.GetRecentActivities)
	http.HandleFunc("/activities/recent/runs", activityHandler.GetRecentRuns)
	http.HandleFunc("/activities/recent/hikes", activityHandler.GetRecentHikes)
	http.HandleFunc("/activities/recent/weight-training", activityHandler.GetRecentWeightTraining)

	//==== COACH ENDPOINTS =====

	http.HandleFunc("/coach/report", coachHandler.GetReport)

	//==== FRONTEND =====

	fileServer := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fileServer)

	//==== FRONTEND DASHBOARD=====

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "./frontend/dashboard.html")

	})

	//==== START SERVER =====

	//verify server is running
	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)

}

/*test endpoint for tokem refresh logic
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
*/
