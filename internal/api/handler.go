package api

import (
	"encoding/json"
	"net/http"

	"github.com/singh-anurag-7991/data-guard/internal/domain"
	"github.com/singh-anurag-7991/data-guard/internal/engine"
)

type IngestRequest struct {
	SourceID string          `json:"source_id"`
	Schema   domain.Schema   `json:"schema"`
	Rules    []domain.Rule   `json:"rules"`
	Data     []domain.Record `json:"data"`
}

type Handler struct {
	executor *engine.Executor
}

func NewHandler(executor *engine.Executor) *Handler {
	return &Handler{
		executor: executor,
	}
}

func (h *Handler) Ingest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SourceID == "" {
		http.Error(w, "source_id is required", http.StatusBadRequest)
		return
	}

	result := h.executor.Validate(req.SourceID, req.Schema, req.Rules, req.Data)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
