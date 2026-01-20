package storage

import (
	"context"

	"github.com/singh-anurag-7991/data-guard/internal/alerting"
	"github.com/singh-anurag-7991/data-guard/internal/domain"
)

type Provider interface {
	SaveResult(ctx context.Context, res domain.ValidationResult) error
	GetLastState(ctx context.Context, sourceID string) (alerting.State, error)
	UpdateState(ctx context.Context, sourceID string, state alerting.State) error
	GetRecentRuns(ctx context.Context, sourceID string, limit int) ([]domain.ValidationResult, error)
}
