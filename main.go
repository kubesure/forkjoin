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

type result struct {
	pc  prospectcompany
	err *myerror
}

type myerror struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func main() {
	fmt.Println("forked...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	checkers := len(configuration)
	multiplexdResult := make(chan result)
	heartBeats := make([]chan interface{}, checkers)

	for i := checkers; i <= checkers; i++ {
		heartBeats = append(heartBeats, make(chan interface{}))
	}
	var wg sync.WaitGroup
	wg.Add(checkers)
	for i, config := range configuration {
		cfg := config
		in := input{ctx: ctx, pc: &prospectcompany{}, wg: &wg}
		go cfg.checker.check(in, multiplexdResult, heartBeats[i])
	}

	/*go func() {
		for {
			select {
			case _, ok := <-heartbeat: // <4>
				if ok == false {
					return
				}
				fmt.Println("pulse received")
			}
		}
	}()*/

	fmt.Println("joined...")
}
