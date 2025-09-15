package agent

import (
	"testing"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/brokers/nats"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestUnmarshalYAML(t *testing.T) {
	y := `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: ligmaMyNats
      type: nats
      address: "localhost:4222"
      auth:
        method: plain 
        user: test
        password: testpass`

	var wrapper struct {
		TelemetryAgent TelemetryAgent `yaml:"telemetryAgent"`
	}
	err := yaml.Unmarshal([]byte(y), &wrapper)
	assert.NoError(t, err)

	// check queue size
	assert.Equal(t, 1000, wrapper.TelemetryAgent.QueueSize, "Wrong queueSize")

	// check if the nats was created on the map
	b, ok := wrapper.TelemetryAgent.Brokers["ligmaMyNats"]
	assert.Equal(t, true, ok)

	assert.Equal(t, auth.BasicType, b.AuthMethod(), "wrong authentication method expected")

	natsBroker, ok := b.(*nats.NATS)
	assert.True(t, ok, "Broker should be a NATS broker")
	assert.Equal(t, "localhost:4222", natsBroker.Config.URL, "Wrong queueSize")

	// assert auth
	basicAuth, ok := natsBroker.Config.Auth.(auth.BasicAuth)
	assert.True(t, ok, "Auth method should be BasicAuth")
	assert.Equal(t, "test", basicAuth.Username, "Wrong username")
	assert.Equal(t, "testpass", basicAuth.Password, "Wrong password")

	assert.Equal(t, b.Type(), nats.NATSType)
}
