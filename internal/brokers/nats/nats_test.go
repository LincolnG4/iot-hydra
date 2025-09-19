package nats

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/auth"
	"github.com/LincolnG4/iot-hydra/internal/message"
	"github.com/alecthomas/assert"
	"github.com/nats-io/nats.go"
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

	natsContainer, err := testnats.Run(ctx,
		"nats:2.9",
		testnats.WithUsername("foo"),
		testnats.WithPassword("bar"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := natsContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	mappedURL, err := natsContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "Connect with BasicAuth",
			cfg: Config{
				URL: mappedURL,
				Auth: &auth.BasicAuth{
					Username: "foo",
					Password: "bar",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create and connect broker
			broker := NewBroker(tt.cfg)
			err := broker.Connect()
			assert.NoError(t, err, "Could not connect to NATS container")
			assert.True(t, broker.isConnected, "Broker should be connected")

			// Publish in background
			go func() {
				time.Sleep(1 * time.Second)
				err := broker.Publish(&message.Message{Topic: "foo", Payload: []byte("Test")})
				assert.NoError(t, err, "Could not publish to NATS")
			}()

			// Subscribe and validate
			msg, err := broker.SubscribeAndWait("foo", 2*time.Second)
			assert.NoError(t, err, "Could not subscribe on NATS")
			assert.Equal(t, "Test", string(msg.Payload))
		})
	}
}

// Lightweight negative cases (no container needed)
func TestNATS_ConnectFailures(t *testing.T) {
	tests := []struct {
		name             string
		cfg              Config
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "Unsupported auth method",
			cfg: Config{
				URL:  "nats://localhost:4222",
				Auth: &unsupportedAuth{},
			},
			expectError:      true,
			expectedErrorMsg: "method &{} not allowed",
		},
		{
			name: "Invalid host",
			cfg: Config{
				URL:  "nats://bad-host:4222",
				Auth: &auth.BasicAuth{Username: "foo", Password: "bar"},
			},
			expectError:      true,
			expectedErrorMsg: "lookup bad-host", // or "connection refused", depending on env
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			broker := NewBroker(tt.cfg)
			err := broker.Connect()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type unsupportedAuth struct{}

func (u *unsupportedAuth) AuthMethod() string { return "unsupported" }
func (u *unsupportedAuth) Validate() error    { return nil }

func TestGetCredentials(t *testing.T) {
	tests := []struct {
		name        string
		authInput   auth.Authenticator
		expectError bool
	}{
		{
			name: "basic auth valid",
			authInput: &auth.BasicAuth{
				Username: "testuser",
				Password: "testpass",
			},
			expectError: false,
		},
		{
			name: "token auth valid",
			authInput: &auth.TokenAuth{
				Token: "sometoken",
			},
			expectError: false,
		},
		{
			name:        "unsupported auth type",
			authInput:   &unsupportedAuth{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := getCredentials(tt.authInput)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(opts) == 0 {
				t.Errorf("expected at least one nats.Option, got none")
			}

			// Apply the option to a dummy nats.Options to verify it doesn't panic
			natsOpts := nats.Options{}
			for _, o := range opts {
				if o == nil {
					t.Errorf("received nil nats.Option")
				}
				err := o(&natsOpts)
				if err != nil {
					t.Errorf("option application failed: %v", err)
				}
			}
		})
	}
}
