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
	// A temporary struct to decode the YAML into
	var raw struct {
		Name    string                 `yaml:"name"`
		Type    string                 `yaml:"type"`
		Address string                 `yaml:"address"`
		Auth    map[string]interface{} `yaml:"auth"`
	}

	if err := value.Decode(&raw); err != nil {
		return err
	}

	if raw.Name == "" {
		return fmt.Errorf("broker config requires a 'name'")
	}
	if raw.Type == "" {
		return fmt.Errorf("broker config requires a 'type'")
	}
	if raw.Address == "" {
		return fmt.Errorf("broker config requires an 'address'")
	}

	c.Name = raw.Name
	c.Type = raw.Type
	c.Address = raw.Address

	// Handle auth based on method
	if raw.Auth != nil {
		method, ok := raw.Auth["method"].(string)
		if !ok {
			return fmt.Errorf("auth 'method' must be a string")
		}

		switch method {
		case auth.BasicType:
			user, ok := raw.Auth["user"].(string)
			if !ok {
				return fmt.Errorf("auth method 'plain' requires a 'user' string")
			}

			password, ok := raw.Auth["password"].(string)
			if !ok {
				return fmt.Errorf("auth method 'plain' requires a 'password' string")
			}

			c.Auth = &auth.BasicAuth{
				Username: user,
				Password: password,
			}

		case auth.TokenType:
			token, ok := raw.Auth["token"].(string)
			if !ok {
				return fmt.Errorf("auth method 'token' requires a 'token' string")
			}

			c.Auth = &auth.Token{
				Token: token,
			}
		default:
			return fmt.Errorf("unsupported auth method: %s", method)
		}
	}

	return nil
}
