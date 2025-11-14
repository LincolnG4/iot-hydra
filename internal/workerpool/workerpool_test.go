package workerpool

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/rs/zerolog"
)

func TestNewWorkerPool(t *testing.T) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	t.Run("Creation sucessed", func(t *testing.T) {
		wp, err := NewPool(context.Background(), 1, 1, &logger)
		assert.NoError(t, err)
		assert.NotNil(t, wp, "workerpool can not be nil")
	})
	t.Run("Failed to create, empty log", func(t *testing.T) {
		wp, err := NewPool(context.Background(), 1, 0, nil)
		assert.Nil(t, wp)
		assert.Error(t, err)
	})
	t.Run("Bad max workers", func(t *testing.T) {
		wp, _ := NewPool(context.Background(), 1, 0, &logger)
		assert.Equal(t, 1, wp.maxWorkers, "when maxworkers is set less than 1, it need to be set as 1")
	})
	t.Run("Bad queueSize", func(t *testing.T) {
		wp, _ := NewPool(context.Background(), -1, 1, &logger)
		assert.Equal(t, 1, wp.maxWorkers, "when queueSize is set less than 0, it need to be set as 1")
	})
	t.Run("Start workerpool", func(t *testing.T) {
		wp, _ := NewPool(context.Background(), 1, 1, &logger)
		wp.Start()
		defer wp.Stop()
	})

	t.Run("Submit: sucessed", func(t *testing.T) {
		wp, _ := NewPool(context.Background(), 1, 1, &logger)
		wp.Start()
		defer wp.Stop()

		job := func() error { return nil }
		err := wp.Submit(job)
		assert.NoError(t, err, "should not error")
	})

	t.Run("Submit failed: closed channel", func(t *testing.T) {
		wp, _ := NewPool(context.Background(), 1, 1, &logger)
		wp.Start()

		job := func() error {
			time.Sleep(1 * time.Second)
			return nil
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			wp.Submit(job)
		}()

		wp.Stop()
		err := wp.Submit(job)
		assert.Error(t, err, "workerpool queue must be closed")
		wg.Wait() // Wait for goroutine before Stop() is called
	})

	t.Run("Submit failed: workerpool not started", func(t *testing.T) {
		wp, _ := NewPool(context.Background(), 1, 1, &logger)
		defer wp.Stop()

		job := func() error {
			return nil
		}

		err := wp.Submit(job)
		assert.Error(t, err, "workerpool queue must be closed")
	})

	t.Run("Get error from worker", func(t *testing.T) {
		wp, _ := NewPool(context.Background(), 1, 1, &logger)
		defer wp.Stop()

		wp.Start()
		job := func() error {
			return errors.New("some error")
		}

		wp.Submit(job)
		err := <-wp.ResultQueue

		assert.Error(t, err.Error, "must return error")
		assert.Contains(t, err.Error.Error(), "some error", "must return a error `some error`")
	})
}
