package storage

import (
	"context"
	"sort"
	"sync"

	"github.com/singh-anurag-7991/data-guard/internal/alerting"
	"github.com/singh-anurag-7991/data-guard/internal/domain"
)

type MemoryStore struct {
	mu          sync.RWMutex
	runs        []domain.ValidationResult
	alertStates map[string]alerting.State
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		runs:        make([]domain.ValidationResult, 0),
		alertStates: make(map[string]alerting.State),
	}
}

func (m *MemoryStore) SaveResult(ctx context.Context, res domain.ValidationResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Prepend or Append? Append is easier, we sort later or reverse iterate.
	m.runs = append(m.runs, res)
	return nil
}

func (m *MemoryStore) GetLastState(ctx context.Context, sourceID string) (alerting.State, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	state, ok := m.alertStates[sourceID]
	if !ok {
		return "", nil // Or specific error? Postgres returns error on no rows. Manager treats error as "StateOK".
	}
	return state, nil
}

func (m *MemoryStore) UpdateState(ctx context.Context, sourceID string, state alerting.State) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alertStates[sourceID] = state
	return nil
}

func (m *MemoryStore) GetRecentRuns(ctx context.Context, sourceID string, limit int) ([]domain.ValidationResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Filter
	var filtered []domain.ValidationResult
	for _, r := range m.runs {
		if sourceID == "" || r.SourceID == sourceID {
			filtered = append(filtered, r)
		}
	}

	// Sort by timestamp descending (newest first)
	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.After(filtered[j].Timestamp)
	})

	// Limit
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}

	return filtered, nil
}
