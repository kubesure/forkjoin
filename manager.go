package main

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

//goroutine waits for reponse from checker on work channel.
func manage(ctx context.Context, i input, multiplexdResultStream chan<- Result) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	workStream := make(chan Result)

	go func() {
		defer i.wg.Done()
		defer close(workStream)
		defer cancel()

		sendResult := func(r Result) {
			r.id = i.id
			multiplexdResultStream <- r
		}

		for {
			select {
			case r := <-workStream:
				sendResult(r)
				return
			case <-ctx.Done():
				r := Result{id: i.id, err: &FJerror{Message: ctx.Err().Error()}}
				sendResult(r)
				return
			}
		}
	}()

	go func() {
		const pulseInterval = 1 * time.Second
		workerDone := make(chan interface{})
		resultStream, pulseStream := dispatch(workerDone, i.x, i.worker, pulseInterval)
		lastPulseT := time.Now()

		//worker moniter loop
		for {
			select {
			case <-ctx.Done():
				close(workerDone)
				return
			case <-pulseStream:
				currPulseT := time.Now()
				diff := int32(currPulseT.Sub(lastPulseT).Seconds())
				lastPulseT = currPulseT
				if diff > 2 {
					log.Printf("heartbeat inconsistent spawning new woker goroutine...\n")
					close(workerDone)
					workerDone = make(chan interface{})
					resultStream, pulseStream = dispatch(workerDone, i.x, i.worker, pulseInterval)
				}
			case r, _ := <-resultStream:
				workStream <- r
				return
			}
		}
	}()
}
