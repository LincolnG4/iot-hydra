package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

type IoTAgent struct {
	config *ServiceConfig
}

func NewIoTAgent(cfg *ServiceConfig) IoTAgent {
	return IoTAgent{
		config: cfg,
	}
}

type ServiceConfig struct {
	Azure  AzureConfig  `json:"azure"`
	AWS    AWSConfig    `json:"aws"`
	Nats   NATSConfig   `json:"nats"`
	Custom CustomConfig `json:"custom"`
}

type NATSConfig struct {
	URL       string `json:"url"`
	Topic     string
	BasicAuth *BasicAuth
}

type BasicAuth struct {
	Username string
	Password string
}

type AzureConfig struct {
	ConnectionString string `json:"connection_string"`
	DeviceID         string `json:"device_id"`
	Enabled          bool   `json:"enabled"`
}

type AWSConfig struct {
	Region    string `json:"region"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	ThingName string `json:"thing_name"`
	Enabled   bool   `json:"enabled"`
}

type CustomConfig struct {
	Endpoint string            `json:"endpoint"`
	Headers  map[string]string `json:"headers"`
	Enabled  bool              `json:"enabled"`
}

type IoTMessage struct {
	ID        string                 `json:"id"`
	DeviceID  string                 `json:"device_id"`
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	// azure, nats, ...
	Target []string `json:"target"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ag *IoTAgent) Route(msg *IoTMessage) error {
	target := msg.Target[0]
	switch target {
	case "nats":
		return ag.Send(msg)
	default:
		return fmt.Errorf("target %s is not allowed", target)
	}
}

func (ag *IoTAgent) Send(msg *IoTMessage) error {
	nc, err := nats.Connect(msg.Target[0])
	if err != nil {
		return err
	}
	return nc.Publish("my.iot", []byte("meu bilau"))
}

func (a *application) websocketIoTHandler(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		a.logger.Error().Err(fmt.Errorf("failed to set websocket upgrade: %+v", err)).Msg("")
		return
	}

	a.logger.Info().Msg("Client connected to Web Socket")
	for {
		msg := IoTMessage{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			a.logger.Error().Err(fmt.Errorf("WebSocket error: %v", err)).Msg("")
			break
		}

		msg.ID = fmt.Sprintf("ws-%d", time.Now().UnixNano())
		msg.Timestamp = time.Now()
		select {
		case a.MessageQueue <- msg:
			a.logger.Debug().Msg("Message received via WebSocket:" + msg.ID)
		default:
			a.logger.Debug().Msg("Message queue full, dropping message: %s" + msg.ID)
		}
	}
}
