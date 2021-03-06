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
	defer cancel()
	defer i.wg.Done()
	workerDone := make(chan interface{})
	defer close(workerDone)

	sendResult := func(r Result) {
		r.ID = i.id
		multiplexdResultStream <- r
	}

	const pulseInterval = 1 * time.Second
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
				log.Println("heartbeat inconsistent spawning new woker goroutine...")
				close(workerDone)
				workerDone = make(chan interface{})
				resultStream, pulseStream = dispatch(workerDone, i, i.worker, pulseInterval)
			}
		case r, _ := <-resultStream:
			sendResult(r)
			return
		default:
		}
	}
}
