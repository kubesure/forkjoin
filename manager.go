package forkjoin

import (
	"context"
	"fmt"
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
	//TODO: rework and test defer and done sequence
	defer i.wg.Done()
	workerDone := make(chan interface{})
	defer close(workerDone)

	log := NewLogger()

	sendResult := func(r Result) {
		r.ID = fmt.Sprint(i.id)
		multiplexdResultStream <- r
	}

	const pulseInterval = 1 * time.Second
	resultStream, pulseStream := dispatch(ctx, i, i.worker, pulseInterval)
	lastPulseT := time.Now()

	//worker moniter loop
	for {
		select {
		case <-pulseStream:
			currPulseT := time.Now()
			diff := int32(currPulseT.Sub(lastPulseT).Seconds())
			lastPulseT = currPulseT
			if diff > 2 {
				log.LogInfo(fmt.Sprint(i.id), "Heartbeat inconsistent spawning new woker goroutine...")
				close(workerDone)
				workerDone = make(chan interface{})
				resultStream, pulseStream = dispatch(ctx, i, i.worker, pulseInterval)
			}
		case r, _ := <-resultStream:
			sendResult(r)
			return
		default:
		}
	}
}
