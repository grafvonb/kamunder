package process

import (
	"errors"
	"fmt"
	"strings"
)

// State is the process-instance state filter.
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

var ErrUnknownStateFilter = errors.New("is unknown (valid: all, active, canceled, completed)")
