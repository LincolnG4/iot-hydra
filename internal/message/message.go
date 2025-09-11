package message

import "time"

type Message struct {
	ID        string         `json:"id"`
	DeviceID  string         `json:"device_id"`
	Timestamp time.Time      `json:"timestamp"`
	Type      string         `json:"type"`
	Payload   map[string]any `json:"payload"`

	Topic string `json:"topic"`
}
