package optimizer

import "github.com/singh-anurag-7991/data-guard/internal/domain"

// ExecutionPlan determines how rules should be executed
type ExecutionPlan struct {
	SQLRules    []domain.Rule // Can be converted to SQL WHERE clauses
	MemoryRules []domain.Rule // Must be checked in Go
}

// Plan separates rules into execution buckets
func Plan(rules []domain.Rule) ExecutionPlan {
	plan := ExecutionPlan{
		SQLRules:    []domain.Rule{},
		MemoryRules: []domain.Rule{},
	}

	for _, rule := range rules {
		if isSQLPushdownSafe(rule) {
			plan.SQLRules = append(plan.SQLRules, rule)
		} else {
			plan.MemoryRules = append(plan.MemoryRules, rule)
		}
	}

	return plan
}

// isSQLPushdownSafe is the decision logic for the optimizer
func isSQLPushdownSafe(rule domain.Rule) bool {
	// 1. If there's a "When" condition, we only support basic SQL operators
	if rule.When != nil {
		if !isOpSafe(rule.When.Op) {
			return false
		}
	}

	// 2. Check all check operators
	for _, check := range rule.Checks {
		if !isOpSafe(check.Op) {
			return false
		}
	}

	return true
}

func isOpSafe(op string) bool {
	switch op {
	case "not_null", "eq", "neq", "gt", "lt", "gte", "lte":
		return true
	// "regex" and "enum" can be supported later but might vary by DB
	default:
		return false
	}
}
