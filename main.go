// The service implements a fork and join integration pattern.
// This service forks an event received on a kakfa topic on multiple goroutines.
// Each goroutine will invoke http service and merge to join point and publish the
// aggregated result to a Kafka topic.

package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("forked...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	checkers := len(configuration)
	multiplexdResult := make(chan result)
	heartBeats := make([]<-chan heartbeat, checkers)

	var wg sync.WaitGroup
	wg.Add(checkers)
	for i, config := range configuration {
		cfg := config
		in := input{id: i, ctx: ctx, pc: &prospectcompany{}, wg: &wg}
		hb := cfg.checker.check(in, multiplexdResult)
		heartBeats = append(heartBeats, hb)
	}

	go func() {
		wg.Wait()
		close(multiplexdResult)
	}()

	for r := range multiplexdResult {
		if r.err != nil {
			fmt.Printf("Error for id: %v %v\n", r.id, r.err.Message)
		} else {
			fmt.Printf("Result for id %v is %v\n", r.id, r.pc.isMatch)
		}
	}
	fmt.Println("joined...")
}
