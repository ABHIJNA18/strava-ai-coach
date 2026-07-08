package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
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

	activities, err := database.GetActivitiesByAthleteID(h.DB, 1)
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

	stats, err := database.GetActivityStats(h.DB, 1)

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

	runs, err := database.GetActivitiesByType(h.DB, 1, "Run")
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

	hikes, err := database.GetActivitiesByType(
		h.DB,
		1,
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

	workouts, err := database.GetActivitiesByType(
		h.DB,
		1,
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

	runs, err := database.GetRecentActivities(
		h.DB,
		1,
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

	runs, err := database.GetRecentActivitiesByType(
		h.DB,
		1,
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

	hikes, err := database.GetRecentActivitiesByType(
		h.DB,
		1,
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

	workouts, err := database.GetRecentActivitiesByType(
		h.DB,
		1,
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
