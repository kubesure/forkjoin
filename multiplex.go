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
func multiplex(pc prospectcompany) <-chan result {
	fmt.Println("forked...")
	ctx, cancel := context.WithCancel(context.Background())
	multiplexdResultStream := make(chan result)

	go func() {
		defer close(multiplexdResultStream)
		var wg sync.WaitGroup
		for i, config := range configuration {
			wg.Add(1)
			cfg := config
			in := input{id: i, ctx: ctx, pc: pc, wg: &wg, chk: cfg.c}
			go work(in, multiplexdResultStream)
		}
		wg.Wait()
		defer cancel()
		fmt.Println("joined...")
	}()
	return multiplexdResultStream
}

//gogroutine waits for reponse from checker on work channel.
func work(i input, multiplexdResultStream chan<- result) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
	i.chk.check(ctx, i.pc, workStream)
}
