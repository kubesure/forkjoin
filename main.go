// The service implements a fork and join integration pattern.
// This service forks an event received on a kakfa topic on multiple goroutines.
// Each goroutine will invoke http service and merge to join point and publish the
// aggregated result to a Kafka topic.

package main

import (
	"fmt"
)

func main() {

	var pc *prospectcompany = &prospectcompany{}
	m := newMultiplexer()
	multiplexdResultSteam := m.multiplex(pc)
	for r := range multiplexdResultSteam {
		if r.err != nil {
			fmt.Printf("Error for id: %v %v\n", r.id, r.err.Message)
		} else {
			fmt.Printf("Result for id %v is %v\n", r.id, r.pc.isMatch)
		}
	}

}
