package optimizer

import (
	"testing"

	"github.com/singh-anurag-7991/data-guard/internal/domain"
)

func TestBuildFailureQuery(t *testing.T) {
	rules := []domain.Rule{
		{
			ID:    "r1",
			Field: "amount",
			Checks: []domain.Check{
				{Op: "not_null"},
				{Op: "gt", Value: 0},
			},
		},
		{
			ID:    "r2",
			Field: "status",
			Checks: []domain.Check{
				{Op: "eq", Value: "active"},
			},
		},
	}

	query, args := BuildFailureQuery("orders", rules)

	// Expected: SELECT * FROM orders WHERE (amount IS NULL OR amount <= $1) OR (status != $2)
	// Note: The order of map iteration in `Plan` wasn't map based, but `rules` is a slice, so order is preserved.

	// Check Args
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d", len(args))
	}
	if args[0] != 0 {
		t.Errorf("expected arg[0] to be 0, got %v", args[0])
	}
	if args[1] != "active" {
		t.Errorf("expected arg[1] to be 'active', got %v", args[1])
	}

	// Check Query Structure (Basic substring check to avoid whitespace brittleness)
	expectedFragments := []string{
		"SELECT * FROM orders WHERE",
		"(amount IS NULL OR amount <= $1)",
		"OR",
		"(status != $2)",
	}

	for _, frag := range expectedFragments {
		if !contains(query, frag) {
			t.Errorf("query missing fragment '%s'. Got: %s", frag, query)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}
