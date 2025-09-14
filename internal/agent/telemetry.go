package agent

import (
	"context"
	"fmt"

	"github.com/LincolnG4/iot-hydra/internal/brokers"
	"github.com/LincolnG4/iot-hydra/internal/message"
	"gopkg.in/yaml.v3"
)

type TelemetryAgent struct {
	ctx context.Context

	// Queue is responsible forward the messages to the external brokers
	Queue     chan *message.Message
	QueueSize int `yaml:"queueSize"`

	// Map of brokers connected
	Brokers map[string]brokers.Broker
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (t *TelemetryAgent) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		QueueSize int              `yaml:"queueSize"`
		Brokers   []brokers.Config `yaml:"brokers"`
	}

	if err := value.Decode(&raw); err != nil {
		return err
	}

	t.QueueSize = raw.QueueSize

	m := make(map[string]brokers.Broker)
	for _, brokerCfg := range raw.Brokers {
		broker, err := brokers.NewBroker(brokerCfg)
		if err != nil {
			return fmt.Errorf("failed to create broker: %v", err)
		}

		name := broker.Name()
		if _, exist := m[name]; exist {
			return fmt.Errorf("duplicate broker name: %s", name)
		}

		m[name] = broker
	}

	t.Brokers = m
	return nil
}
