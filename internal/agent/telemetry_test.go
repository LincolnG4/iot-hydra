package agent

import (
	"testing"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/brokers"
	"github.com/LincolnG4/iot-hydra/internal/brokers/nats"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestUnmarshalYAML_Success(t *testing.T) {
	y := `telemetryAgent:
  queueSize: 1000
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
        method: natsToken
        token: my-secret-token`

	var wrapper struct {
		TelemetryAgent TelemetryAgent `yaml:"telemetryAgent"`
	}
	err := yaml.Unmarshal([]byte(y), &wrapper)
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
		assertion        func(t *testing.T, b brokers.Broker)
	}{
		{
			name:             "NATS broker with plain auth",
			brokerName:       "ligmaMyNats",
			expectedType:     "nats",
			expectedAddress:  "localhost:4222",
			expectedAuthType: "plain",
			assertion: func(t *testing.T, b brokers.Broker) {
				natsBroker, ok := b.(*nats.NATS)
				assert.True(t, ok, "Broker should be a NATS broker")

				basicAuth, ok := natsBroker.Config.Auth.(auth.BasicAuth)
				assert.True(t, ok, "Auth method should be BasicAuth")
				assert.Equal(t, "test", basicAuth.Username, "Wrong username")
				assert.Equal(t, "testpass", basicAuth.Password, "Wrong password")
			},
		},
		{
			name:             "NATS broker with token auth",
			brokerName:       "foo",
			expectedType:     "nats",
			expectedAddress:  "localhost:9999",
			expectedAuthType: "natsToken",
			assertion: func(t *testing.T, b brokers.Broker) {
				natsBroker, ok := b.(*nats.NATS)
				assert.True(t, ok, "Broker should be a NATS broker")

				tokenAuth, ok := natsBroker.Config.Auth.(auth.NatsToken)
				assert.True(t, ok, "Auth method should be TokenAuth")
				assert.Equal(t, "my-secret-token", tokenAuth.Token, "Wrong token")
			},
		},
	}

	// Iterate through all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Check if the broker was created on the map
			b, ok := wrapper.TelemetryAgent.Brokers[tc.brokerName]
			assert.True(t, ok, "Broker '%s' not found in map", tc.brokerName)

			// Assert common properties
			assert.Equal(t, tc.expectedType, b.Type(), "Wrong broker type")
			assert.Equal(t, tc.expectedAuthType, b.AuthMethod(), "Wrong authentication method")

			// Assert specific broker address
			if natsBroker, ok := b.(*nats.NATS); ok {
				assert.Equal(t, tc.expectedAddress, natsBroker.Config.URL, "Wrong broker address")
			}

			// Run specific assertions for this broker type
			tc.assertion(t, b)
		})
	}
}

func TestUnmarshalYAML_Fail(t *testing.T) {
	y := `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: ligmaMyNats
      type: nats
      address: "localhost:4222"
      auth:
        method: notExistMethod 
        user: test
        password: testpass
    - name: foo
      address: "localhost:9999"
      auth:
        method: natsToken
        token: my-secret-token`

	var wrapper struct {
		TelemetryAgent TelemetryAgent `yaml:"telemetryAgent"`
	}
	err := yaml.Unmarshal([]byte(y), &wrapper)
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
		assertion        func(t *testing.T, b brokers.Broker)
	}{
		{
			name:             "NATS broker with plain auth",
			brokerName:       "ligmaMyNats",
			expectedType:     "nats",
			expectedAddress:  "localhost:4222",
			expectedAuthType: "plain",
			assertion: func(t *testing.T, b brokers.Broker) {
				natsBroker, ok := b.(*nats.NATS)
				assert.True(t, ok, "Broker should be a NATS broker")

				basicAuth, ok := natsBroker.Config.Auth.(auth.BasicAuth)
				assert.True(t, ok, "Auth method should be BasicAuth")
				assert.Equal(t, "test", basicAuth.Username, "Wrong username")
				assert.Equal(t, "testpass", basicAuth.Password, "Wrong password")
			},
		},
		{
			name:             "NATS broker with token auth",
			brokerName:       "foo",
			expectedType:     "nats",
			expectedAddress:  "localhost:9999",
			expectedAuthType: "natsToken",
			assertion: func(t *testing.T, b brokers.Broker) {
				natsBroker, ok := b.(*nats.NATS)
				assert.True(t, ok, "Broker should be a NATS broker")

				tokenAuth, ok := natsBroker.Config.Auth.(auth.NatsToken)
				assert.True(t, ok, "Auth method should be TokenAuth")
				assert.Equal(t, "my-secret-token", tokenAuth.Token, "Wrong token")
			},
		},
	}

	// Iterate through all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Check if the broker was created on the map
			b, ok := wrapper.TelemetryAgent.Brokers[tc.brokerName]
			assert.True(t, ok, "Broker '%s' not found in map", tc.brokerName)

			// Assert common properties
			assert.Equal(t, tc.expectedType, b.Type(), "Wrong broker type")
			assert.Equal(t, tc.expectedAuthType, b.AuthMethod(), "Wrong authentication method")

			// Assert specific broker address
			if natsBroker, ok := b.(*nats.NATS); ok {
				assert.Equal(t, tc.expectedAddress, natsBroker.Config.URL, "Wrong broker address")
			}

			// Run specific assertions for this broker type
			tc.assertion(t, b)
		})
	}
}
