package nats

import (
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

func (n *NATS) Publish(msg *message.Message) error {
	// TODO: Add from where the message came
	return n.conn.Publish(msg.Topic, msg.Payload)
}

func (n *NATS) SubscribeAndWait(topic string, waitSecond time.Duration) (*message.Message, error) {
	s, err := n.conn.SubscribeSync(topic)
	if err != nil {
		return nil, err
	}
	defer s.Unsubscribe()

	msg, err := s.NextMsg(waitSecond)
	if err != nil {
		return nil, err
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
	case *auth.Token:
		natsOpts = append(natsOpts, nats.Token(authConfig.Token))
	default:
		return nil, fmt.Errorf("method %s not allowed", authConfig)
	}

	return natsOpts, nil
}
