package nats

import (
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
	config      Config
}

func NewBroker(cfg Config) *NATS {
	return &NATS{
		isConnected: false,
		config:      cfg,
	}
}

type Config struct {
	Name string             `json:"name" yaml:"name"`
	URL  string             `json:"url" yaml:"url"`
	Auth auth.Authenticator `json:"auth" yaml:"auth"`
}

func (n *NATS) Name() string {
	return n.config.Name
}

func (n *NATS) Type() string {
	return NATSType
}

func (n *NATS) Connect() error {
	opts := &auth.ConnectOptions{}
	if n.config.Auth != nil {
		if err := n.config.Auth.Apply(opts); err != nil {
			return err
		}
	}

	// Build NATS options
	ncOpts := []nats.Option{}
	if opts.Token != "" {
		ncOpts = append(ncOpts, nats.Token(opts.Token))
	}
	if opts.Username != "" && opts.Password != "" {
		ncOpts = append(ncOpts, nats.UserInfo(opts.Username, opts.Password))
	}

	var err error
	n.conn, err = nats.Connect(n.config.URL, ncOpts...)
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
		Payload: msg.Data,
		Topic:   topic,
		Type:    NATSType,
	}, nil
}
