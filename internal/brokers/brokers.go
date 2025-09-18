package brokers

import (
	"fmt"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/brokers/nats"
	"github.com/LincolnG4/iot-hydra/internal/message"
	"gopkg.in/yaml.v3"
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
	Publish(*message.Message) error

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

// NewBroker returns a Broker interface based on the config type (nats, iothub,...)
func NewBroker(cfg Config) (Broker, error) {
	switch cfg.Type {
	case nats.NATSType:
		return nats.NewBroker(nats.Config{
			Name: cfg.Name,
			URL:  cfg.Address,
			Auth: cfg.Auth,
		}), nil
	default:
		return nil, fmt.Errorf("broker type %s not allowed", cfg.Type)
	}
}

// Custom unmarshaling for Config
func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	// Create a temporary struct with raw auth data
	var raw struct {
		Name    string                 `yaml:"name"`
		Type    string                 `yaml:"type"`
		Address string                 `yaml:"address"`
		Auth    map[string]interface{} `yaml:"auth"`
	}

	if err := value.Decode(&raw); err != nil {
		return err
	}

	c.Name = raw.Name
	c.Type = raw.Type
	c.Address = raw.Address

	// Handle auth based on method
	if raw.Auth != nil {
		method, ok := raw.Auth["method"].(string)
		if !ok {
			return fmt.Errorf("auth method must be a string")
		}

		switch method {
		case auth.BasicType:
			// TODO: fix
			user, _ := raw.Auth["user"].(string)
			password, _ := raw.Auth["password"].(string)
			c.Auth = &auth.BasicAuth{
				Username: user,
				Password: password,
			}
		case auth.TokenType:
			token, _ := raw.Auth["token"].(string)
			c.Auth = &auth.Token{
				Token: token,
			}
		default:
			return fmt.Errorf("unsupported auth method: %s", method)
		}
	}

	return nil
}
