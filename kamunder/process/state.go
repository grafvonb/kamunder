package process

import "strings"

type State string

const (
	StateAll        State = "ALL"
	StateActive     State = "ACTIVE"
	StateCompleted  State = "COMPLETED"
	StateCanceled   State = "CANCELED"
	StateTerminated State = "TERMINATED"
)

func (s State) String() string { return string(s) }

func (s State) EqualsIgnoreCase(other State) bool {
	return strings.EqualFold(s.String(), other.String())
}

func (s State) In(states ...State) bool {
	for _, st := range states {
		if s.EqualsIgnoreCase(st) {
			return true
		}
	}
	return false
}

func ParseState(in string) (State, bool) {
	switch strings.ToLower(in) {
	case "all":
		return StateAll, true
	case "active":
		return StateActive, true
	case "completed":
		return StateCompleted, true
	case "canceled", "cancelled":
		return StateCanceled, true
	case "terminated":
		return StateTerminated, true
	default:
		return "", false
	}
}

func (s State) IsTerminal() bool {
	return s.In(StateCompleted, StateCanceled, StateTerminated)
}
