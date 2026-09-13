package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/auth"
	"github.com/ABHIJNA18/strava-ai-coach/internal/coach"
	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
	"github.com/ABHIJNA18/strava-ai-coach/internal/handlers"
	"github.com/ABHIJNA18/strava-ai-coach/internal/middleware"
	"github.com/ABHIJNA18/strava-ai-coach/internal/strava"
	"github.com/ABHIJNA18/strava-ai-coach/internal/webhook"
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

	//========create connect to databse==============
	db, err := database.NewPostgresConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//========PORTS=================
	//use PORT defined in env file for production, else fall back to localhost:8080 for development
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	grpcAddress := os.Getenv("PYTHON_GRPC_ADDRESS")
	if grpcAddress == "" {
		grpcAddress = "localhost:50051"
	}

	//========SESSIONS=================
	//Create the session manager after creating the database connection

	sessionSecure := strings.EqualFold(
		os.Getenv("APP_ENV"),
		"production",
	)

	//Check if explicitly whether the cookie should be Secure
	//If yes, it overrides the automatic APP_ENV decision.

	if configuredSecure := os.Getenv(
		"SESSION_COOKIE_SECURE",
	); configuredSecure != "" {
		parsedSecure, parseErr := strconv.ParseBool(
			configuredSecure,
		)
		if parseErr == nil {
			sessionSecure = parsedSecure
		}
	}

	sessionCookieName := os.Getenv(
		"SESSION_COOKIE_NAME",
	)

	//if session cookie name isn't set, use default values based on whether the cookie is secure or not
	if sessionCookieName == "" {
		if sessionSecure {
			sessionCookieName = "__Host-session"
		} else {
			sessionCookieName = "session_id"
		}
	}

	sessionManager := auth.NewSessionManager(
		db,
		sessionCookieName,
		sessionSecure,
	)

	// ============Connect to Python Coach Service====================

	coachClient, coachConn, err := coach.NewClient(grpcAddress)
	if err != nil {
		panic(err)
	}
	defer coachConn.Close()

	coachService := coach.NewService(db, coachClient)
	coachHandler := handlers.NewCoachHandler(coachService)

	//========= STRAVA VARIABLES ==========================
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

	//======== STRAVA REDIRECT URI ==========================
	redirectURI := os.Getenv("STRAVA_REDIRECT_URI")

	if redirectURI == "" {
		panic("STRAVA_REDIRECT_URI environment variable is not set")
	}

	// ============Sync Service====================
	syncService := strava.NewSyncService(db, clientID, clientSecret)

	// ============WEBHOOK SERVICE====================

	webhookSigningSecret := os.Getenv("STRAVA_WEBHOOK_SIGNING_SECRET")
	// Signature verification is currently optional because Strava's webhook
	// signature support is not reliably verifiable. The value is retained for
	// future re-enablement without blocking webhook startup today.

	webhookVerifyToken := os.Getenv(
		"STRAVA_WEBHOOK_VERIFY_TOKEN",
	)

	if webhookVerifyToken == "" {
		panic(
			"STRAVA_WEBHOOK_VERIFY_TOKEN environment variable is not set",
		)
	}

	webhookService := webhook.NewWebhookService(
		db,
		syncService,
	)

	webhookHandler := webhook.NewHandler(
		webhookService,
		webhookVerifyToken,
		webhookSigningSecret,
	)
	//=============auth and activity handlers=========================

	activityHandler := handlers.NewActivityHandler(db)
	authHandler := handlers.NewAuthHandler(db, clientID, clientSecret, redirectURI, sessionManager, syncService, sessionSecure)

	//stats handler
	statsHandler := handlers.NewStatsHandler(db)

	//Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Strava AI Coach running")
	})

	//===============Create a protected-handler helper=======================
	//protected() takes the handler you want to protect, wraps it with RequireAuth, and returns the protected handler.

	protected := func(
		handler http.HandlerFunc,
	) http.Handler {
		return middleware.RequireAuth(
			sessionManager,
			http.HandlerFunc(handler),
		)
	}

	//==== AUTH ENDPOINTS  =====

	http.HandleFunc("/login", authHandler.Login)
	//callback endpoint
	http.HandleFunc("/oauth/callback", authHandler.Callback)

	//==== STATS ENDPOINTS =====

	http.Handle("/stats", protected(activityHandler.GetStats))
	http.Handle("/stats/top-sport", protected(statsHandler.GetTopSport))

	//==== ACTIVITY ENDPOINTS =====
	http.Handle("/activities", protected(activityHandler.GetActivities))
	http.Handle("/activities/runs", protected(activityHandler.GetRuns))
	http.Handle("/activities/hikes", protected(activityHandler.GetHikes))
	http.Handle("/activities/weight-training", protected(activityHandler.GetWeightTraining))
	http.Handle("/activities/recent", protected(activityHandler.GetRecentActivities))
	http.Handle("/activities/recent/runs", protected(activityHandler.GetRecentRuns))
	http.Handle("/activities/recent/hikes", protected(activityHandler.GetRecentHikes))
	http.Handle("/activities/recent/weight-training", protected(activityHandler.GetRecentWeightTraining))

	//==== COACH ENDPOINTS =====

	http.Handle("/coach/report", protected(coachHandler.GetReport))
	http.Handle("/coach/coaching", protected(coachHandler.GetCoaching))

	//==== FRONTEND =====

	fileServer := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fileServer)

	//==== FRONTEND DASHBOARD=====

	http.Handle(
		"/dashboard", middleware.RequireAuth(sessionManager, http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", "no-store")
				http.ServeFile(
					w,
					r,
					"./frontend/dashboard.html",
				)
			},
		),
		),
	)

	//==== LOGOUT=====
	http.HandleFunc("/logout", authHandler.Logout)

	//==== WEBHOOKS =====
	http.Handle("/webhooks/strava", webhookHandler)

	//==== START SERVER =====
	/*
		if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("HTTP server stopped:", err)
		}*/

	//replacing the above with explicit http server below
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	fmt.Printf("Server running on port %s...\n", port)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		fmt.Println("HTTP server stopped:", err)
	}

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

	fmt.Printf("Access Token Retrieved Successfully")
})
*/
