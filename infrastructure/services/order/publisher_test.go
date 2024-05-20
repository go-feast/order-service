package order

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-feast/topics"
	"github.com/stretchr/testify/assert"
	"service/domain/order"
	"sync"
	"testing"
)

// TODO: integration testing with kafka

func TestPublisherService_PublishOrderCreated(t *testing.T) {
	t.Run("assert publisher creates order", func(t *testing.T) {
		var (
			logger = watermill.NopLogger{}
		)

		ch := gochannel.NewGoChannel(gochannel.Config{}, logger)

		service := NewPublisherService(ch)

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

		ctx, cancel := context.WithCancel(context.Background())

		messages, err := ch.Subscribe(ctx, topics.OrderCreated.String())

		assert.NoError(t, err)

		var wg sync.WaitGroup

		wg.Add(1)

		go func() {
			defer wg.Done()

			msg := <-messages

			msg.Ack()

			cancel()

			var p PostOrder

			err = json.Unmarshal(msg.Payload, &p)

			assert.NoError(t, err)

			assert.Equal(t, o.ID.GetID(), p.ID)
			assert.Equal(t, o.RestaurantID.GetID(), p.RestaurantID)
			assert.Equal(t, o.CustomerID.GetID(), p.CustomerID)
			assert.ElementsMatch(t, o.GetMealsID(), p.Meals)
			assert.Equal(t, o.Destination.Latitude, p.Destination.Lat)
			assert.Equal(t, o.Destination.Longitude, p.Destination.Long)
		}()

		err = service.PublishOrderCreated(o)

		assert.NoError(t, err)

		wg.Wait()
	})
}
