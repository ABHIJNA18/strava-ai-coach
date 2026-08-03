package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CoachService interface {
	AnalyzeRecentRuns(ctx context.Context, athleteID int64) (string, error)
}

type CoachHandler struct {
	coachService CoachService
}

func NewCoachHandler(coachService CoachService) *CoachHandler {
	return &CoachHandler{
		coachService: coachService,
	}
}

type coachReportResponse struct {
	Summary string `json:"summary"`
}

func (h *CoachHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("coach report request started:", r.Method, r.URL.Path, time.Now().UnixNano())
	summary, err := h.coachService.AnalyzeRecentRuns(r.Context(), 1)
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
