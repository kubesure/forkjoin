package main

import (
	"fmt"
	"math/rand"
	"time"
)

var configuration = []config{
	{
		c: &policechecker{},
	},
	/*{
		c: &centralbankchecker{},
	},
	{
		c: &creditratingchecker{},
	},*/
}

// worker checker
func (p *policechecker) check(done <-chan interface{}, pc prospectcompany, pulseInterval time.Duration) (<-chan result, <-chan heartbeat) {
	pulseStream := make(chan heartbeat)
	resultStream := make(chan result)

	go func() {
		defer close(resultStream)
		go mockResult(pc, resultStream)
		for {
			select {
			case <-done:
				fmt.Println("Quite received....")
				return
			case <-resultStream:
			default:
			}
		}
	}()

	go func() {
		defer close(pulseStream)
		pulse := time.Tick(pulseInterval)
		sendPulse := func() {
			fmt.Printf("Pluse sent from %v\n", pc.id)
			select {
			case pulseStream <- heartbeat{}:
			default:
			}
		}

		for {
			select {
			case <-done:
				fmt.Println("quiting pulse")
				return
			case <-pulse:
				sendPulse()
			}
		}
	}()
	return resultStream, pulseStream
}

func mockResult(pc prospectcompany, resultStream chan<- result) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(15)
	fmt.Printf("Sleeping %d seconds...\n", n)
	time.Sleep(time.Duration(n) * time.Second)
	pc.isMatch = false
	resultStream <- result{pc: pc}
}
