package message

import "time"

type Message struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`

	// content of the message
	Payload []byte `json:"payload"`

	// NATS, IoTHub, ...
	Type  string `json:"type"`
	Topic string `json:"topic"`
}
