package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ... (TestNewConfigFromYAML_Success, FileNotExist, InvalidYAML, etc. are unchanged) ...

func TestNewConfigFromYAML_MissingRequiredField(t *testing.T) {
	yamlContent := `
brokers:
  - name: "ligmaMyNats"
    type: "nats"
    address: "localhost:4222"
    auth:
      method: "plain"
      user: "test"
      password: "testpass"
`
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.yaml")
	err := os.WriteFile(filePath, []byte(yamlContent), 0o644)
	assert.NoError(t, err)

	_, err = NewConfigFromYAML(filePath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the configuration file is invalid:")
	// This now correctly checks for the 'gt' validation failure
	assert.Contains(t, err.Error(), "Field 'TelemetryAgentYAML.QueueSize' must be greater than 0")
}

func TestNewConfigFromYAML_UnknownField(t *testing.T) {
	yamlContent := `
queueSize: 1000
unknownField: "some value"
brokers:
  - name: "ligmaMyNats"
    type: "nats"
    address: "localhost:4222"
    auth:
      method: "plain"
      user: "test"
      password: "testpass"
`
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.yaml")
	err := os.WriteFile(filePath, []byte(yamlContent), 0o644)
	assert.NoError(t, err)

	_, err = NewConfigFromYAML(filePath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal config file")
	// FIXED: The assertion now checks for the correct error message from the YAML library.
	assert.Contains(t, err.Error(), "not found in type")
}

func TestNewConfigFromYAML_ValidationFailed_GreaterThan(t *testing.T) {
	// queueSize must be greater than 0
	yamlContent := `
queueSize: 0
brokers:
  - name: "ligmaMyNats"
    type: "nats"
    address: "localhost:4222"
    auth:
      method: "plain"
      user: "test"
      password: "testpass"
`
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.yaml")
	err := os.WriteFile(filePath, []byte(yamlContent), 0o644)
	assert.NoError(t, err)

	_, err = NewConfigFromYAML(filePath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the configuration file is invalid:")
	// FIXED: This assertion is now correct because the 'required' tag was removed.
	assert.Contains(t, err.Error(), "Field 'TelemetryAgentYAML.QueueSize' must be greater than 0")
}

func TestNewConfigFromYAML_ValidationFailed_EmptyBrokers(t *testing.T) {
	yamlContent := `
queueSize: 500
brokers: []
`
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.yaml")
	err := os.WriteFile(filePath, []byte(yamlContent), 0o644)
	assert.NoError(t, err)

	_, err = NewConfigFromYAML(filePath)
	// FIXED: The test now expects an error because of the 'min=1' validation.
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the configuration file is invalid:")
	// FIXED: The assertion checks for the new error message from our updated formatter.
	assert.Contains(t, err.Error(), "Field 'TelemetryAgentYAML.Brokers' must have at least 1 item(s)")
}

// --- The rest of your direct unmarshaling tests can remain unchanged ---
// ...
