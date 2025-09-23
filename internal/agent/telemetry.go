package agent

import (
	"context"
	"fmt"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/brokers"
	"github.com/LincolnG4/iot-hydra/internal/config"
	"github.com/LincolnG4/iot-hydra/internal/message"
	"github.com/LincolnG4/iot-hydra/internal/workerpool"
)

type TelemetryAgent struct {
	ctx context.Context

	// Queue telemetry messages
	Queue chan *message.Message

	// Map of brokers connected
	Brokers map[string]brokers.Broker

	WorkerPool *workerpool.Workerpool
}

// NewTelemetryAgent creates and configures a new TelemetryAgent.
func NewTelemetryAgent(cfg *config.TelemetryAgentYAML) (*TelemetryAgent, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	// The map of brokers is created to hold the initialized brokers.
	brokerMap := make(map[string]brokers.Broker)

	// loop through each broker configuration provided in the YAML file.
	for _, brokerCfg := range cfg.Brokers {
		//  create the authenticator
		authenticator, err := auth.NewAuthenticator(brokerCfg.Auth)
		if err != nil {
			return nil, fmt.Errorf("failed to create authenticator for broker '%s': %w", brokerCfg.Name, err)
		}

		// create the broker .
		broker, err := brokers.NewBroker(brokers.Config{
			Name:    brokerCfg.Name,
			Type:    brokerCfg.Type,
			Address: brokerCfg.Address,
			Auth:    authenticator,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create broker '%s': %w", brokerCfg.Name, err)
		}

		// check for duplicate broker names to avoid conflicts.
		name := broker.Name()
		if _, exist := brokerMap[name]; exist {
			return nil, fmt.Errorf("duplicate broker name: %s", name)
		}

		err = broker.Connect()
		if err != nil {
			return nil, fmt.Errorf("could not connect")
		}
		brokerMap[name] = broker
	}

	ctx := context.Background()

	// The agent is assembled with the created brokers and a properly sized message queue.
	agent := &TelemetryAgent{
		Queue:      make(chan *message.Message, cfg.QueueSize),
		Brokers:    brokerMap,
		ctx:        context.Background(),
		WorkerPool: workerpool.New(ctx, cfg.QueueSize, 3),
	}

	return agent, nil
}

// func (t *TelemetryAgent) RouteToBrokers(msg *message.Message) error {
// 	var errors []error
//
// 	for _, brokerName := range msg.TargetBrokers {
// 		b, exist := t.Brokers[brokerName]
// 		if !exist {
// 			err := fmt.Errorf("broker '%s' is not configured", brokerName)
// 			errors = append(errors, err)
// 			continue
// 		}
//
// 		if err := b.Publish(msg); err != nil {
// 			err = fmt.Errorf("failed to publish message to broker '%s': %w", brokerName, err)
// 			errors = append(errors, err)
// 		}
// 	}
//
// 	if len(errors) > 0 {
// 		return fmt.Errorf("failed to route message to %d broker(s): %v", len(errors), errors)
// 	}
//
// 	return nil
// }
