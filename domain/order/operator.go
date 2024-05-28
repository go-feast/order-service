package order

import "github.com/pkg/errors"

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
// If order is already setted, it returns the order with setted flag set to true.
// If order is closed, it returns an error.
func (s *StateOperator) CancelOrder() (*Order, bool, error) {
	return s.trySetState(Canceled)
}

// CloseOrder set orders`s state to [Closed].
// If order is already setted, it returns the order with setted flag set to true.
// If order is closed, it returns a nil error.
func (s *StateOperator) CloseOrder() (*Order, bool, error) {
	return s.trySetState(Closed)
}

// PayOrder set orders`s state to [Paid].
// If order is already setted, it returns the order with setted flag set to true.
// If order is closed, it returns an error.
func (s *StateOperator) PayOrder() (*Order, bool, error) {
	return s.trySetState(Paid)
}

// CookOrder set orders`s state to [Cooking].
// If order is already setted, it returns the order with setted flag set to true.
// If order is closed, it returns an error.
func (s *StateOperator) CookOrder() (*Order, bool, error) {
	return s.trySetState(Cooking)
}

// OrderFinished set orders`s state to [Finished].
// If order is already setted, it returns the order with setted flag set to true.
// If order is closed, it returns an error.
func (s *StateOperator) OrderFinished() (*Order, bool, error) {
	return s.trySetState(Finished)
}

// WaitForCourier set orders`s state to [WaitingForCourier].
// If order is already setted, it returns the order with setted flag set to true.
// If order is closed, it returns an error.
func (s *StateOperator) WaitForCourier() (*Order, bool, error) {
	return s.trySetState(WaitingForCourier)
}

// CourierTookOrder set orders`s state to [CourierTook].
// If order is already setted, it returns the order with setted flag set to true.
// If order is closed, it returns an error.
func (s *StateOperator) CourierTookOrder() (*Order, bool, error) {
	return s.trySetState(CourierTook)
}

// DeliveringOrder set orders`s state to [Delivering].
// If order is already setted, it returns the order with setted flag set to true.
// If order is closed, it returns an error.
func (s *StateOperator) DeliveringOrder() (*Order, bool, error) {
	return s.trySetState(Delivering)
}

// OrderDelivered set orders`s state to [Delivered].
// If order is already setted, it returns the order with setted flag set to true.
// If order is closed, it returns an error.
func (s *StateOperator) OrderDelivered() (*Order, bool, error) {
	return s.trySetState(Delivered)
}

// trySetState tries set state to the order.
// If next State equals to Order`s state. The returned boolean will be true and error is nil.
// As if verb of the State were done.
// If next state is [Canceled] or [Closed] it sets it immediately.
// Otherwise, it checks if the next state is the same as the current order state.
// If it is, it sets the next state.
func (s *StateOperator) trySetState(next State) (*Order, bool, error) {
	orderState := s.o.state

	if orderState == next {
		return s.o, true, nil
	}

	if s.o.Is(Closed) {
		return nil, false, errors.Wrapf(ErrOrderClosed, "cannot set state: %s", next.Name)
	}

	if s.o.Is(Canceled) {
		return s.o, false, errors.Wrapf(ErrOrderCanceled, "cannot set state: %s", next.Name)
	}

	if next == Canceled || next == Closed {
		s.setState(next)
		return s.o, true, nil
	}

	if orderState.Next.Name != next.Name {
		return nil, false, errors.Wrapf(ErrInvalidState, "cannot set state %#v", next)
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
