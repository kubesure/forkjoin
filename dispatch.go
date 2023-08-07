package forkjoin

import (
	"context"
	"time"
)

// implements the template worker algo for all workers and dispatches work to worker
func dispatch(ctx context.Context, i input, w Worker, pulseInterval time.Duration) (<-chan Result, <-chan Heartbeat) {

	pulseStream := make(<-chan Heartbeat)
	result := make(<-chan Result)
	resultStream := make(chan Result)
	quitPulseStream := make(chan interface{})

	go func() {
		defer close(quitPulseStream)
		defer close(resultStream)
		result, pulseStream = w.Work(ctx, i.x)
		resultStream <- <-result
		quitPulseStream <- struct{}{}
	}()

	return resultStream, pulseStream
}
