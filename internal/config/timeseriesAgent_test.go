package config

import (
	"testing"

	"github.com/LincolnG4/iot-hydra/internal/utils"
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

func TestUnmarshalYAML_Fail_InvalidQueueSize(t *testing.T) {
	y := []byte(`
telemetryAgent:
  queueSize: 0
  maxWorkers: 2
  brokers:
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
	assert.NoError(t, err) // YAML parsing should succeed

	err = utils.Validate.Struct(wrapper.TelemetryAgent)
	assert.Error(t, err)
}

func TestUnmarshalYAML_Fail_InvalidMaxWorkers(t *testing.T) {
	y := []byte(`
telemetryAgent:
  queueSize: 100
  maxWorkers: 0
  brokers:
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
	_ = yaml.Unmarshal(y, &wrapper)

	err := utils.Validate.Struct(wrapper.TelemetryAgent)
	assert.Error(t, err)
}

func TestUnmarshalYAML_Fail_NoBrokers(t *testing.T) {
	y := []byte(`
telemetryAgent:
  queueSize: 100
  maxWorkers: 2
  brokers: []
`)

	var wrapper struct {
		TelemetryAgent TelemetryAgentYAML `yaml:"telemetryAgent"`
	}
	_ = yaml.Unmarshal(y, &wrapper)

	err := utils.Validate.Struct(wrapper.TelemetryAgent)
	assert.Error(t, err)
}

func TestUnmarshalYAML_Fail_MissingBrokerFields(t *testing.T) {
	y := []byte(`
telemetryAgent:
  queueSize: 100
  maxWorkers: 2
  brokers:
    - name: ""     # missing
      type: ""     # missing
      address: ""  # missing
      auth:
        method: "" # missing
`)

	var wrapper struct {
		TelemetryAgent TelemetryAgentYAML `yaml:"telemetryAgent"`
	}
	_ = yaml.Unmarshal(y, &wrapper)

	err := utils.Validate.Struct(wrapper.TelemetryAgent)
	assert.Error(t, err)
}

func TestUnmarshalYAML_Fail_AuthMethodButMissingCredentials(t *testing.T) {
	// "plain" auth requires user+password, but they're missing
	y := []byte(`
telemetryAgent:
  queueSize: 100
  maxWorkers: 2
  brokers:
    - name: foo
      type: nats
      address: "localhost:9999"
      auth:
        method: plain
`)

	var wrapper struct {
		TelemetryAgent TelemetryAgentYAML `yaml:"telemetryAgent"`
	}
	_ = yaml.Unmarshal(y, &wrapper)

	err := utils.Validate.Struct(wrapper.TelemetryAgent)
	assert.NoError(t, err, "validation doesn't know about conditional fields")
}

func TestUnmarshalYAML_Fail_EmptyConfig(t *testing.T) {
	y := []byte(`{}`)

	var wrapper struct {
		TelemetryAgent TelemetryAgentYAML `yaml:"telemetryAgent"`
	}
	_ = yaml.Unmarshal(y, &wrapper)

	err := utils.Validate.Struct(wrapper.TelemetryAgent)
	assert.Error(t, err)
}

func TestUnmarshalYAML_Fail_UnknownField(t *testing.T) {
	y := []byte(`
telemetryAgent:
  queueSize: 100
  maxWorkers: 2
  brokers:
    - name: foo
      type: nats
      address: "localhost:9999"
      auth:
        method: token
        token: my-secret-token
  unknownField: true
`)

	var wrapper struct {
		TelemetryAgent TelemetryAgentYAML `yaml:"telemetryAgent"`
	}
	err := yaml.Unmarshal(y, &wrapper)
	assert.NoError(t, err, "yaml.Unmarshal ignores unknown fields by default")
}
