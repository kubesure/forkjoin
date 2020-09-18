// The service implements a fork and join integration pattern.
// This service forks an event received on a kakfa topic on multiple goroutines.
// Each goroutine will invoke http service and merge to join point and publish the
// aggregated result to a Kafka topic.

package main

import (
	"context"
	"fmt"
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
	var pc prospectcompany
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, config := range configuration {
		c := config
		heartbeat, result := c.checker.check(ctx, &pc)

		for {
			select {
			case _, ok := <-heartbeat: // <4>
				if ok == false {
					return
				}
				fmt.Println("pulse received")
			case r, ok := <-result: // <5>
				if ok == false {
					return
				}
				if r.err != nil {
					fmt.Println(r.err.Message)
					return
				}
				fmt.Printf("results %v\n", r.pc.isMatch)
				return
			}
		}
	}
	fmt.Println("joined...")
}
