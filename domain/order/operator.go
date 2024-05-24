package order

// StateOperator provides methods to operate with order state.
// USe StateOperator to operate on order state.
type StateOperator struct {
	o *Order
}

// NewStateOperator creates a new StateOperator.
func NewStateOperator(o *Order) *StateOperator {
	return &StateOperator{o: o}
}

// CancelOrder set orders`s state to [Canceled].
// If order is already canceled, it returns the order with canceled flag set to true.
// If order is closed, it returns an error.
//
// Returning values: order, canceled(bool), error
func (s *StateOperator) CancelOrder() (*Order, bool, error) {
	if s.o.IsCanceled() {
		return s.o, true, nil
	}

	return s.trySetState(Canceled)
}

// trySetState tries set state to the order.
// If next state is [Canceled] or [Closed] it sets it immediately.
// Otherwise, it checks if the next state is the same as the current order state.
// If it is, it sets the next state.
// Returning values: Order, changed(bool), error
func (s *StateOperator) trySetState(next State) (*Order, bool, error) {
	if s.o.IsClosed() {
		return nil, false, ErrOrderClosed
	}

	if next == Canceled || next == Closed {
		s.setState(next)
		return s.o, true, nil
	}

	orderState := s.o.state

	if orderState.Next.Name != next.Name {
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
