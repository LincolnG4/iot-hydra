package main

import (
	"context"

	"github.com/LincolnG4/iot-hydra/internal/agent"
	"github.com/LincolnG4/iot-hydra/internal/message"
)

func (a *application) startTelemetryAgent(ctx context.Context) error {
	var err error
	a.TelemetryAgent, err = agent.NewTelemetryAgent(&a.config.TelemetryAgent)
	if err != nil {
		return err
	}

	// start workerpool
	a.TelemetryAgent.WorkerPool.Start()

	a.ctx = ctx
	// go func to read the Telemetry Messages and route to brokers.
	go func() {
		for {
			select {
			case msg := <-a.TelemetryAgent.Queue: // Read messsages from the Channel
				err := a.RouteMessage(ctx, msg)
				if err != nil {
					a.logger.Error().Err(err).Str("message_id", msg.ID)
				}
			case failedResult := <-a.TelemetryAgent.WorkerPool.ResultQueue: // Log worker error
				a.logger.Error().Err(failedResult.Error)
			case <-ctx.Done(): // Context Canceled, finalizing Channel
				a.logger.Info().Msg("telemetry agent stopping")
				a.TelemetryAgent.WorkerPool.Stop()
				return
			}
		}
	}()
	return nil
}

// RouteMessage distribute the message over all brokers.
func (a *application) RouteMessage(ctx context.Context, msg *message.Message) error {
	// Distribute message for the routers
	for _, brokerName := range msg.TargetBrokers {
		// Check if router exist
		b, exist := a.TelemetryAgent.Brokers[brokerName]
		if !exist {
			a.logger.Error().Str("broker", brokerName).Str("device_id", msg.DeviceID).Str("topic", msg.Topic).Str("message_id", msg.ID).Msg("broker not configured")
			continue
		}

		// Submit messsage to the router
		submitted := a.TelemetryAgent.WorkerPool.Submit(func() error {
			a.logger.Debug().Str("broker", brokerName).Str("device_id", msg.DeviceID).Str("topic", msg.Topic).Str("message_id", msg.ID).Msg("publishing telemetry")
			return b.Publish(ctx, msg)
		})

		// Check if the channel is not closed
		if !submitted {
			a.logger.Error().Str("broker", brokerName).Str("device_id", msg.DeviceID).Str("topic", msg.Topic).Str("message_id", msg.ID).Msg("failed to enqueue publish job")
		}
	}
	return nil
}
