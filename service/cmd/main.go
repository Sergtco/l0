package main

import (
	"context"
	"service/api"
	"service/nats"
)

func main() {
	bg := context.Background()
	ctx, cancelFunc := context.WithCancel(bg)
	defer cancelFunc()
	go func() {
		if err := nats.StartListening(ctx); err != nil {
			panic(err)
		}
	}()
	if err := api.RunServer(); err != nil {
		panic(err)
	}

}
