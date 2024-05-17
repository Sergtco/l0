package data

import (
	"encoding/json"

	"service/data/cache"
	"service/models"

	"github.com/nats-io/nats.go/jetstream"
)

var Cache *cache.Cache

func init() {
	Cache = cache.NewCache(1024)
}

func Process(message jetstream.Msg) error {
	data := message.Data()
	order := &models.Order{}
	err := json.Unmarshal(data, order)
	if err != nil {
		return err
	}
	Cache.Store(order.Id, order)
	message.Ack()
	return nil
}

func GetOrder(id string) (*models.Order, error) {
	order, err := Cache.Get(id)
	if err != nil {
		return nil, err
	}
	return order, nil
}
