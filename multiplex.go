package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

//creates new basic multiplexer
func newMultiplexer() multiplexer {
	return multiplexer{}
}

//starts N goroutines configured in []config
func multiplex(ctx context.Context, pc prospectcompany) <-chan result {
	fmt.Println("forked...")
	multiplexdResultStream := make(chan result)

	go func() {
		defer close(multiplexdResultStream)
		var wg sync.WaitGroup
		for i, config := range configuration {
			wg.Add(1)
			cfg := config
			in := input{id: i, pc: pc, wg: &wg, chk: cfg.c}
			go work(ctx, in, multiplexdResultStream)
		}
		wg.Wait()
		fmt.Println("joined...")
	}()
	return multiplexdResultStream
}

//gogroutine waits for reponse from checker on work channel.
func work(ctx context.Context, i input, multiplexdResultStream chan<- result) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	workStream := make(chan result)
	wardDone := make(chan interface{})

	go func() {
		defer i.wg.Done()
		defer close(workStream)
		defer cancel()

		sendResult := func(r result) {
			r.id = i.id
			multiplexdResultStream <- r
			//close(wardDone)
		}

		for {
			select {
			case r := <-workStream:
				sendResult(r)
				return
			case <-ctx.Done():
				r := result{id: i.id, err: &myerror{Message: ctx.Err().Error()}}
				close(wardDone)
				sendResult(r)
				return
			}
		}
	}()

	go func() {
		const pulseInterval = 1 * time.Second
		resultStream, pulseStream := i.chk.check(wardDone, i.pc, pulseInterval)

		go func() {
			lastPulse := time.Now()
			for {
				select {
				case _, ok := <-pulseStream:
					if ok == false {
						return
					}
					fmt.Printf("pulse received from %v..\n", i.id)
					currPulse := time.Now()
					diff := int32(currPulse.Sub(lastPulse).Seconds())
					lastPulse = currPulse
					if diff > 2 {
						fmt.Println("heart beat inconsistent spawn new goroutine")
						close(wardDone)
						wardDone = make(chan interface{})
						resultStream, pulseStream = i.chk.check(wardDone, i.pc, pulseInterval)

					}
				case r, ok := <-resultStream:
					if ok == false {
						return
					}
					workStream <- r
					close(wardDone)
					return
				}
			}
		}()

	}()
}
