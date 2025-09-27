package config

import (
	"testing"

	"github.com/alecthomas/assert"
	"gopkg.in/yaml.v3"
)

func TestUnmarshalYAML_Success(t *testing.T) {
	y := []byte(`
telemetryAgent:
  queueSize: 1000
  maxWorkers: 2
  brokers:
    - name: ligmaMyNats
      type: nats
      address: "localhost:4222"
      auth:
        method: plain
        user: test
        password: testpass
    - name: foo
      type: nats
      address: "localhost:9999"
      auth:
        method: token
        token: my-secret-token
`)
	var wrapper struct {
		TelemetryAgent TelemetryAgentYAML `yaml:"telemetryAgent"`
	}
	err := yaml.Unmarshal(y, &wrapper)
	assert.NoError(t, err)

	// check global properties
	assert.Equal(t, 1000, wrapper.TelemetryAgent.QueueSize, "Wrong queueSize")
	assert.Len(t, wrapper.TelemetryAgent.Brokers, 2, "Expected to parse 2 brokers")

	// Define test cases for each broker
	testCases := []struct {
		name             string
		brokerName       string
		expectedType     string
		expectedAddress  string
		expectedAuthType string
	}{
		{
			name:             "NATS broker with plain auth",
			brokerName:       "ligmaMyNats",
			expectedType:     "nats",
			expectedAddress:  "localhost:4222",
			expectedAuthType: "plain",
		},
		{
			name:             "NATS broker with token auth",
			brokerName:       "foo",
			expectedType:     "nats",
			expectedAddress:  "localhost:9999",
			expectedAuthType: "token",
		},
	}

	// Iterate through all test cases
	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Check if the broker was created on the map
			b := wrapper.TelemetryAgent.Brokers[i]
			assert.Equal(t, tc.brokerName, b.Name, "Wrong broker name")                 // Assert NAME
			assert.Equal(t, tc.expectedType, b.Type, "Wrong broker type")               // Assert Type
			assert.Equal(t, tc.expectedAddress, b.Address, "Wrong broker address")      // Assert Address
			assert.Equal(t, tc.expectedAuthType, b.Auth.Method, "Wrong broker address") // Assert Auth Method
		})
	}
}
