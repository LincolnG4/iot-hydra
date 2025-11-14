package config_test

import (
	"os"
	"testing"

	"github.com/LincolnG4/iot-hydra/internal/config"
	"github.com/stretchr/testify/require"
)

func writeTmpFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "cfg-*.yaml")
	require.NoError(t, err)
	_, err = f.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}

func TestNewConfigFromYAML_Success(t *testing.T) {
	yaml := `
apiService:
  address: ":8080"
telemetryAgent:
  queueSize: 6
  maxWorkers: 1
  brokers:
    - name: ligmaNats
      type: nats
      address: "nats_main:4222"
      auth:
        method: plain
        user: test
        password: test
`
	path := writeTmpFile(t, yaml)
	_, err := config.NewConfigFromYAML(path)
	require.NoError(t, err)
}

func TestNewConfigFromYAML_FileNotFound(t *testing.T) {
	_, err := config.NewConfigFromYAML("nonexistent.yaml")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to read config file")
}

func TestNewConfigFromYAML_InvalidYAML(t *testing.T) {
	yaml := `
apiService:
  address: ":8080"
telemetryAgent:
  queueSize: not-a-number
`
	path := writeTmpFile(t, yaml)

	_, err := config.NewConfigFromYAML(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to unmarshal config file")
}

func TestNewConfigFromYAML_ValidationError(t *testing.T) {
	// Missing apiService → violates `required`
	yaml := `
telemetryAgent:
  queueSize: 6
  maxWorkers: 1
`
	path := writeTmpFile(t, yaml)

	_, err := config.NewConfigFromYAML(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "the configuration file is invalid")
}
