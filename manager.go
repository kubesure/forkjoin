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

//implements the worker algo for all workers and dispatches work to worker
func dispatch(done <-chan interface{}, i input, w Worker, pulseInterval time.Duration) (<-chan Result, <-chan heartbeat) {

	pulseStream := make(chan heartbeat)
	resultStream := make(chan Result)
	c := make(chan Result)
	quit := make(chan interface{})

	go func() {
		defer close(c)
		defer close(resultStream)
		//dispatches work to worker
		go w.work(done, i.x, c)
		for {
			select {
			case <-done:
				return
			case r := <-c:
				resultStream <- r
				quit <- struct{}{}
				return
			default:
			}
		}
	}()

	go func() {
		defer close(pulseStream)
		pulse := time.Tick(pulseInterval)
		sendPulse := func() {
			select {
			case pulseStream <- heartbeat{}:
			default:
			}
		}

		for {
			select {
			case <-done:
				return
			case <-quit:
				return
			case <-pulse:
				sendPulse()
			default:
			}
		}
	}()
	return resultStream, pulseStream
}
