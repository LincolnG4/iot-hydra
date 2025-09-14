package agent

import (
	"fmt"
	"testing"

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
        password: test`

	var wrapper struct {
		TelemetryAgent TelemetryAgent `yaml:"telemetryAgent"`
	}
	err := yaml.Unmarshal([]byte(y), &wrapper)
	assert.NoError(t, err)
	fmt.Println(wrapper.TelemetryAgent)
}
