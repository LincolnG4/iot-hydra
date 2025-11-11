package workerpool

import (
	"context"
	"os"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/rs/zerolog"
)

func TestNewWorkerPool(t *testing.T) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	t.Run("Creation Sucessed", func(t *testing.T) {
		wp, err := New(context.Background(), 1, 1, &logger)
		assert.NoError(t, err)
		assert.NotNil(t, wp, "workerpool can not be nil")
	})
	t.Run("Failed to create, empty log", func(t *testing.T) {
		wp, err := New(context.Background(), 1, 0, nil)
		assert.Nil(t, wp)
		assert.Error(t, err)
	})
	t.Run("Bad max workers", func(t *testing.T) {
		wp, _ := New(context.Background(), 1, 0, &logger)
		assert.Equal(t, 1, wp.maxWorkers, "when maxworkers is set less than 1, it need to be set as 1")
	})
	t.Run("Bad queueSize", func(t *testing.T) {
		wp, _ := New(context.Background(), -1, 1, &logger)
		assert.Equal(t, 1, wp.maxWorkers, "when queueSize is set less than 0, it need to be set as 1")
	})
	t.Run("Start workerpool", func(t *testing.T) {
		wp, _ := New(context.Background(), 1, 1, &logger)
		wp.Start()
		defer wp.Stop()
	})
	// t.Run("Submit: sucessed", func(t *testing.T) {
	// 	wp, _ := New(context.Background(), 1, 1, &logger)
	// 	wp.Start()
	// 	defer wp.Stop()
	//
	// 	job := func() error { return nil }
	// 	wp.Submit(job)
	// })
	// t.Run("Submit: failed", func(t *testing.T) {
	// 	wp, _ := New(context.Background(), 1, 1, &logger)
	// 	wp.Start()
	// 	defer wp.Stop()
	// })
}
