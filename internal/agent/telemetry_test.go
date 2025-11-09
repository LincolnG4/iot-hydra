package agent

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/brokers/nats"
	"github.com/LincolnG4/iot-hydra/internal/config"
	"github.com/alecthomas/assert"
	"github.com/rs/zerolog"
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

	// Starting logs
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	ag, err := NewTelemetryAgent(context.Background(), cfg, &logger)
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

	// Starting logs
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTelemetryAgent(context.Background(), &tt.input, &logger)
			if strings.Contains(tt.expectedError, err.Error()) {
				t.Errorf("expected `%s`, got `%s`", tt.expectedError, err.Error())
			}
		})
	}
}

func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
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
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	// Starting logs
	ag, _ := NewTelemetryAgent(ctx, cfg, &logger)
	ag.WorkerPool.Start()
	time.Sleep(1 * time.Second)
	cancel()
}
