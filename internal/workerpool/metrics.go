package workerpool

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const name = "workerpool"

var (
	meter       = otel.Meter(name)
	workerCnt   metric.Int64UpDownCounter
	jobQueueCnt metric.Int64UpDownCounter
)

func init() {
	var err error

	workerCnt, err = meter.Int64UpDownCounter("workerpool.active_workers",
		metric.WithDescription("Number of workers running on the workerpool"),
		metric.WithUnit("{workers}"),
	)
	if err != nil {
		// TODO: NOT PANIC
		panic(err)
	}

	jobQueueCnt, err = meter.Int64UpDownCounter("workerpool.jobs_in_queue",
		metric.WithDescription("Number of jobs in the queue of the workerpool"),
		metric.WithUnit("{jobs}"),
	)
	if err != nil {
		// TODO: NOT PANIC
		panic(err)
	}
}
