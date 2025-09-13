package brokers

import (
	"testing"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/brokers/nats"
	"github.com/alecthomas/assert"
)

func TestNewBroker(t *testing.T) {
	cfg := Config{
		Type:    nats.NATSType,
		Address: "nats://localhost:4222",
		Auth: auth.BasicAuth{
			Username: "foo",
			Password: "bar",
		},
	}

	b, err := NewBroker(cfg)
	assert.NoError(t, err, "Could not create nats")

	assert.Equal(t, nats.NATSType, b.Type(), "Broker type doesn't match")
}
