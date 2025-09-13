package brokers

import (
	"fmt"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/brokers/nats"
	"github.com/LincolnG4/iot-hydra/internal/message"
)

type Broker interface {
	// Broker type
	Type() string

	// Connect to the broker
	Connect() error

	// Stop the connection to the broker
	Stop() error

	// Publish the message to the broker
	Publish(*message.Message) error

	// Subscribe to broker and wait T seconds to receive the message, otherwise
	// returns nil and timeout
	SubscribeAndWait(string, time.Duration) (*message.Message, error)
}

type Config struct {
	Type    string
	Address string
	Auth    auth.Authenticator
}

// NewBroker returns a Broker interface based on the config type (nats, iothub,...)
func NewBroker(cfg Config) (Broker, error) {
	switch cfg.Type {
	case nats.NATSType:
		return nats.NewBroker(nats.Config{
			URL:  cfg.Address,
			Auth: cfg.Auth,
		}), nil
	default:
		return &nats.NATS{}, fmt.Errorf("broker type %s not allowed", cfg.Type)
	}
}
