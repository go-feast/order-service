package order

import (
	"slices"
)

type StateOperator struct {
	o *Order
}

func NewStateOperator(o *Order) *StateOperator {
	return &StateOperator{o: o}
}

// CancelOrder set orders`s state to [Canceled].
// Returning values: order, canceled(bool), error
func (s *StateOperator) CancelOrder() (*Order, bool, error) {
	if s.o.IsCanceled() {
		return s.o, true, nil
	}

	return s.trySetState(Canceled)
}

// trySetState tries set state to the order.
// If state to replace is behind current state, ErrCannotSetState is returned.
func (s *StateOperator) trySetState(state State) (*Order, bool, error) {
	if s.o.IsClosed() {
		return nil, false, ErrOrderClosed
	}

	if state == Canceled || state == Closed {
		s.setState(state)
		return s.o, true, nil
	}

	orderState := s.o.state

	if orderState.Next.Name != state.Name {
		return nil, false, ErrCannotSetState
	}

	s.nextState()

	return s.o, true, nil
}

// nextState sets order`s state to the next
func (s *StateOperator) nextState() {
	s.o.state = *s.o.state.Next
}

// setState sets provided state to the current order.
func (s *StateOperator) setState(state State) {
	s.o.state = state
}

// getStateIndexes returns current order`s state index and replacing state index from stateList.
func (s *StateOperator) getStateIndexes(state State) (currentStateIndex, replacingState int, err error) {
	currentStateIndex = slices.Index(stateList, s.o.state)
	replacingState = slices.Index(stateList, state)

	if currentStateIndex == -1 || replacingState == -1 {
		return -1, -1, ErrNoSuchState
	}

	return
}
