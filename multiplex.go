package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func newMultiplexer() multiplexer {
	return multiplexer{}
}

func multiplex(pc *prospectcompany) <-chan result {
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

func work(i input, multiplexdResultStream chan<- result) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	work := make(chan result)

	go func() {
		defer i.wg.Done()
		defer close(work)
		defer cancel()

		go func() {
			checkStream := i.chk.check(ctx, i.pc)
			/*for {
				select {
				case <-ctx.Done():
					return
				case w := <-checkStream:
					work <- w
				}
			}*/
			for r := range checkStream {
				work <- r
			}
		}()

		sendResult := func(r result) {
			multiplexdResultStream <- r
		}

		for {
			select {
			case r := <-work:
				sendResult(r)
				return
			case <-ctx.Done():
				r := result{id: i.id, err: &myerror{Message: ctx.Err().Error()}}
				sendResult(r)
				return
			}
		}
	}()
}
