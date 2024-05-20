package order

type Publisher interface {
	PublishOrderCreated(o *Order) error
}
