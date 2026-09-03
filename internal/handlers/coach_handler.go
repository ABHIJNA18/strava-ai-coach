// protobuf still returns the structured coaching JSON as a string. This handler converts that string into a real JSON object for the browser.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ABHIJNA18/strava-ai-coach/internal/middleware"
)

type CoachService interface {
	AnalyzeRecentRuns(ctx context.Context, athleteID int64) (string, error)
	GenerateCoaching(ctx context.Context, athleteID int64, goal string) (string, error)
}

type CoachHandler struct {
	coachService CoachService
}

func NewCoachHandler(coachService CoachService) *CoachHandler {
	return &CoachHandler{
		coachService: coachService,
	}
}

// for 30 day summary
type coachReportResponse struct {
	Summary string `json:"summary"`
}

// for personalised coaching
type coachingRequest struct {
	Goal string `json:"goal"`
}

type coachingResponse struct {
	Coaching json.RawMessage `json:"coaching"`
}

func (h *CoachHandler) GetReport(w http.ResponseWriter, r *http.Request) {

	//get athleteID from the context, set by middleware

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
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("coach report request started:", r.Method, r.URL.Path, time.Now().UnixNano())
	summary, err := h.coachService.AnalyzeRecentRuns(r.Context(), athleteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("coach report request finished:", time.Now().UnixNano())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coachReportResponse{
		Summary: summary,
	})
}

// GetCoaching receives a goal and returns structured personalized coaching.
func (h *CoachHandler) GetCoaching(w http.ResponseWriter, r *http.Request) {

	//get athleteID from the context, set by middleware

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

	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var request coachingRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid JSON body",
			http.StatusBadRequest,
		)
		return
	}

	request.Goal = strings.TrimSpace(
		request.Goal,
	)

	if request.Goal == "" {
		http.Error(
			w,
			"goal cannot be empty",
			http.StatusBadRequest,
		)
		return
	}

	coaching, err := h.coachService.GenerateCoaching(r.Context(), athleteID, request.Goal)

	if err != nil {
		http.Error(
			w,
			"failed to generate coaching",
			http.StatusInternalServerError,
		)
		return
	}

	if !json.Valid([]byte(coaching)) {
		http.Error(
			w,
			"invalid structured coaching response",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(
		coachingResponse{
			Coaching: json.RawMessage(coaching),
		},
	)
	if err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}
