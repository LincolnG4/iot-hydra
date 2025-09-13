package nats

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/message"
	"github.com/alecthomas/assert"
	"github.com/nats-io/nats.go"
	"github.com/testcontainers/testcontainers-go"
	testnats "github.com/testcontainers/testcontainers-go/modules/nats"
)

// MockNATSConn is a mock implementation of the NATSConn interface
type MockNATSConn struct {
	SubscribeSyncFunc func(subj string) (*nats.Subscription, error)
	PublishFunc       func(subject string, data []byte) error
	CloseFunc         func()
}

func (m *MockNATSConn) SubscribeSync(subject string) (*nats.Subscription, error) {
	if m.SubscribeSyncFunc != nil {
		return m.SubscribeSyncFunc(subject)
	}
	return nil, nil
}

func (m *MockNATSConn) Publish(subject string, data []byte) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(subject, data)
	}
	return nil
}

func (m *MockNATSConn) Close() {
	if m.CloseFunc != nil {
		m.CloseFunc()
	}
}

func TestNATS_Publish(t *testing.T) {
	tests := []struct {
		name          string
		msg           message.Message
		publishFunc   func(subject string, data []byte) error
		expectedError bool
	}{
		{
			name: "Successful Publish",
			msg:  message.Message{Topic: "test/topic"},
			publishFunc: func(subject string, data []byte) error {
				return nil
			},
			expectedError: false,
		},
		{
			name: "Failed Publish",
			msg:  message.Message{Topic: "test/topic"},
			publishFunc: func(subject string, data []byte) error {
				return errors.New("publish error")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &MockNATSConn{
				PublishFunc: tt.publishFunc,
			}

			n := &NATS{
				conn:        mockConn,
				isConnected: true,
			}

			err := n.Publish(&tt.msg)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNATS_Stop(t *testing.T) {
	closeCalled := false
	mockConn := &MockNATSConn{
		CloseFunc: func() {
			closeCalled = true
		},
	}

	n := &NATS{
		conn: mockConn,
	}

	err := n.Stop()

	assert.NoError(t, err)
	assert.True(t, closeCalled, "expected Close to be called")
}

func TestNATS_Integration(t *testing.T) {
	ctx := context.Background()

	c, err := testnats.Run(ctx,
		"nats:2.9",
		testnats.WithArgument("server_name", "nats://localhost:4222"),
		testnats.WithUsername("foo"), testnats.WithPassword("bar"))
	defer func() {
		if err := testcontainers.TerminateContainer(c); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		t.Fatal(err)
	}

	cfg := Config{
		URL: "nats://localhost:4222",
		Auth: auth.BasicAuth{
			Username: "foo",
			Password: "bar",
		},
	}

	// Create a NewBroker
	nc := NewBroker(cfg)
	err = nc.Connect()
	assert.NoError(t, err, "Could not connect to the NATS container")

	go func() {
		time.Sleep(1 * time.Second)
		err := nc.Publish(&message.Message{Topic: "foo", Payload: []byte("Test")})
		assert.NoError(t, err, "Could not publish on topic on NATS")
	}()

	msg, err := nc.SubscribeAndWait("foo", 2*time.Second)
	assert.NoError(t, err, "Could not subcribe to topic on NATS")

	// check if message is correct
	assert.Equal(t, "Test", string(msg.Payload), "Expected doest match with received")
}

func BenchmarkNATS_Publish(b *testing.B) {
	// Setup: Connect to a NATS server (consider using a test-specific server)
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		b.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	broker := &NATS{
		conn:        nc,
		isConnected: true,
	}

	msg := &message.Message{
		Topic:   "benchmark-topic",
		Payload: []byte("hello world"),
	}

	b.ResetTimer() // Start the timer after setup

	// The benchmark loop
	for i := 0; i < b.N; i++ {
		if err := broker.Publish(msg); err != nil {
			b.Fatalf("Failed to publish: %v", err)
		}
	}
}

func BenchmarkNATS_PublishSubscribeLatency(b *testing.B) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		b.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	broker := &NATS{
		conn:        nc,
		isConnected: true,
	}

	topic := "latency-benchmark"
	payload := []byte("test message")

	// Set up a subscriber
	var wg sync.WaitGroup
	sub, err := nc.Subscribe(topic, func(m *nats.Msg) {
		wg.Done()
	})
	if err != nil {
		b.Fatalf("Failed to subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		err := broker.Publish(&message.Message{Topic: topic, Payload: payload})
		if err != nil {
			b.Fatalf("Publish failed: %v", err)
		}
		wg.Wait() // Wait for the subscriber to receive the message
	}
}
