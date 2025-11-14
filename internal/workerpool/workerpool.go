package workerpool

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
)

type (
	Job func() error
)

type FailedResult struct {
	Error error
}

type Workerpool struct {
	ctx         context.Context
	cancel      context.CancelFunc
	wg          *sync.WaitGroup
	logger      *zerolog.Logger
	maxWorkers  int      // Number of workers in the pool
	JobQueue    chan Job // Receives the worker's jobs
	isClosed    bool
	mu          sync.Mutex
	ResultQueue chan FailedResult // Output of the workers
}

// NewPool creates a new workerpool, where size of queue of jobs need to be defined
// and the maximum number of workers. A error queue is set with same size of queue.
// if queueSize be zero, it will be no buffered channel, if it be < zero, it will set as 1
// It maxWorkers need to be > 0, otherwise it will be set as 1.
func NewPool(ctx context.Context, queueSize int, maxWorkers int, parentLogger *zerolog.Logger) (*Workerpool, error) {
	if parentLogger == nil {
		return nil, errors.New("logger can't be new")
	}
	logger := parentLogger.With().Str("component", "workerpool").Logger()

	if maxWorkers <= 0 {
		maxWorkers = 1
		logger.Warn().Msg("max workers can't be less or equal zero. Set to 1")
	}

	if queueSize < 0 {
		queueSize = 1
		logger.Warn().Msg("queue size can't be less than zero. Set to 1")
	}

	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	return &Workerpool{
		ctx:         ctx,
		cancel:      cancel,
		wg:          &wg,
		logger:      &logger,
		JobQueue:    make(chan Job, queueSize),
		ResultQueue: make(chan FailedResult, queueSize),
		maxWorkers:  maxWorkers,
		isClosed:    true,
	}, nil
}

// Start spawns workers into the pools.
func (w *Workerpool) Start() {
	w.logger.Info().Msg(fmt.Sprintf("starting %d workers.", w.maxWorkers))
	w.isClosed = false
	// spawn workers
	for i := range w.maxWorkers {
		w.wg.Add(1)
		go w.worker(i)
	}
}

func (w *Workerpool) worker(id int) {
	w.logger.Debug().Str("component", "worker").Int("worker ID", id).Msg("worker started")
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			w.logger.Info().Msg("workerpool context done. worker stopped")
			return

		case job, ok := <-w.JobQueue:
			if !ok {
				w.logger.Debug().Msg("channel is closed.")
				return
			}
			w.logger.Debug().Msg("starting job...")

			err := job()
			if err != nil {
				res := FailedResult{
					Error: fmt.Errorf("worker: %w", err),
				}
				w.ResultQueue <- res
			}

			w.logger.Debug().Msg("finish")
		}
	}
}

func (w *Workerpool) Stop() {
	w.logger.Info().Msg("stopping workerpool")

	// Signal that no more jobs will be submitted
	close(w.JobQueue)
	w.cancel()

	w.mu.Lock()
	w.isClosed = true
	w.mu.Unlock()

	// Wait for all workers to finish processing remaining jobs
	// Ensure any waiting workers are released
	w.wg.Wait()

	// Close results after all workers have exited
	close(w.ResultQueue)
	w.logger.Info().Msg("workerpool stopped")
}

// Submit enqueues a job for execution. It returns an error if the queue
// is closed or full.
func (w *Workerpool) Submit(j Job) error {
	w.mu.Lock()
	isClosed := w.isClosed
	w.mu.Unlock()

	// Check if the queue still open
	if isClosed {
		return errors.New("workerpool is closed")
	}

	// ok go on
	select {
	case <-w.ctx.Done():
		return w.ctx.Err()
	case w.JobQueue <- j:
		return nil
	default:
		return errors.New("workerpool queue is full")
	}
}
