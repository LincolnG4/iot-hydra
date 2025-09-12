package nats

import "github.com/nats-io/nats.go"

type Connector interface {
	Publish(string, []byte) error
	SubscribeSync(string) (*nats.Subscription, error)

	Close()
}
