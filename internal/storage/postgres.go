package storage

import (
	"context"
	"fmt"

	"github.com/singh-anurag-7991/data-guard/internal/alerting"
	"github.com/singh-anurag-7991/data-guard/internal/domain"
	"github.com/singh-anurag-7991/data-guard/internal/ingest/postgres"
)

type Repository struct {
	client *postgres.Client
}

func NewRepository(client *postgres.Client) *Repository {
	return &Repository{client: client}
}

// SaveResult persists the validation run and any errors
func (r *Repository) SaveResult(ctx context.Context, res domain.ValidationResult) error {
	tx, err := r.client.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Insert Run
	var runID int
	err = tx.QueryRow(ctx, `
		INSERT INTO validation_runs (source_id, status, records_checked, rules_failed, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		res.SourceID, res.Status, res.RecordsChecked, res.RulesFailed, res.Timestamp,
	).Scan(&runID)
	if err != nil {
		return fmt.Errorf("failed to insert validation run: %w", err)
	}

	// 2. Insert Errors (if any)
	if len(res.Errors) > 0 {
		for _, e := range res.Errors {
			// Convert value to string safety
			valStr := fmt.Sprintf("%v", e.Value)
			_, err := tx.Exec(ctx, `
				INSERT INTO validation_errors (run_id, rule_id, field, fail_value, reason)
				VALUES ($1, $2, $3, $4, $5)`,
				runID, e.RuleID, e.Field, valStr, e.Reason,
			)
			if err != nil {
				return fmt.Errorf("failed to insert error: %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}

// GetLastState returning status for alerting manager
func (r *Repository) GetLastState(ctx context.Context, sourceID string) (alerting.State, error) {
	var status string
	err := r.client.Pool().QueryRow(ctx, `
		SELECT last_status FROM alert_states WHERE source_id = $1`, sourceID,
	).Scan(&status)

	if err != nil {
		// If no row found (pgx does return err for no rows), return error to indicate "unknown/first run"
		return "", err
	}
	return alerting.State(status), nil
}

// UpdateState updates the alert state
func (r *Repository) UpdateState(ctx context.Context, sourceID string, state alerting.State) error {
	_, err := r.client.Pool().Exec(ctx, `
		INSERT INTO alert_states (source_id, last_status, last_alerted_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (source_id) DO UPDATE 
		SET last_status = EXCLUDED.last_status, last_alerted_at = NOW()`,
		sourceID, string(state),
	)
	return err
}

// GetRecentRuns fetches the latest validation runs, optionally filtered by sourceID
func (r *Repository) GetRecentRuns(ctx context.Context, sourceID string, limit int) ([]domain.ValidationResult, error) {
	query := `
		SELECT id, source_id, status, records_checked, rules_failed, created_at
		FROM validation_runs
		WHERE ($1 = '' OR source_id = $1)
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.client.Pool().Query(ctx, query, sourceID, limit)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []domain.ValidationResult
	for rows.Next() {
		var res domain.ValidationResult
		var id int // Not currently part of domain model, but good to know
		err := rows.Scan(&id, &res.SourceID, &res.Status, &res.RecordsChecked, &res.RulesFailed, &res.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		results = append(results, res)
	}
	return results, nil
}
