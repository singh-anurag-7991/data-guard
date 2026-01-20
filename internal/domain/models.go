package domain

import (
	"time"
)

// Record represents a single data record (row or JSON object)
type Record map[string]interface{}

// Schema defines the expected structure of the data
type Schema map[string]string // specific field -> expected type (e.g., "string", "number", "timestamp")

// Condition defines when a rule should be applied
type Condition struct {
	Field string      `json:"field"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

// Check defines the actual validation logic
type Check struct {
	Op    string      `json:"op"`
	Value interface{} `json:"value,omitempty"`
}

// Rule defines a validation rule
type Rule struct {
	ID       string      `json:"id"`
	Field    string      `json:"field"`
	When     *Condition  `json:"when,omitempty"` // Pointer to allow null (always apply)
	Checks   []Check     `json:"checks"`
	Severity string      `json:"severity"` // "error", "warning", "info"
}

// ValidationResult represents the outcome of a validation run
type ValidationResult struct {
	SourceID       string        `json:"source_id"`
	Status         string        `json:"status"` // "PASS", "FAIL"
	RecordsChecked int           `json:"records_checked"`
	RulesFailed    int           `json:"rules_failed"`
	Errors         []ErrorDetail `json:"errors,omitempty"`
	Timestamp      time.Time     `json:"timestamp"`
}

// ErrorDetail captures specific validation failures
type ErrorDetail struct {
	RuleID   string      `json:"rule_id"`
	Field    string      `json:"field"`
	Value    interface{} `json:"value"`
	Reason   string      `json:"reason"`
	RecordID string      `json:"record_id,omitempty"` // specific identifier if available
}
