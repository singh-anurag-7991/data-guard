package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/singh-anurag-7991/data-guard/internal/storage"
)

type DashboardHandler struct {
	repo *storage.Repository
}

func NewDashboardHandler(repo *storage.Repository) *DashboardHandler {
	return &DashboardHandler{repo: repo}
}

// ListRuns returns specific validation history
func (h *DashboardHandler) ListRuns(w http.ResponseWriter, r *http.Request) {
	sourceID := r.URL.Query().Get("source_id")
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	runs, err := h.repo.GetRecentRuns(r.Context(), sourceID, limit)
	if err != nil {
		http.Error(w, "Failed to fetch runs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow generic CORS for local dev
	json.NewEncoder(w).Encode(runs)
}
