package engine

import (
	"fmt"
	"time"

	"github.com/singh-anurag-7991/data-guard/internal/domain"
	"github.com/singh-anurag-7991/data-guard/internal/operators"
)

// Executor is responsible for running validations
type Executor struct{}

// NewExecutor creates a new validation executor
func NewExecutor() *Executor {
	return &Executor{}
}

// Validate executes the rules against the provided records
func (e *Executor) Validate(sourceID string, schema domain.Schema, rules []domain.Rule, records []domain.Record) domain.ValidationResult {
	result := domain.ValidationResult{
		SourceID:       sourceID,
		Status:         "PASS",
		RecordsChecked: len(records),
		Errors:         []domain.ErrorDetail{},
		Timestamp:      time.Now(),
	}

	for _, record := range records {
		// 1. Schema Validation (First Gate)
		if err := e.validateSchema(record, schema); err != nil {
			result.Status = "FAIL"
			result.RulesFailed++ // Counting schema failure as a rule failure
			result.Errors = append(result.Errors, *err)
			continue // Skip processing rules if schema fails
		}

		// 2. Rule Execution
		for _, rule := range rules {
			// Check 'When' condition
			if !e.evaluateCondition(record, rule.When) {
				continue // Skip rule if condition not met
			}

			// Execute Rule Checks
			for _, check := range rule.Checks {
				val, exists := record[rule.Field]
				// If field is missing and op is not 'not_null', it might be valid or invalid depending on rule.
				// For simplicity, if field missing and we check it, we treat as nil.
				if !exists {
					val = nil
				}

				opFunc, found := operators.Get(check.Op)
				if !found {
					result.Status = "FAIL"
					result.RulesFailed++
					result.Errors = append(result.Errors, domain.ErrorDetail{
						RuleID: rule.ID,
						Field:  rule.Field,
						Reason: fmt.Sprintf("unknown operator: %s", check.Op),
					})
					continue
				}

				pass, reason := opFunc(val, check)
				if !pass {
					result.Status = "FAIL"
					result.RulesFailed++
					result.Errors = append(result.Errors, domain.ErrorDetail{
						RuleID: rule.ID,
						Field:  rule.Field,
						Value:  val,
						Reason: reason,
					})
				}
			}
		}
	}

	return result
}

// validateSchema checks if fields match expected types (Basic implementation)
func (e *Executor) validateSchema(record domain.Record, schema domain.Schema) *domain.ErrorDetail {
	for field, expectedType := range schema {
		val, exists := record[field]
		if !exists {
			return &domain.ErrorDetail{
				Field:  field,
				Reason: "field missing",
			}
		}

		// Very basic type check
		switch expectedType {
		case "string":
			if _, ok := val.(string); !ok {
				return &domain.ErrorDetail{Field: field, Reason: "expected string"}
			}
		case "number":
			if _, ok := operators.ToFloat(val); !ok {
				return &domain.ErrorDetail{Field: field, Reason: "expected number"}
			}
		case "boolean":
			if _, ok := val.(bool); !ok {
				return &domain.ErrorDetail{Field: field, Reason: "expected boolean"}
			}
		}
	}
	return nil
}

// evaluateCondition checks if the "When" condition is met
func (e *Executor) evaluateCondition(record domain.Record, cond *domain.Condition) bool {
	if cond == nil {
		return true // Always run if no condition
	}

	val, exists := record[cond.Field]
	if !exists {
		val = nil
	}

	check := domain.Check{Op: cond.Op, Value: cond.Value}
	opFunc, found := operators.Get(cond.Op)
	if !found {
		return false // Fail safe
	}

	pass, _ := opFunc(val, check)
	return pass
}
