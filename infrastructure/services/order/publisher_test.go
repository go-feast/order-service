package order

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-feast/topics"
	"github.com/stretchr/testify/assert"
	"service/domain/order"
	"sync"
	"testing"
	"time"
)

// TODO: integration testing with kafka

func pubsubMessageSending(ctx context.Context, topic string, pub order.Publisher, sub message.Subscriber) func(t *testing.T) {
	return func(t *testing.T) {
		messages, err := sub.Subscribe(ctx, topic)

		assert.NoError(t, err)

		var wg sync.WaitGroup

		wg.Add(1)

		var receivedCounter int

		go func(messages <-chan *message.Message) {
			defer wg.Done()

			for msg := range messages {
				var p PostOrder

				err := json.Unmarshal(msg.Payload, &p)

				assert.NoError(t, err)

				msg.Ack()

				receivedCounter++
			}
		}(messages)

		var sentCounter int

		wg.Add(1)

		go func(ctx context.Context) {
			defer wg.Done()

			for {
				o, err := order.NewOrder(
					gofakeit.UUID(),
					gofakeit.UUID(),
					gofakeit.UUID(),
					[]string{
						gofakeit.UUID(),
						gofakeit.UUID(),
					},
					false,
					gofakeit.Latitude(),
					gofakeit.Longitude(),
				)

				assert.NoError(t, err)

				o.Create()

				select {
				case <-ctx.Done():
					return
				default:
				}

				err = pub.PublishOrderCreated(o)
				assert.NoError(t, err)

				sentCounter++

				select {
				case <-ctx.Done():
					return
				default:
				}
			}
		}(ctx)

		// TODO: ask
		// receivedCounter == sentCounter

		wg.Wait()
	}
}

func TestPublisherService_PublishOrderCreated_gochannels(t *testing.T) {
	var (
		logger  = watermill.NopLogger{}
		ch      = gochannel.NewGoChannel(gochannel.Config{}, logger)
		service = NewPublisherService(ch)
	)

	t.Cleanup(func() {
		err := ch.Close()

		assert.NoError(t, err)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	t.Run(
		"assert publisher creates order",
		pubsubMessageSending(ctx, topics.OrderCreated.String(), service, ch),
	)
}
