package nats

import (
	"context"
	"log"
	"os"
	"service/data"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func StartListening(parentCtx context.Context) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
	defer cancel()
	url := os.Getenv("NATS_URL")
	if len(url) == 0 {
		url = "127.0.0.1"
	}
	nc, err := nats.Connect(url)
	for err != nil {
		log.Println(err)
		time.Sleep(time.Second)
		nc, err = nats.Connect(url)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}
	cons, err := js.CreateOrUpdateConsumer(ctx, "ORDERS", jetstream.ConsumerConfig{
		Name:          "service",
		AckPolicy:     jetstream.AckExplicitPolicy,
		Durable:       "service",
		FilterSubject: "ORDERS.new",
	})
	if err != nil {
		return err
	}
	iter, err := cons.Messages()
	if err != nil {
		return err
	}
	defer iter.Stop()
	for {
		msg, err := iter.Next()
		if err != nil {
			return err
		}
		go func() {
			err := data.Process(&msg)
			if err != nil {
				log.Println("In nats: ", err)
				return
			}
			log.Println("Successfuly parsed message ", msg)
		}()
	}
}
