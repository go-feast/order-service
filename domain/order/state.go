package order

type State struct { //nolint:govet
	Name string
	Next *State
}

func (s *State) String() string {
	if s == nil {
		panic("state can`t be nil")
	}

	return s.Name
}

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
