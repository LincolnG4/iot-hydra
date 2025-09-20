package message

import "time"

type Message struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`

	// content of the message
	Payload []byte `json:"payload"`

	// Slice of all brokers where the message will be fowarded
	TargetBrokers []string `json:"target_brokers"`
	// Slice of from which broker the message came
	SourceBroker string `json:"source_broker"`
	Topic        string `json:"topic"`
}
