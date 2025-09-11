package brokers

import (
	"fmt"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/brokers/nats"
	"github.com/LincolnG4/iot-hydra/internal/message"
)

const (
	TypeNATS = "nats"
)

type Broker interface {
	Connect() error
	Stop() error
	Publish(message.Message) error
}

type Config struct {
	Type    string
	Address string
	Auth    auth.Authenticator
}

func NewBroker(cfg Config) (Broker, error) {
	switch cfg.Type {
	case TypeNATS:
		return nats.NewBroker(nats.Config{
			URL:  cfg.Address,
			Auth: cfg.Auth,
		}), nil
	default:
		return &nats.NATS{}, fmt.Errorf("broker type %s not allowed", cfg.Type)
	}
}
