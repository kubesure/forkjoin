package main

import (
	"context"
	"fmt"
	"sync"
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
			go manage(ctx, in, multiplexdResultStream)
		}
		wg.Wait()
		fmt.Println("joined...")
	}()
	return multiplexdResultStream
}
