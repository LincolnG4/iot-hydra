package workerpool

import (
	"context"
	"sync"
)

type (
	job func() error
)

type Workerpool struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Number of workers in the pool
	maxWorkers int
	// Receives the worker's jobs
	JobQueue chan job
	// Output of the workers
	ResultQueue chan error
}

// New workerpool, where size of queue of jobs need to be defined
// and the maximum number of workers.
func New(ctx context.Context, queueSize int, maxWorkers int) *Workerpool {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	return &Workerpool{
		ctx:         ctx,
		cancel:      cancel,
		wg:          wg,
		JobQueue:    make(chan job, queueSize),
		ResultQueue: make(chan error, queueSize),
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
			w.ResultQueue <- job()
		case <-w.ctx.Done():
			return
		}
	}
}

func (w *Workerpool) Stop() {
	w.cancel()
	w.wg.Wait()

	close(w.JobQueue)
	close(w.ResultQueue)
}
