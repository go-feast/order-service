package outbox

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/topics"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"service/domain/order"
	"service/event"
)

type Outbox struct {
	publisher  message.Publisher
	marshaller event.Marshaler
	repository order.Repository
}

func NewOutbox(
	publisher message.Publisher,
	repository order.Repository,
	marshaller event.Marshaler,
) *Outbox {
	return &Outbox{
		publisher:  publisher,
		marshaller: marshaller,
		repository: repository,
	}
}

func (ob *Outbox) Save(
	ctx context.Context,
	o *order.Order,
) error {
	err := ob.repository.Create(ctx, o)
	if err != nil {
		return errors.Wrap(err, "outbox: saving: failed to create order")
	}

	bytes, err := ob.marshaller.Marshal(
		o.ToEvent().ToJSON())
	if err != nil {
		return errors.Wrap(err, "outbox: saving: failed to marshal event")
	}

	msg := message.NewMessage(uuid.NewString(), bytes)

	msg.SetContext(ctx)

	err = ob.publisher.Publish(topics.OrderCreated.String(), msg)
	if err != nil {
		return errors.Wrap(err, "outbox: saving: failed to publish event")
	}

	return nil
}
