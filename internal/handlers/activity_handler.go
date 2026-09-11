package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
	"github.com/ABHIJNA18/strava-ai-coach/internal/middleware"
)

// define the struct
type ActivityHandler struct {
	DB *sql.DB
}

// fucntion to create a handler which has GetRuns, GetActivities etc methods attached to it
func NewActivityHandler(db *sql.DB) *ActivityHandler {
	return &ActivityHandler{
		DB: db,
	}
}

//attach method GetActivities to this  struct

func (h *ActivityHandler) GetActivities(w http.ResponseWriter, r *http.Request) {

	//retireve the athleteID from the context, which was set by the middleware
	athleteID, ok := middleware.AthleteIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			"unauhthorized",
			http.StatusUnauthorized,
		)
		return
	}

	activities, err := database.GetActivitiesByAthleteID(h.DB, athleteID)
	if err != nil {
		fmt.Println("Failed to load activities:", err)
		http.Error(w, "failed to load activities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(activities)

	if err != nil {
		fmt.Println("Failed to encode activities response:", err)
		http.Error(w, "failed to encode activities response", http.StatusInternalServerError)
		return
	}

}
func (h *ActivityHandler) GetStats(

	w http.ResponseWriter,
	r *http.Request,

) {

	//retireve the athleteID from the context, which was set by the middleware
	athleteID, ok := middleware.AthleteIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			"unauhthorized",
			http.StatusUnauthorized,
		)
		return
	}

	stats, err := database.GetActivityStats(h.DB, athleteID)

	if err != nil {
		fmt.Println("Failed to load activity stats:", err)
		http.Error(w, "failed to load activity statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		fmt.Println("Failed to encode activity stats response:", err)
		http.Error(w, "failed to encode activity statistics", http.StatusInternalServerError)
		return
	}
}

func (h *ActivityHandler) GetRuns(w http.ResponseWriter, r *http.Request) {

	//retireve the athleteID from the context, which was set by the middleware
	athleteID, ok := middleware.AthleteIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			"unauhthorized",
			http.StatusUnauthorized,
		)
		return
	}

	runs, err := database.GetActivitiesByType(h.DB, athleteID, "Run")
	if err != nil {
		fmt.Println("Failed to load runs:", err)
		http.Error(w, "failed to load runs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(runs)
	if err != nil {
		fmt.Println("Failed to encode runs response:", err)
		http.Error(w, "failed to encode runs response", http.StatusInternalServerError)
		return
	}
}

func (h *ActivityHandler) GetHikes(
	w http.ResponseWriter,
	r *http.Request,
) {

	//retireve the athleteID from the context, which was set by the middleware
	athleteID, ok := middleware.AthleteIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			"unauhthorized",
			http.StatusUnauthorized,
		)
		return
	}

	hikes, err := database.GetActivitiesByType(
		h.DB,
		athleteID,
		"Hike",
	)

	if err != nil {
		fmt.Println("Failed to load hikes:", err)
		http.Error(w, "failed to load hikes", http.StatusInternalServerError)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(
		hikes,
	)

	if err != nil {
		fmt.Println("Failed to encode hikes response:", err)
		http.Error(w, "failed to encode hikes response", http.StatusInternalServerError)
		return
	}
}

func (h *ActivityHandler) GetWeightTraining(
	w http.ResponseWriter,
	r *http.Request,
) {

	//retireve the athleteID from the context, which was set by the middleware
	athleteID, ok := middleware.AthleteIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			"unauhthorized",
			http.StatusUnauthorized,
		)
		return
	}

	workouts, err := database.GetActivitiesByType(
		h.DB,
		athleteID,
		"WeightTraining",
	)

	if err != nil {
		fmt.Println("Failed to load weight-training activities:", err)
		http.Error(w, "failed to load weight-training activities", http.StatusInternalServerError)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(
		workouts,
	)

	if err != nil {
		fmt.Println("Failed to encode weight-training response:", err)
		http.Error(w, "failed to encode weight-training response", http.StatusInternalServerError)
		return
	}
}

// ===RECENT ACTIVITIES ====

func (h *ActivityHandler) GetRecentActivities(
	w http.ResponseWriter,
	r *http.Request,
) {

	//retireve the athleteID from the context, which was set by the middleware
	athleteID, ok := middleware.AthleteIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			"unauhthorized",
			http.StatusUnauthorized,
		)
		return
	}

	runs, err := database.GetRecentActivities(
		h.DB,
		athleteID,
		10,
	)

	if err != nil {
		fmt.Println("Failed to load recent activities:", err)
		http.Error(w, "failed to load recent activities", http.StatusInternalServerError)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(runs)
	if err != nil {
		fmt.Println("Failed to encode recent activities response:", err)
		http.Error(w, "failed to encode recent activities response", http.StatusInternalServerError)
		return
	}
}

func (h *ActivityHandler) GetRecentRuns(
	w http.ResponseWriter,
	r *http.Request,
) {

	//retireve the athleteID from the context, which was set by the middleware
	athleteID, ok := middleware.AthleteIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			"unauhthorized",
			http.StatusUnauthorized,
		)
		return
	}

	runs, err := database.GetRecentActivitiesByType(
		h.DB,
		athleteID,
		"Run",
		10,
	)

	if err != nil {
		fmt.Println("Failed to load recent runs:", err)
		http.Error(w, "failed to load recent runs", http.StatusInternalServerError)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(runs)
	if err != nil {
		fmt.Println("Failed to encode recent runs response:", err)
		http.Error(w, "failed to encode recent runs response", http.StatusInternalServerError)
		return
	}
}

func (h *ActivityHandler) GetRecentHikes(
	w http.ResponseWriter,
	r *http.Request,
) {

	//retireve the athleteID from the context, which was set by the middleware
	athleteID, ok := middleware.AthleteIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			"unauhthorized",
			http.StatusUnauthorized,
		)
		return
	}

	hikes, err := database.GetRecentActivitiesByType(
		h.DB,
		athleteID,
		"Hike",
		10,
	)

	if err != nil {
		fmt.Println("Failed to load recent hikes:", err)
		http.Error(w, "failed to load recent hikes", http.StatusInternalServerError)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(hikes)
	if err != nil {
		fmt.Println("Failed to encode recent hikes response:", err)
		http.Error(w, "failed to encode recent hikes response", http.StatusInternalServerError)
		return
	}
}

func (h *ActivityHandler) GetRecentWeightTraining(
	w http.ResponseWriter,
	r *http.Request,
) {

	//retireve the athleteID from the context, which was set by the middleware
	athleteID, ok := middleware.AthleteIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			"unauhthorized",
			http.StatusUnauthorized,
		)
		return
	}

	workouts, err := database.GetRecentActivitiesByType(
		h.DB,
		athleteID,
		"WeightTraining",
		10,
	)

	if err != nil {
		fmt.Println("Failed to load recent weight-training activities:", err)
		http.Error(w, "failed to load recent weight-training activities", http.StatusInternalServerError)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(workouts)
	if err != nil {
		fmt.Println("Failed to encode recent weight-training response:", err)
		http.Error(w, "failed to encode recent weight-training response", http.StatusInternalServerError)
		return
	}
}
