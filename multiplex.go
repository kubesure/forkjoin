package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func newMultiplexer() multiplexer {
	return multiplexer{}
}

func (m *multiplexer) multiplex(pc *prospectcompany) <-chan result {
	fmt.Println("forked...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	checkers := len(configuration)
	multiplexdResultStream := make(chan result)
	resultStream := make(chan result)

	var wg sync.WaitGroup
	wg.Add(checkers)
	for i, config := range configuration {
		cfg := config
		in := input{id: i, ctx: ctx, pc: pc, wg: &wg, chk: cfg.c}
		go m.work(in, multiplexdResultStream)
	}

	go func() {
		wg.Wait()
		close(multiplexdResultStream)
		fmt.Println("joined...")
	}()

	go func() {
		for r := range multiplexdResultStream {
			resultStream <- r
		}
	}()

	return resultStream
}

func (m *multiplexer) work(i input, multiplexdResultStream chan<- result) {

	ctx, cancel := context.WithTimeout(i.ctx, 10*time.Second)
	work := make(chan result)

	doWork := func(pc *prospectcompany) {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(15)
		fmt.Printf("Sleeping %d seconds...\n", n)
		time.Sleep(time.Duration(n) * time.Second)
		work <- result{id: i.id, pc: prospectcompany{isMatch: true}}
		return
	}

	go func(doWork func(pc *prospectcompany)) {
		defer i.wg.Done()
		defer close(work)
		defer cancel()

		go doWork(i.pc)

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
	}(doWork)
}
