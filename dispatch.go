package forkjoin

import (
	"context"
	"time"
)

//implements the template worker algo for all workers and dispatches work to worker
func dispatch(ctx context.Context, i input, w Worker, pulseInterval time.Duration) (<-chan Result, <-chan heartbeat) {

	pulseStream := make(chan heartbeat)
	resultStream := make(chan Result)
	quitPulstStream := make(chan interface{})

	go func() {
		defer close(quitPulstStream)
		defer close(resultStream)
		result := w.Work(ctx, i.x)
		resultStream <- <-result
		quitPulstStream <- struct{}{}
	}()

	go func() {
		defer close(pulseStream)
		pulse := time.Tick(pulseInterval)

		sendPulse := func() {
			select {
			case pulseStream <- heartbeat{id: i.id}:
			default:
			}
		}
		//send pulse at a interval or 1 second
		for {
			select {
			case <-quitPulstStream:
				return
			case <-pulse:
				sendPulse()
			}
		}
	}()
	return resultStream, pulseStream
}
