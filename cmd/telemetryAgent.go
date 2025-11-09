package main

import (
	"context"

	ag "github.com/LincolnG4/iot-hydra/internal/agent"
)

// startTelemetryAgent inits the agent responsible to receive messages and delivery to brokers.
func (a *application) startTelemetryAgent(ctx context.Context) error {
	var err error

	// Create the agent that will handle messages
	logger := a.logger.With().Str("service", "telemetry agent").Logger()
	agent, err := ag.NewTelemetryAgent(ctx, &a.config.TelemetryAgent, &logger)
	if err != nil {
		return err
	}

	// Start workerpool
	agent.StartWorkerPool()

	// go routine read the Telemetry Messages and route to brokers.
	agent.Start()

	a.TelemetryAgent = agent
	return nil
}
