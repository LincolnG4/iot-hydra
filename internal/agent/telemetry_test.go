package agent

import (
	"context"
	"testing"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/brokers"
	"github.com/LincolnG4/iot-hydra/internal/message"
	"github.com/rs/zerolog/log"
)

// TODO: Create tests
func TestSendMessage_Success(t *testing.T) {
	cfg := brokers.Config{
		Name:    "nats",
		Type:    "nats",
		Address: "localhost:4222",
		Auth: &auth.BasicAuth{
			Username: "test",
			Password: "test",
		},
	}
	b, _ := brokers.NewBroker(cfg)
	ag := TelemetryAgent{
		Queue: make(chan *message.Message, 1000),
		Brokers: map[string]brokers.Broker{
			"nats": b,
		},
	}

	ag.WorkerPool.Start()
	ctx, cancel := context.WithCancel(context.Background())
	go func(context.CancelFunc) {
		time.Sleep(10 * time.Second)
		cancel()
	}(cancel)

	for {
		select {
		case msg := <-ag.Queue:
			if err := ag.Submit(msg); err {
				log.Error().Msg("queue should not be full")
			}
		case <-ctx.Done():
			log.Error().Msg("terminated")
		}
	}
}
