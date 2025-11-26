package state

import (
	"fmt"
	"slices"
)

type stateType[st ~string | ~int] interface {
	~string | ~int
	Transitions() map[st][]st
}

type State[T stateType[T]] struct {
	current T
}

func NewState[T stateType[T]](s T) (State[T], error) {
	state := State[T]{}

	if _, valid := s.Transitions()[s]; !valid {
		return state, fmt.Errorf(
			"could not create state: invalid initial state %q", s,
		)
	}

	state.current = s

	return state, nil
}

func (state *State[T]) Transition(to T) error {
	transitions, ok := to.Transitions()[state.current]

	if !ok {
		return fmt.Errorf(
			"cannot transition state: invalid current state %q", state.current,
		)
	}

	if !slices.Contains(transitions, to) {
		return fmt.Errorf(
			"cannot transition state from %q to %q, invalid transition",
			state.current, to,
		)
	}

	state.current = to

	return nil
}

func (state *State[T]) Get() T {
	return state.current
}
