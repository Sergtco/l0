package data

import (
	"encoding/json"

	"service/data/cache"
	"service/models"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go/jetstream"
)

var Cache *cache.Cache
var validate validator.Validate

func init() {
	Cache = cache.NewCache(1024)
	validate = *validator.New(validator.WithRequiredStructEnabled())
}

// Processes message. Either accepts it or declines.
func Process(message *jetstream.Msg) error {
	data := (*message).Data()
	order := &models.Order{}
	err := json.Unmarshal(data, order)
	if err != nil {
		return err
	}
	err = validate.Struct(order)
	if err != nil {
		return err
	}

	err = Cache.Store(order)
	if err != nil {
		return err
	}
	(*message).Ack()
	return nil
}

// REturns Order. Else orror returned.
func GetOrder(id string) (*models.Order, error) {
	order, err := Cache.Get(id)
	if err != nil {
		return nil, err
	}
	return order, nil
}
