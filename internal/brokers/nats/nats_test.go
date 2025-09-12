package nats

import (
	"errors"
	"testing"

	"github.com/LincolnG4/iot-hydra/internal/message"
	"github.com/alecthomas/assert"
)

// MockNATSConn is a mock implementation of the NATSConn interface
type MockNATSConn struct {
	PublishFunc func(subject string, data []byte) error
	CloseFunc   func()
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

			err := n.Publish(tt.msg)

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
