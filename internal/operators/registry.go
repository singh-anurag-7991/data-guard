package operators

import (
	"fmt"
	"regexp"

	"github.com/singh-anurag-7991/data-guard/internal/domain"
)

// OperatorFunc defines the signature for a validation check
type OperatorFunc func(value interface{}, check domain.Check) (bool, string)

// Registry holds all available operators
var Registry = map[string]OperatorFunc{
	"not_null": notNull,
	"eq":       equal,
	"neq":      notEqual,
	"gt":       greaterThan,
	"lt":       lessThan,
	"regex":    regexMatch,
	"enum":     enumMatch,
}

// Helper to get operator
func Get(op string) (OperatorFunc, bool) {
	fn, exists := Registry[op]
	return fn, exists
}

// --- Implementations ---

func notNull(value interface{}, check domain.Check) (bool, string) {
	if value == nil {
		return false, "value is null"
	}
	return true, ""
}

func equal(value interface{}, check domain.Check) (bool, string) {
	if value == check.Value {
		return true, ""
	}
	return false, fmt.Sprintf("expected %v, got %v", check.Value, value)
}

func notEqual(value interface{}, check domain.Check) (bool, string) {
	if value != check.Value {
		return true, ""
	}
	return false, fmt.Sprintf("expected not %v, got %v", check.Value, value)
}

func greaterThan(value interface{}, check domain.Check) (bool, string) {
	v, ok := ToFloat(value)
	if !ok {
		return false, "value is not a number"
	}
	t, ok := ToFloat(check.Value)
	if !ok {
		return false, "threshold is not a number"
	}
	if v > t {
		return true, ""
	}
	return false, fmt.Sprintf("value %v is not greater than %v", v, t)
}

func lessThan(value interface{}, check domain.Check) (bool, string) {
	v, ok := ToFloat(value)
	if !ok {
		return false, "value is not a number"
	}
	t, ok := ToFloat(check.Value)
	if !ok {
		return false, "threshold is not a number"
	}
	if v < t {
		return true, ""
	}
	return false, fmt.Sprintf("value %v is not less than %v", v, t)
}

func regexMatch(value interface{}, check domain.Check) (bool, string) {
	vStr, ok := value.(string)
	if !ok {
		return false, "value is not a string"
	}
	pattern, ok := check.Value.(string)
	if !ok {
		return false, "pattern is not a string"
	}
	matched, err := regexp.MatchString(pattern, vStr)
	if err != nil {
		return false, fmt.Sprintf("invalid regex: %v", err)
	}
	if matched {
		return true, ""
	}
	return false, fmt.Sprintf("value %s does not match pattern %s", vStr, pattern)
}

func enumMatch(value interface{}, check domain.Check) (bool, string) {
	allowedList, ok := check.Value.([]interface{})
	if !ok {
		// Try string slice
		if strList, ok := check.Value.([]string); ok {
			for _, item := range strList {
				if item == value {
					return true, ""
				}
			}
			return false, fmt.Sprintf("value %v not in enum list", check.Value)
		}
		return false, "enum list must be an array"
	}

	for _, item := range allowedList {
		if item == value {
			return true, ""
		}
	}
	return false, fmt.Sprintf("value %v not in enum list", check.Value)
}

// Utility to convert numbers to float64 safely
func ToFloat(i interface{}) (float64, bool) {
	switch v := i.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}
