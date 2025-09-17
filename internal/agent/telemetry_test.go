package agent

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/brokers"
	"github.com/LincolnG4/iot-hydra/internal/brokers/nats"
	"github.com/LincolnG4/iot-hydra/internal/message"
	"github.com/rs/zerolog/log"
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

// TODO: Create fail tests
func TestUnmarshalYAML_Fail(t *testing.T) {
	testCases := []struct {
		name        string
		yamlContent string
		expectError string
	}{
		{
			name: "Invalid YAML syntax",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
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
        token: my-secret-token
      invalid_indentation`,
			expectError: "yaml: line",
		},
		{
			name: "Missing required broker name",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - type: nats
      address: "localhost:4222"
      auth:
        method: plain
        user: test
        password: testpass`,
			expectError: "name",
		},
		{
			name: "Missing required broker type",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      address: "localhost:4222"
      auth:
        method: plain
        user: test
        password: testpass`,
			expectError: "type",
		},
		{
			name: "Missing required broker address",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: nats
      auth:
        method: plain
        user: test
        password: testpass`,
			expectError: "address",
		},
		{
			name: "Unsupported broker type",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: kafka
      address: "localhost:9092"
      auth:
        method: plain
        user: test
        password: testpass`,
			expectError: "unsupported broker type",
		},
		{
			name: "Invalid auth method",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"
      auth:
        method: oauth2
        user: test
        password: testpass`,
			expectError: "unsupported auth method",
		},
		{
			name: "Missing auth credentials for plain method",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"
      auth:
        method: plain`,
			expectError: "username",
		},
		{
			name: "Missing password for plain auth",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"
      auth:
        method: plain
        user: test`,
			expectError: "password",
		},
		{
			name: "Missing token for natsToken auth",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"
      auth:
        method: natsToken`,
			expectError: "token",
		},
		{
			name: "Invalid queueSize type",
			yamlContent: `telemetryAgent:
  queueSize: "not-a-number"
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"
      auth:
        method: plain
        user: test
        password: testpass`,
			expectError: "cannot unmarshal",
		},
		{
			name: "Negative queueSize",
			yamlContent: `telemetryAgent:
  queueSize: -1000
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"
      auth:
        method: plain
        user: test
        password: testpass`,
			expectError: "queue size must be positive",
		},
		{
			name: "Zero queueSize",
			yamlContent: `telemetryAgent:
  queueSize: 0
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"
      auth:
        method: plain
        user: test
        password: testpass`,
			expectError: "queue size must be positive",
		},
		{
			name: "Duplicate broker names",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"
      auth:
        method: plain
        user: test
        password: testpass
    - name: test
      type: nats
      address: "localhost:9999"
      auth:
        method: natsToken
        token: my-secret-token`,
			expectError: "duplicate broker name",
		},
		{
			name: "Empty brokers array",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers: []`,
			expectError: "at least one broker must be configured",
		},
		{
			name: "Missing brokers section",
			yamlContent: `telemetryAgent:
  queueSize: 1000`,
			expectError: "brokers configuration is required",
		},
		{
			name: "Invalid address format",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: nats
      address: "invalid-address"
      auth:
        method: plain
        user: test
        password: testpass`,
			expectError: "invalid address format",
		},
		{
			name: "Empty broker name",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: ""
      type: nats
      address: "localhost:4222"
      auth:
        method: plain
        user: test
        password: testpass`,
			expectError: "broker name cannot be empty",
		},
		{
			name: "Empty auth section",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"
      auth: {}`,
			expectError: "auth method is required",
		},
		{
			name: "Missing auth section",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"`,
			expectError: "auth configuration is required",
		},
		{
			name: "Extra fields in plain auth",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"
      auth:
        method: plain
        user: test
        password: testpass
        token: should-not-be-here`,
			expectError: "unexpected field 'token' for plain auth",
		},
		{
			name: "Extra fields in token auth",
			yamlContent: `telemetryAgent:
  queueSize: 1000
  brokers:
    - name: test
      type: nats
      address: "localhost:4222"
      auth:
        method: natsToken
        token: my-token
        user: should-not-be-here`,
			expectError: "unexpected field 'user' for token auth",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var wrapper struct {
				TelemetryAgent TelemetryAgent `yaml:"telemetryAgent"`
			}

			err := yaml.Unmarshal([]byte(tc.yamlContent), &wrapper)

			// Assert that an error occurred
			assert.Error(t, err, "Expected an error for test case: %s", tc.name)

			// Assert that the error message contains the expected substring
			if err != nil {
				assert.Contains(t, strings.ToLower(err.Error()), strings.ToLower(tc.expectError),
					"Error message should contain '%s', got: %s", tc.expectError, err.Error())
			}
		})
	}
}

func TestSendMessage_Success(t *testing.T) {
	cfg := brokers.Config{
		Name:    "nats",
		Type:    "nats",
		Address: "localhost:4222",
		Auth: auth.BasicAuth{
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

	ag.Start()
	ctx, cancel := context.WithCancel(context.Background())
	go func(context.CancelFunc) {
		time.Sleep(10 * time.Second)
		cancel()
	}(cancel)

	for {
		select {
		case msg := <-ag.Queue:
			if err := ag.RouteToBrokers(msg); err != nil {
				log.Error().Err(err)
			}
		case <-ctx.Done():
			log.Error().Msg("terminated")
		}
	}
}
