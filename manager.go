package forkjoin

import (
	"context"
	"log"
	"os"
	"time"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)
}

//goroutine waits for response from worker on work channel.
func manage(ctx context.Context, i input, multiplexdResultStream chan<- Result) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	sendResult := func(r Result) {
		r.ID = i.id
		multiplexdResultStream <- r
	}

	//start a new worker, moniters progress and returns response from worker
	//moniters for result or context done from client cancel or timeout
	go func() {
		defer cancel()
		defer i.wg.Done()
		const pulseInterval = 1 * time.Second
		workerDone := make(chan interface{})
		defer close(workerDone)
		resultStream, pulseStream := dispatch(workerDone, i, i.worker, pulseInterval)
		lastPulseT := time.Now()

		//worker moniter loop
		for {
			select {
			case <-ctx.Done():
				r := Result{ID: i.id, Err: &FJerror{Code: ConcurrencyContextError, Message: ctx.Err().Error()}}
				sendResult(r)
				return
			case <-pulseStream:
				currPulseT := time.Now()
				diff := int32(currPulseT.Sub(lastPulseT).Seconds())
				lastPulseT = currPulseT
				if diff > 2 {
					log.Printf("heartbeat inconsistent spawning new woker goroutine...\n")
					close(workerDone)
					workerDone = make(chan interface{})
					resultStream, pulseStream = dispatch(workerDone, i, i.worker, pulseInterval)
				}
			case r, _ := <-resultStream:
				sendResult(r)
				return
			}
		}
	}()
}
