package order

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/topics"
	"github.com/google/uuid"
	"service/domain/order"
)

type PublisherService struct {
	publisher message.Publisher
}

func NewPublisherService(publisher message.Publisher) *PublisherService {
	return &PublisherService{publisher: publisher}
}

type PostOrder struct {
	ID           string   `json:"id"`
	RestaurantID string   `json:"restaurant_id"`
	CustomerID   string   `json:"customer_id"`
	Meals        []string `json:"meals"`
	Destination  struct {
		Lat, Long float64
	} `json:"destination"`
}

func (p *PublisherService) PublishOrderCreated(o *order.Order) error {
	po := PostOrder{
		ID:           o.ID.GetID(),
		RestaurantID: o.RestaurantID.GetID(),
		CustomerID:   o.CustomerID.GetID(),
		Meals:        o.GetMealsID(),
		Destination: struct{ Lat, Long float64 }{
			Lat:  o.Destination.Latitude,
			Long: o.Destination.Longitude,
		},
	}

	// TODO: ask
	bytes, err := json.Marshal(po)
	if err != nil {
		return err
	}

	msg := message.NewMessage(uuid.NewString(), bytes)

	return p.publisher.Publish(topics.OrderCreated.String(), msg)
}
