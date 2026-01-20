package engine

import (
	"testing"

	"github.com/singh-anurag-7991/data-guard/internal/domain"
)

func TestExecutor_Validate(t *testing.T) {
	e := NewExecutor()

	schema := domain.Schema{
		"amount": "number",
		"status": "string",
	}

	rules := []domain.Rule{
		{
			ID:    "amount_positive",
			Field: "amount",
			Checks: []domain.Check{
				{Op: "not_null"},
				{Op: "gt", Value: 0},
			},
		},
		{
			ID:    "status_enum",
			Field: "status",
			Checks: []domain.Check{
				{Op: "enum", Value: []string{"active", "pending", "failed"}},
			},
		},
	}

	tests := []struct {
		name        string
		record      domain.Record
		expectPass  bool
		failedRules int
	}{
		{
			name: "valid_record",
			record: domain.Record{
				"amount": 100,
				"status": "active",
			},
			expectPass:  true,
			failedRules: 0,
		},
		{
			name: "invalid_amount",
			record: domain.Record{
				"amount": -10,
				"status": "active",
			},
			expectPass:  false,
			failedRules: 1, // gt failed
		},
		{
			name: "missing_field_schema_error",
			record: domain.Record{
				"status": "active",
			},
			expectPass:  false,
			failedRules: 1, // schema error
		},
		{
			name: "invalid_type_schema_error",
			record: domain.Record{
				"amount": "not_a_number",
				"status": "active",
			},
			expectPass:  false,
			failedRules: 1, // type error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.Validate("test_source", schema, rules, []domain.Record{tt.record})

			if tt.expectPass && result.Status != "PASS" {
				t.Errorf("expected PASS, got %s (errors: %v)", result.Status, result.Errors)
			}
			if !tt.expectPass && result.Status != "FAIL" {
				t.Errorf("expected FAIL, got PASS")
			}
			if result.RulesFailed != tt.failedRules {
				t.Errorf("expected %d failures, got %d", tt.failedRules, result.RulesFailed)
			}
		})
	}
}

func TestExecutor_ConditionalRule(t *testing.T) {
	e := NewExecutor()
	schema := domain.Schema{"amount": "number", "type": "string"}

	// Rule: If type == "credit", amount must be > 0. If "debit", amount must be < 0
	rules := []domain.Rule{
		{
			ID:    "credit_positive",
			Field: "amount",
			When:  &domain.Condition{Field: "type", Op: "eq", Value: "credit"},
			Checks: []domain.Check{
				{Op: "gt", Value: 0},
			},
		},
	}

	// Case 1: type=credit, amount=100 (Should apply, pass)
	res1 := e.Validate("source", schema, rules, []domain.Record{{"type": "credit", "amount": 100}})
	if res1.Status != "PASS" {
		t.Errorf("Case 1 failed")
	}

	// Case 2: type=credit, amount=-10 (Should apply, fail)
	res2 := e.Validate("source", schema, rules, []domain.Record{{"type": "credit", "amount": -10}})
	if res2.Status != "FAIL" {
		t.Errorf("Case 2 passed but should fail")
	}

	// Case 3: type=debit, amount=-10 (Should NOT apply, pass)
	res3 := e.Validate("source", schema, rules, []domain.Record{{"type": "debit", "amount": -10}})
	if res3.Status != "PASS" {
		t.Errorf("Case 3 failed but condition should have skipped validation")
	}
}
