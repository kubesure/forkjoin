package main

import (
	"context"
	"fmt"
	"time"
)

//goroutine waits for reponse from checker on work channel.
func manage(ctx context.Context, i input, multiplexdResultStream chan<- result) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	workStream := make(chan result)

	go func() {
		defer i.wg.Done()
		defer close(workStream)
		defer cancel()

		sendResult := func(r result) {
			r.id = i.id
			multiplexdResultStream <- r
		}

		for {
			select {
			case r := <-workStream:
				sendResult(r)
				return
			case <-ctx.Done():
				r := result{id: i.id, err: &myerror{Message: ctx.Err().Error()}}
				sendResult(r)
				return
			}
		}
	}()

	go func() {
		const pulseInterval = 1 * time.Second
		workerDone := make(chan interface{})
		resultStream, pulseStream := i.chk.work(workerDone, i.pc, pulseInterval)
		lastPulse := time.Now()

		//worker moniter loop
		for {
			select {
			case <-ctx.Done():
				close(workerDone)
				return
			case <-pulseStream:
				currPulse := time.Now()
				diff := int32(currPulse.Sub(lastPulse).Seconds())
				lastPulse = currPulse
				if diff > 2 {
					fmt.Println("heartbeat inconsistent spawning new goroutine...")
					close(workerDone)
					workerDone = make(chan interface{})
					resultStream, pulseStream = i.chk.work(workerDone, i.pc, pulseInterval)
				}
			case r, _ := <-resultStream:
				workStream <- r
				return
			}
		}
	}()
}
