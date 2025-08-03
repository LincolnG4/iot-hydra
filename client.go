package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type IoTRuntime struct {
	nats      *nats.Conn
	jetstream jetstream.JetStream
}

func startIoTRuntime(url string) IoTRuntime {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal("could not connect to server", err)
	}
	defer nc.Drain()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}
	err = setupStream(js)
	if err != nil {
		log.Fatal(err)
	}
	return IoTRuntime{
		nats:      nc,
		jetstream: js,
	}
}

func setupStream(js jetstream.JetStream) error {
	_, err := js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name: "iot",
	})
	if err != nil {
		if err == jetstream.ErrStreamNameAlreadyInUse {
			return nil
		}
		return fmt.Errorf("CreateStream failed: %s", err)
	}
	return nil
}

func main() {
	rt := startIoTRuntime("dfjijddi24uhjd29834ijrr0345jo0r3j034n@localhost:4222")

	rt.jetstream.Publish(context.Background(), "messages", []byte("test"))
}
