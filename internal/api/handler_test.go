package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/singh-anurag-7991/data-guard/internal/domain"
	"github.com/singh-anurag-7991/data-guard/internal/engine"
)

func TestHandler_Ingest(t *testing.T) {
	exec := engine.NewExecutor()
	handler := NewHandler(exec, nil)

	reqBody := IngestRequest{
		SourceID: "test_source",
		Schema: domain.Schema{
			"amount": "number",
		},
		Rules: []domain.Rule{
			{
				ID:    "positive_amount",
				Field: "amount",
				Checks: []domain.Check{
					{Op: "gt", Value: 0},
				},
			},
		},
		Data: []domain.Record{
			{"amount": 100},
			{"amount": -50},
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/ingest/api", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Ingest(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result domain.ValidationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Status != "FAIL" {
		t.Errorf("expected FAIL, got %s", result.Status)
	}
	if result.RulesFailed != 1 {
		t.Errorf("expected 1 rule failure, got %d", result.RulesFailed)
	}
}
