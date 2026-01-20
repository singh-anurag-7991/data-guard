package operators

import (
	"testing"

	"github.com/singh-anurag-7991/data-guard/internal/domain"
)

func TestOperators(t *testing.T) {
	tests := []struct {
		name     string
		op       string
		val      interface{}
		checkVal interface{}
		want     bool
	}{
		// not_null
		{"not_null_valid", "not_null", "foo", nil, true},
		{"not_null_invalid", "not_null", nil, nil, false},

		// eq
		{"eq_string_match", "eq", "hello", "hello", true},
		{"eq_string_mismatch", "eq", "hello", "world", false},
		{"eq_int_match", "eq", 10, 10, true},

		// gt
		{"gt_int_valid", "gt", 15, 10, true},
		{"gt_int_invalid", "gt", 5, 10, false},
		{"gt_float_valid", "gt", 10.5, 10.0, true},

		// lt
		{"lt_int_valid", "lt", 5, 10, true},
		{"lt_int_invalid", "lt", 15, 10, false},

		// regex
		{"regex_match_email", "regex", "test@example.com", `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, true},
		{"regex_fail_email", "regex", "invalid-email", `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, false},

		// enum
		{"enum_valid", "enum", "active", []string{"active", "inactive"}, true},
		{"enum_invalid", "enum", "deleted", []string{"active", "inactive"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opFunc, exists := Get(tt.op)
			if !exists {
				t.Fatalf("operator %s not found", tt.op)
			}
			check := domain.Check{Value: tt.checkVal}
			got, _ := opFunc(tt.val, check)
			if got != tt.want {
				t.Errorf("op %s(%v, %v) = %v, want %v", tt.op, tt.val, tt.checkVal, got, tt.want)
			}
		})
	}
}
