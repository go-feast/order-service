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
) (err error) {
	var repositoryError, marshallError, publishError error

	defer func() {
		if marshallError != nil || publishError != nil {
			err = ob.repository.Delete(ctx, o)
		}
	}()

	repositoryError = ob.repository.Create(ctx, o)
	if repositoryError != nil {
		return errors.Wrap(repositoryError, "outbox: saving: failed to create order")
	}

	bytes, marshallError := ob.marshaller.Marshal(
		o.ToEvent().JSONEventOrderCreated())
	if marshallError != nil {
		return errors.Wrap(marshallError, "outbox: saving: failed to marshal event")
	}

	msg := message.NewMessage(uuid.NewString(), bytes)

	msg.SetContext(ctx)

	publishError = ob.publisher.Publish(topics.OrderCreated.String(), msg)
	if publishError != nil {
		return errors.Wrap(publishError, "outbox: saving: failed to publish event")
	}

	return nil
}
