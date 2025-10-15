package domain

import (
	"fmt"
	"strings"
)

type ProcessInstance struct {
	BpmnProcessId             string
	EndDate                   string
	Incident                  bool
	Key                       int64
	ParentFlowNodeInstanceKey int64
	ParentKey                 int64
	ProcessDefinitionKey      int64
	ProcessVersion            int32
	ProcessVersionTag         string
	StartDate                 string
	State                     State
	TenantId                  string
}

type ProcessInstanceSearchFilterOpts struct {
	Key               int64
	BpmnProcessId     string
	ProcessVersion    int32
	ProcessVersionTag string
	State             State
	ParentKey         int64
}

type State string

func (s State) EqualsIgnoreCase(other State) bool {
	return strings.EqualFold(string(s), string(other))
}

const (
	StateAll       State = "all"
	StateActive    State = "active"
	StateCompleted State = "completed"
	StateCanceled  State = "canceled"
)

func (s State) String() string { return string(s) }

// ParseState parses a string (case-insensitive) into a State.
func ParseState(in string) (State, error) {
	switch strings.ToLower(in) {
	case "all":
		return StateAll, nil
	case "active":
		return StateActive, nil
	case "canceled":
		return StateCanceled, nil
	case "completed":
		return StateCompleted, nil
	default:
		return "", fmt.Errorf("%q %w", in, ErrUnknownStateFilter)
	}
}

type CancelResponse struct {
	StatusCode int
	Status     string
}

type ChangeStatus struct {
	Deleted int64
	Message string
}

func (c ChangeStatus) String() string {
	return fmt.Sprintf("deleted: %d, message: %s", c.Deleted, c.Message)
}
