package order

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

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
// If order is closed, it returns an error.
func (s *StateOperator) CancelOrder() (bool, error) {
	return s.trySetState(Canceled)
}

// CloseOrder set orders`s state to [Closed].
// If order is closed, it returns a nil error.
func (s *StateOperator) CloseOrder() (bool, error) {
	return s.trySetState(Closed)
}

// PayOrder set orders`s state to [Paid].
// If order is closed, it returns an error.
func (s *StateOperator) PayOrder(transactionID uuid.UUID) (bool, error) {
	set, err := s.trySetState(Paid)
	if err != nil || set == false {
		return false, err
	}

	s.o.transactionID = transactionID

	return true, nil
}

// CookOrder set orders`s state to [Cooking].
// If order is closed, it returns an error.
func (s *StateOperator) CookOrder() (bool, error) {
	return s.trySetState(Cooking)
}

// OrderFinished set orders`s state to [Finished].
// If order is closed, it returns an error.
func (s *StateOperator) OrderFinished() (bool, error) {
	return s.trySetState(Finished)
}

// WaitForCourier set orders`s state to [WaitingForCourier].
// If order is closed, it returns an error.
func (s *StateOperator) WaitForCourier() (bool, error) {
	return s.trySetState(WaitingForCourier)
}

// CourierTookOrder set orders`s state to [CourierTook].
// If order is closed, it returns an error.
func (s *StateOperator) CourierTookOrder() (bool, error) {
	return s.trySetState(CourierTook)
}

// DeliveringOrder set orders`s state to [Delivering].
// If order is closed, it returns an error.
func (s *StateOperator) DeliveringOrder() (bool, error) {
	return s.trySetState(Delivering)
}

// OrderDelivered set orders`s state to [Delivered].
// If order is closed, it returns an error.
func (s *StateOperator) OrderDelivered() (bool, error) {
	return s.trySetState(Delivered)
}

// trySetState tries set state to the order.
// If next State equals to Order`s state. The returned boolean will be true and error is nil.
// As if verb of the State were done.
// If next state is [Canceled] or [Closed] it sets it immediately.
// Otherwise, it checks if the next state is the same as the current order state.
// If it is, it sets the next state.
func (s *StateOperator) trySetState(next State) (bool, error) {
	orderState := s.o.state

	if orderState == next {
		return true, nil
	}

	if s.o.Is(Closed) {
		return false, errors.Wrapf(ErrOrderClosed, "cannot set state: %s", next.Name)
	}

	if s.o.Is(Canceled) {
		return false, errors.Wrapf(ErrOrderCanceled, "cannot set state: %s", next.Name)
	}

	if next == Canceled || next == Closed {
		s.setState(next)
		return true, nil
	}

	if orderState.Next.Name != next.Name {
		return false, errors.Wrapf(ErrInvalidState, "cannot set state %#v", next)
	}

	s.nextState()

	return true, nil
}

// nextState sets order`s state to the next
func (s *StateOperator) nextState() {
	s.o.state = *s.o.state.Next
}

// setState sets provided state to the current order.
func (s *StateOperator) setState(state State) {
	s.o.state = state
}
