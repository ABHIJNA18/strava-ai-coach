package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/database"
	"github.com/ABHIJNA18/strava-ai-coach/internal/middleware"
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

	//get the athleteID from the context, which was set by the middleware
	athleteID, ok := middleware.AthleteIDFromContext(
		r.Context(),
	)
	if !ok {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

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
		athleteID,
		since,
	)

	if err != nil {
		fmt.Println("Failed to load top sport:", err)
		http.Error(w, "failed to load top sport", http.StatusInternalServerError)
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
		fmt.Println("Failed to encode top sport response:", err)
		http.Error(w, "failed to encode top sport response", http.StatusInternalServerError)
		return
	}

}
