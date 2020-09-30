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
	done := make(chan interface{})

	go func() {
		defer i.wg.Done()
		defer close(workStream)
		defer cancel()

		sendResult := func(r result) {
			r.id = i.id
			multiplexdResultStream <- r
			close(done)
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
		resultStream, pulseStream := i.chk.check(done, i.pc, pulseInterval)

		go func() {
			for p := range pulseStream {
				fmt.Printf("pulse received from %v..\n", p.id)
			}
		}()

		for r := range resultStream {
			workStream <- r
		}
	}()
}
