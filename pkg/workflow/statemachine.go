package workflow

// State represents a workflow phase.
type State string

const (
	StatePlan      State = "plan"
	StateImplement State = "implement"
	StateVerify    State = "verify"
	StateFinish    State = "finish"
)

// StateMachine defines valid transitions between workflow states.
type StateMachine struct {
	States      []State
	Transitions map[State][]State
}

// NewStateMachine creates the default 4-phase state machine.
func NewStateMachine() *StateMachine {
	return &StateMachine{
		States: []State{StatePlan, StateImplement, StateVerify, StateFinish},
		Transitions: map[State][]State{
			StatePlan:      {StateImplement},
			StateImplement: {StateVerify},
			StateVerify:    {StateImplement, StateFinish},
			StateFinish:    {},
		},
	}
}

// CanTransition checks if a transition from one state to another is valid.
func (sm *StateMachine) CanTransition(from, to State) bool {
	valid, ok := sm.Transitions[from]
	if !ok {
		return false
	}
	for _, s := range valid {
		if s == to {
			return true
		}
	}
	return false
}

// NextStates returns all valid next states from the given state.
func (sm *StateMachine) NextStates(from State) []State {
	return sm.Transitions[from]
}
