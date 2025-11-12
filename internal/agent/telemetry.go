package agent

import (
	"context"
	"fmt"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/brokers"
	"github.com/LincolnG4/iot-hydra/internal/config"
	"github.com/LincolnG4/iot-hydra/internal/message"
	"github.com/LincolnG4/iot-hydra/internal/workerpool"
	"github.com/rs/zerolog"
)

type TelemetryAgent struct {
	ctx        context.Context
	Cancel     context.CancelFunc
	logger     *zerolog.Logger
	Queue      chan *message.Message     // Queue telemetry messages
	Brokers    map[string]brokers.Broker // Map of brokers connected
	WorkerPool *workerpool.Workerpool
}

// NewTelemetryAgent creates and configures a new TelemetryAgent.
func NewTelemetryAgent(ctx context.Context, cfg *config.TelemetryAgentYAML, logger *zerolog.Logger) (*TelemetryAgent, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	// Startup all brokers
	brokerMap, err := setupBrokers(cfg.Brokers)
	if err != nil {
		return nil, err
	}

	wp, err := workerpool.New(ctx, cfg.QueueSize, cfg.MaxWorkers, logger)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	// The agent is assembled with the created brokers and a properly sized message queue.
	agent := &TelemetryAgent{
		Queue:      make(chan *message.Message, cfg.QueueSize),
		Brokers:    brokerMap,
		ctx:        ctx,
		Cancel:     cancel,
		WorkerPool: wp,
		logger:     logger,
	}

	return agent, nil
}

// setupBrokers iterates through all broker from yaml, startup and returns a map[string]Broker that points to each broker configured
func setupBrokers(config []config.BrokerYAML) (map[string]brokers.Broker, error) {
	// The map of brokers is created to hold the initialized brokers.
	brokerMap := make(map[string]brokers.Broker)

	// Loop through each broker configuration provided in the YAML file.
	for _, brokerCfg := range config {
		//  Create the authenticator
		authenticator, err := auth.NewAuthenticator(brokerCfg.Auth)
		if err != nil {
			return nil, fmt.Errorf("failed to create authenticator for broker '%s': %w", brokerCfg.Name, err)
		}

		// Create the broker.
		broker, err := brokers.NewBroker(brokers.Config{
			Name:    brokerCfg.Name,
			Type:    brokerCfg.Type,
			Address: brokerCfg.Address,
			Auth:    authenticator,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create broker '%s': %w", brokerCfg.Name, err)
		}

		// Check for duplicate broker names to avoid conflicts.
		name := broker.Name()
		if _, exist := brokerMap[name]; exist {
			return nil, fmt.Errorf("duplicate broker name: %s", name)
		}

		err = broker.Connect()
		if err != nil {
			return nil, fmt.Errorf("could not connect to broker '%s': %w", name, err)
		}
		brokerMap[name] = broker
	}
	return brokerMap, nil
}

func (t *TelemetryAgent) StartWorkerPool() {
	t.WorkerPool.Start()
}

// Start initiate a go routine that will receive message from the Queue. The function only if context is cancel
func (t *TelemetryAgent) Start() {
	go func() {
		for {
			select {
			case msg := <-t.Queue: // Read messsages from the Channel
				err := t.RouteMessage(msg)
				if err != nil {
					t.logger.Error().Err(err).Str("message_id", msg.ID)
				}
			case failedResult := <-t.WorkerPool.ResultQueue: // Log worker error
				t.logger.Error().Err(failedResult.Error)
			case <-t.ctx.Done(): // Context Canceled, finalizing Channel
				t.logger.Info().Msg("telemetry agent stopping")
				t.WorkerPool.Stop()
				return
			}
		}
	}()
}

// RouteMessage distribute the message over all brokers.
func (t *TelemetryAgent) RouteMessage(msg *message.Message) error {
	// Distribute message for the routers
	for _, brokerName := range msg.TargetBrokers {
		// Check if router exist
		b, exist := t.Brokers[brokerName]
		if !exist {
			t.logger.Error().Str("broker", brokerName).Str("device_id", msg.DeviceID).Str("topic", msg.Topic).Str("message_id", msg.ID).Msg("broker not configured")
			continue
		}

		// Submit messsage to the router
		err := t.WorkerPool.Submit(
			func() error {
				t.logger.Debug().Str("broker", brokerName).Str("device_id", msg.DeviceID).Str("topic", msg.Topic).Str("message_id", msg.ID).Msg("publishing telemetry")
				return b.Publish(t.ctx, msg)
			},
		)
		if err != nil {
			t.logger.Error().Err(err).Str("broker", brokerName).Str("device_id", msg.DeviceID).Str("topic", msg.Topic).Str("message_id", msg.ID).Msg("failed to enqueue publish job")
		}
	}
	return nil
}

// Submit tries to send message to the Queue, if successed it returns True,
// otherwise it will return false meaning that queue is full.
func (t *TelemetryAgent) Submit(m *message.Message) bool {
	select {
	case t.Queue <- m:
		return true
	case <-t.ctx.Done():
		return false
	}
}
