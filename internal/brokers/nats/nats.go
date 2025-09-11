package nats

import (
	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/message"
	"github.com/nats-io/nats.go"
)

type NATS struct {
	conn        *nats.Conn
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
	URL  string             `json:"url"`
	Auth auth.Authenticator `json:"auth"`
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

	conn, err := nats.Connect(n.config.URL, ncOpts...)
	if err != nil {
		return err
	}

	n.conn = conn
	n.isConnected = true
	return nil
}

func (n *NATS) Stop() error {
	if n.conn != nil {
		n.conn.Close()
	}
	return nil
}

func (n *NATS) Publish(msg message.Message) error {
	return nil
}
