package brokers

import (
	"context"
	"fmt"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/brokers/nats"
	"github.com/LincolnG4/iot-hydra/internal/message"
)

type Broker interface {
	// Unique name
	Name() string
	// Broker type
	Type() string

	// Connect to the broker
	Connect() error

	// Stop the connection to the broker
	Stop() error

	// Publish the message to the broker
	Publish(context.Context, *message.Message) error

	// Subscribe to broker and wait T seconds to receive the message, otherwise
	// returns nil and timeout
	SubscribeAndWait(string, time.Duration) (*message.Message, error)
}

type Config struct {
	Name    string             `yaml:"name" validate:"required,min=1,max=255"`
	Type    string             `yaml:"type" validate:"required"`
	Address string             `yaml:"address" validate:"required"`
	Auth    auth.Authenticator `yaml:"auth"`
}

// NewBroker is a factory that returns a specific broker implementation.
func NewBroker(cfg Config) (Broker, error) {
	switch cfg.Type {
	case nats.NATSType:
		return nats.NewBroker(nats.Config{
			Name: cfg.Name,
			URL:  cfg.Address,
			Auth: cfg.Auth,
		}), nil
	default:
		return nil, fmt.Errorf("broker type '%s' is not supported", cfg.Type)
	}
}
