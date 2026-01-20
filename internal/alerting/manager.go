package alerting

import (
	"fmt"
	"time"

	"github.com/singh-anurag-7991/data-guard/internal/domain"
)

type State string

const (
	StateOK   State = "PASS"
	StateFail State = "FAIL"
)

type Notifier interface {
	Send(title, message, color string) error
}

type StateManager interface {
	GetLastState(sourceID string) (State, error)
	UpdateState(sourceID string, state State) error
}

type Manager struct {
	notifier     Notifier
	stateManager StateManager
}

func NewManager(n Notifier, sm StateManager) *Manager {
	return &Manager{
		notifier:     n,
		stateManager: sm,
	}
}

// ProcessResult decides whether to send an alert based on the result and previous state
func (m *Manager) ProcessResult(res domain.ValidationResult) error {
	lastState, err := m.stateManager.GetLastState(res.SourceID)
	if err != nil {
		// If no state found, assume OK (first run)
		lastState = StateOK
	}

	currentState := State(res.Status)

	// State Machine Logic
	if lastState == StateOK && currentState == StateFail {
		// New Failure -> Alert
		msg := fmt.Sprintf("Source '%s' has failed validation.\nRules Failed: %d\nTime: %s",
			res.SourceID, res.RulesFailed, res.Timestamp.Format(time.RFC3339))
		if err := m.notifier.Send("🚨 Data Validation Failed", msg, "#FF0000"); err != nil {
			return err
		}
		return m.stateManager.UpdateState(res.SourceID, StateFail)
	}

	if lastState == StateFail && currentState == StateOK {
		// Recovery -> Alert
		msg := fmt.Sprintf("Source '%s' has recovered and is passing validation.", res.SourceID)
		if err := m.notifier.Send("✅ Data Validation Recovered", msg, "#36a64f"); err != nil {
			return err
		}
		return m.stateManager.UpdateState(res.SourceID, StateOK)
	}

	// OK -> OK or FAIL -> FAIL: Do nothing (Suppress)
	return nil
}
