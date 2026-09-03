package handlers

import (
	"database/sql"
	"encoding/json"
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

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(activities)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(runs)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(runs)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(runs)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(hikes)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(workouts)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
}
