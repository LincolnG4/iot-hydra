package agent

import (
	"strings"
	"testing"

	"github.com/LincolnG4/iot-hydra/internal/brokers/nats"
	"github.com/LincolnG4/iot-hydra/internal/config"
	"github.com/alecthomas/assert"
)

func TestNewTelemetryAgent_Success(t *testing.T) {
	cfg := &config.TelemetryAgentYAML{
		QueueSize:  10,
		MaxWorkers: 2,
		Brokers: []config.BrokerYAML{
			{
				Name:    "deezeNats",
				Type:    "nats",
				Address: "localhost:4222",
				Auth: config.AuthYAML{
					Method:   "plain",
					User:     "test",
					Password: "pwd",
				},
			},
		},
	}

	ag, err := NewTelemetryAgent(cfg)
	assert.NoError(t, err, "Should not fail to create a New Telemetrt Agent")

	b, exist := ag.Brokers["deezeNats"].(*nats.NATS)
	assert.True(t, exist, "Broker should exist")

	assert.Equal(t, cfg.Brokers[0].Name, b.Name())
	assert.Equal(t, cfg.Brokers[0].Type, b.Type())
	assert.Equal(t, cfg.Brokers[0].Type, b.Type())
}

func TestNewTelemetryAgent_Fail(t *testing.T) {
	tests := []struct {
		name          string
		input         config.TelemetryAgentYAML
		expectedError string
	}{
		{
			"Bad broker config",
			config.TelemetryAgentYAML{
				QueueSize:  10,
				MaxWorkers: 2,
				Brokers: []config.BrokerYAML{
					{
						Name:    "deezeNats",
						Type:    "wrongType",
						Address: "localhost:4222",
						Auth: config.AuthYAML{
							Method:   "plain",
							User:     "test",
							Password: "pwd",
						},
					},
				},
			},
			"failed to create broker",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTelemetryAgent(&tt.input)
			if strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("expected `%s`, got `%s`", tt.expectedError, err.Error())
			}
		})
	}
}
