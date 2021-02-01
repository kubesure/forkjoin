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
	workStream := make(chan Result)

	//moniters for result or context done from client cancel or timeout
	go func() {
		defer i.wg.Done()
		defer close(workStream)
		defer cancel()

		sendResult := func(r Result) {
			r.ID = i.id
			multiplexdResultStream <- r
		}

		for {
			select {
			case r := <-workStream:
				sendResult(r)
				return
			case <-ctx.Done():
				r := Result{ID: i.id, Err: &FJerror{Message: ctx.Err().Error()}}
				sendResult(r)
				return
			}
		}
	}()

	//start a new worker, moniters progress and returns response from worker
	go func() {
		const pulseInterval = 1 * time.Second
		workerDone := make(chan interface{})
		defer close(workerDone)
		resultStream, pulseStream := dispatch(workerDone, i, i.worker, pulseInterval)
		lastPulseT := time.Now()

		//worker moniter loop
		for {
			select {
			case <-ctx.Done():
				return
			case <-pulseStream:
				currPulseT := time.Now()
				diff := int32(currPulseT.Sub(lastPulseT).Seconds())
				lastPulseT = currPulseT
				//starts new goroutine if current pulse is delayed more than
				//2 seconds than last pulse received
				if diff > 2 {
					log.Printf("heartbeat inconsistent spawning new woker goroutine...\n")
					close(workerDone)
					workerDone = make(chan interface{})
					resultStream, pulseStream = dispatch(workerDone, i, i.worker, pulseInterval)
				}
			case r, _ := <-resultStream:
				workStream <- r
				return
			}
		}
	}()
}
