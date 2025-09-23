package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/message"
	"github.com/nats-io/nats.go"
)

const (
	NATSType = "nats"
)

type NATS struct {
	conn        Connector
	isConnected bool
	Config      Config
}

func NewBroker(cfg Config) *NATS {
	return &NATS{
		isConnected: false,
		Config:      cfg,
	}
}

type Config struct {
	Name string             `json:"name" yaml:"name"`
	URL  string             `json:"url" yaml:"url"`
	Auth auth.Authenticator `json:"auth" yaml:"auth"`
}

func (n *NATS) Name() string {
	return n.Config.Name
}

func (n *NATS) Type() string {
	return NATSType
}

func (n *NATS) Connect() error {
	natsOpts, err := getCredentials(n.Config.Auth)
	if err != nil {
		return err
	}

	n.conn, err = nats.Connect(n.Config.URL, natsOpts...)
	if err != nil {
		return err
	}

	n.isConnected = true
	return nil
}

func (n *NATS) Stop() error {
	if n.conn != nil {
		n.conn.Close()
	}
	return nil
}

func (n *NATS) Publish(ctx context.Context, msg *message.Message) error {
	if n.conn == nil {
		return fmt.Errorf("NATS connection is not established for broker '%s'", n.Config.Name)
	}

	if !n.isConnected {
		return fmt.Errorf("NATS broker '%s' is not connected", n.Config.Name)
	}

	if err := n.conn.Publish(msg.Topic, msg.Payload); err != nil {
		return fmt.Errorf("failed to publish message to topic '%s' on broker '%s': %w", msg.Topic, n.Config.Name, err)
	}

	return nil
}

func (n *NATS) SubscribeAndWait(topic string, waitSecond time.Duration) (*message.Message, error) {
	if n.conn == nil {
		return nil, fmt.Errorf("NATS connection is not established for broker '%s'", n.Config.Name)
	}

	if !n.isConnected {
		return nil, fmt.Errorf("NATS broker '%s' is not connected", n.Config.Name)
	}

	s, err := n.conn.SubscribeSync(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic '%s' on broker '%s': %w", topic, n.Config.Name, err)
	}
	defer s.Unsubscribe()

	msg, err := s.NextMsg(waitSecond)
	if err != nil {
		return nil, fmt.Errorf("failed to receive message from topic '%s' on broker '%s': %w", topic, n.Config.Name, err)
	}

	return &message.Message{
		Payload:      msg.Data,
		Topic:        topic,
		SourceBroker: n.Name(),
	}, nil
}

// getCredentials identify the type of authentication and returns the credentials for the broker.
func getCredentials(a auth.Authenticator) ([]nats.Option, error) {
	var natsOpts []nats.Option

	switch authConfig := a.(type) {
	case *auth.BasicAuth:
		natsOpts = append(natsOpts, nats.UserInfo(authConfig.Username, authConfig.Password))
	case *auth.TokenAuth:
		natsOpts = append(natsOpts, nats.Token(authConfig.Token))
	default:
		return nil, fmt.Errorf("method %s not allowed", authConfig)
	}

	return natsOpts, nil
}
