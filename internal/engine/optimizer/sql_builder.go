package optimizer

import (
	"fmt"
	"strings"

	"github.com/singh-anurag-7991/data-guard/internal/domain"
)

// BuildFailureQuery constructs a SQL query to find records that FAIL the rules.
// Logic: If Rule is "amount > 0", Failure is "amount <= 0 OR amount IS NULL".
func BuildFailureQuery(tableName string, rules []domain.Rule) (string, []interface{}) {
	if len(rules) == 0 {
		return "", nil
	}

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	for _, rule := range rules {
		ruleConditions := []string{}

		for _, check := range rule.Checks {
			cond, val := invertCheckToSQL(rule.Field, check)
			if cond != "" {
				// Only append arg if val is not nil (some ops like IS NULL don't need args)
				if val != nil {
					ruleConditions = append(ruleConditions, fmt.Sprintf("%s $%d", cond, argCounter))
					args = append(args, val)
					argCounter++
				} else {
					ruleConditions = append(ruleConditions, cond)
				}
			}
		}

		// A rule fails if ANY of its inverted checks match (wait, rules usually imply AND? No, a single check failure fails the rule)
		// Actually, standard rule: all checks must pass. So if ANY check fails, rule fails.
		// So we OR the inverted conditions.
		if len(ruleConditions) > 0 {
			// Wrap in parens: (amount <= 0 OR amount IS NULL)
			clause := fmt.Sprintf("(%s)", strings.Join(ruleConditions, " OR "))
			whereClauses = append(whereClauses, clause)
		}
	}

	if len(whereClauses) == 0 {
		return "", nil
	}

	// We want to return rows that fail ANY rule.
	fullWhere := strings.Join(whereClauses, " OR ")
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", tableName, fullWhere)
	return query, args
}

// invertCheckToSQL returns the INVERTED condition (what makes it fail)
func invertCheckToSQL(field string, check domain.Check) (string, interface{}) {
	switch check.Op {
	case "not_null":
		// Fail if IS NULL
		return fmt.Sprintf("%s IS NULL", field), nil
	case "eq":
		// Fail if !=
		return fmt.Sprintf("%s !=", field), check.Value
	case "neq":
		// Fail if =
		return fmt.Sprintf("%s =", field), check.Value
	case "gt":
		// Fail if <=
		return fmt.Sprintf("%s <=", field), check.Value
	case "lt":
		// Fail if >=
		return fmt.Sprintf("%s >=", field), check.Value
	case "gte":
		// Fail if <
		return fmt.Sprintf("%s <", field), check.Value
	case "lte":
		// Fail if >
		return fmt.Sprintf("%s >", field), check.Value
	default:
		return "", nil
	}
}
