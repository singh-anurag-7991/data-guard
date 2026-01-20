package optimizer

import (
	"testing"

	"github.com/singh-anurag-7991/data-guard/internal/domain"
)

func TestPlan(t *testing.T) {
	rules := []domain.Rule{
		{
			ID:    "sql_safe",
			Field: "amount",
			Checks: []domain.Check{
				{Op: "gt", Value: 10},
			},
		},
		{
			ID:    "memory_only",
			Field: "email",
			Checks: []domain.Check{
				{Op: "unknown_op", Value: "foo"},
			},
		},
		{
			ID:    "sql_safe_conditional",
			Field: "status",
			When: &domain.Condition{
				Field: "type",
				Op:    "eq",
				Value: "new",
			},
			Checks: []domain.Check{
				{Op: "eq", Value: "active"},
			},
		},
	}

	plan := Plan(rules)

	if len(plan.SQLRules) != 2 {
		t.Errorf("expected 2 SQL rules, got %d", len(plan.SQLRules))
	}
	if len(plan.MemoryRules) != 1 {
		t.Errorf("expected 1 Memory rule, got %d", len(plan.MemoryRules))
	}
	if plan.MemoryRules[0].ID != "memory_only" {
		t.Errorf("expected 'memory_only' to be in memory rules")
	}
}
