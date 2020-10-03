// The service implements a fork and join integration pattern.
// This service forks multiplexes an event received on a kakfa topic to N goroutines.
// each goroutine will invoke http service multiplex will fanout and fan in and
// publish aggregated result to a Kafka topic.

package main

import (
	"context"
	"fmt"
)

var configuration = []config{
	{
		worker: &policechecker{},
	},
	/*{
		worker: &centralbankchecker{},
	},
	{
		worker: &creditratingchecker{},
	},*/
}

func main() {
	var pc prospectcompany = prospectcompany{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resultStream := multiplex(ctx, pc)
	for r := range resultStream {
		if r.err != nil {
			fmt.Printf("Error for id: %v %v\n", r.id, r.err.Message)
		} else {
			pc, ok := r.x.(prospectcompany)
			if !ok {
				fmt.Println("type assertion err prospectcompany not found")
			} else {
				fmt.Printf("Result for id %v is %v\n", r.id, pc.isMatch)
			}

		}
	}
}
