package order

// State represents the current state of the order.
// It contains the name of the state and the next state.
// State is a value object.
type State struct { //nolint:govet
	Name string
	Next *State
}

func (s State) String() string { return s.Name }

// State machine for an order:
// Created -> Paid -> Cooking -> Finished -> WaitingForCourier -> CourierTook -> Delivering -> Delivered -> Closed.
// Every State can go into Canceled State. But the only way where Canceled can go into is Closed.
// Canceled -> Closed.
var (
	Canceled = State{"order.canceled", &Closed}

	Created = State{"order.created", &Paid}

	Paid = State{"order.paid", &Cooking}

	Cooking = State{"order.cooking", &Finished}

	Finished = State{"order.cooking.finished", &WaitingForCourier}

	WaitingForCourier = State{"order.waiting", &CourierTook}

	CourierTook = State{"order.taken", &Delivering}

	Delivering = State{"order.delivering", &Delivered}

	Delivered = State{"order.delivered", &Closed}

	Closed = State{"order.closed", nil}
)

// mapStates is a map of order state names and order states.
var mapStates = map[string]State{ //nolint:unused
	"order.created":          Created,
	"order.paid":             Paid,
	"order.cooking":          Cooking,
	"order.cooking.finished": Finished,
	"order.waiting":          WaitingForCourier,
	"order.taken":            CourierTook,
	"order.delivering":       Delivering,
	"order.delivered":        Delivered,
	"order.closed":           Closed,
}
