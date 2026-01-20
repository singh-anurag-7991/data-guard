package alerting

import (
	"testing"
	"time"

	"github.com/singh-anurag-7991/data-guard/internal/domain"
)

// Mock objects
type mockNotifier struct {
	sentCount int
}

func (m *mockNotifier) Send(title, message, color string) error {
	m.sentCount++
	return nil
}

type mockStateManager struct {
	state State
}

func (m *mockStateManager) GetLastState(sourceID string) (State, error) {
	return m.state, nil
}

func (m *mockStateManager) UpdateState(sourceID string, state State) error {
	m.state = state
	return nil
}

func TestManager_ProcessResult(t *testing.T) {
	mockNotif := &mockNotifier{}
	mockState := &mockStateManager{state: StateOK}
	manager := NewManager(mockNotif, mockState)

	// 1. Pass -> Fail (Should Alert)
	failRes := domain.ValidationResult{SourceID: "src", Status: "FAIL", RulesFailed: 1, Timestamp: time.Now()}
	_ = manager.ProcessResult(failRes)

	if mockNotif.sentCount != 1 {
		t.Errorf("expected 1 alert for new failure, got %d", mockNotif.sentCount)
	}
	if mockState.state != StateFail {
		t.Errorf("state should update to FAIL")
	}

	// 2. Fail -> Fail (Should Suppress)
	_ = manager.ProcessResult(failRes)
	if mockNotif.sentCount != 1 {
		t.Errorf("expected alert count to stay 1 (suppressed), got %d", mockNotif.sentCount)
	}

	// 3. Fail -> Pass (Should Alert Recovery)
	passRes := domain.ValidationResult{SourceID: "src", Status: "PASS", Timestamp: time.Now()}
	_ = manager.ProcessResult(passRes)

	if mockNotif.sentCount != 2 {
		t.Errorf("expected 2 alerts (1 fail + 1 recovery), got %d", mockNotif.sentCount)
	}
	if mockState.state != StateOK {
		t.Errorf("state should update to OK")
	}
}
