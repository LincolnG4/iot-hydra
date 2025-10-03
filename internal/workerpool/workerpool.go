package workerpool

import (
	"context"
	"fmt"
	"sync"
)

type (
	job func() error
)

type FailedResult struct {
	WorkerID int
	Error    error
}

type Workerpool struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	maxWorkers  int               // Number of workers in the pool
	JobQueue    chan job          // Receives the worker's jobs
	ResultQueue chan FailedResult // Output of the workers
}

// New workerpool, where size of queue of jobs need to be defined
// and the maximum number of workers.
func New(ctx context.Context, queueSize int, maxWorkers int) *Workerpool {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	return &Workerpool{
		ctx:         ctx,
		cancel:      cancel,
		wg:          &wg,
		JobQueue:    make(chan job, queueSize),
		ResultQueue: make(chan FailedResult, queueSize),
		maxWorkers:  maxWorkers,
	}
}

// Start spawns workers into the pools.
func (w *Workerpool) Start() {
	// spawn workers
	for i := range w.maxWorkers {
		w.wg.Add(1)
		go w.worker(i)
	}
}

func (w *Workerpool) worker(id int) {
	defer w.wg.Done()
	for {
		select {
		case job, ok := <-w.JobQueue:
			if !ok {
				return
			}
			err := job()
			if err != nil {
				res := FailedResult{
					WorkerID: id,
					Error:    fmt.Errorf("worker: %w", err),
				}
				w.ResultQueue <- res
			}
		case <-w.ctx.Done():
			return
		}
	}
}

func (w *Workerpool) Stop() {
	// Signal that no more jobs will be submitted
	close(w.JobQueue)
	// Ensure any waiting workers are released
	w.cancel()
	// Wait for all workers to finish processing remaining jobs
	w.wg.Wait()
	// Close results after all workers have exited
	close(w.ResultQueue)
}

// Submit enqueues a job for execution. It returns false if the pool is stopping.
func (w *Workerpool) Submit(j job) bool {
	select {
	case w.JobQueue <- j:
		return true
	case <-w.ctx.Done():
		return false
	}
}
