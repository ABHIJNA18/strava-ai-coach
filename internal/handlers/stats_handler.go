package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
)

type StatsHandler struct {
	DB *sql.DB
}

func NewStatsHandler(db *sql.DB) *StatsHandler {
	return &StatsHandler{
		DB: db,
	}
}

type TopSportResponse struct {
	Sports []database.TopSport `json:"sports"`
}

func (h *StatsHandler) GetTopSport(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	since := time.Now().AddDate(0, 0, -30)

	topSports, err := database.GetTopSportSince(
		h.DB,
		1,
		since,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	//Filling the struct with sport and count value
	//this converts it to JSON and writes it to the HTTP response.
	err = json.NewEncoder(w).Encode(
		TopSportResponse{
			Sports: topSports,
		},
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
