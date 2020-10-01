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
	c := make(chan result)

	go func() {
		defer close(resultStream)
		go mockResult(pc, resultStream)
		for {
			select {
			case <-done:
				//fmt.Printf("goroutine quitting\n")
				return
			case r := <-c:
				fmt.Println("result received")
				resultStream <- r
			default:
			}
		}
	}()

	go func() {
		defer close(pulseStream)
		pulse := time.Tick(pulseInterval)
		sendPulse := func() {
			n := randInt(4)
			fmt.Printf("random pulse %v\n", n)
			time.Sleep(time.Duration(n) * time.Second)
			fmt.Printf("pluse sent from %v\n", pc.id)
			select {
			case pulseStream <- heartbeat{}:
			default:
			}
		}

		for {
			select {
			case <-done:
				//fmt.Printf("pulse quitting\n")
				return
			case <-pulse:
				sendPulse()
			default:
			}
		}
	}()
	return resultStream, pulseStream
}

func mockResult(pc prospectcompany, resultStream chan<- result) {
	n := randInt(15)
	fmt.Printf("Sleeping %d seconds...\n", n)
	time.Sleep(time.Duration(n) * time.Second)
	pc.isMatch = false
	resultStream <- result{pc: pc}
}

func randInt(inrange int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(inrange)
	return n
}
