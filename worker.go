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
func (p *policechecker) work(done <-chan interface{}, pc prospectcompany, pulseInterval time.Duration) (<-chan result, <-chan heartbeat) {
	pulseStream := make(chan heartbeat)
	resultStream := make(chan result)
	c := make(chan result)
	quit := make(chan interface{})

	go func() {
		defer close(c)
		defer close(resultStream)
		go mockResult(done, pc, c)
		for {
			select {
			case <-done:
				fmt.Printf("goroutine asked to quit\n")
				return
			case r := <-c:
				fmt.Println("result received")
				resultStream <- r
				quit <- struct{}{}
				return
			default:
			}
		}
	}()

	go func() {
		defer close(pulseStream)
		pulse := time.Tick(pulseInterval)
		sendPulse := func() {
			//n := randInt(4)
			//fmt.Printf("random pulse %v\n", n)
			//time.Sleep(time.Duration(n) * time.Second)
			//fmt.Printf("pluse sent from %v\n", pc.id)
			select {
			case pulseStream <- heartbeat{}:
			default:
			}
		}

		for {
			select {
			case <-done:
				fmt.Printf("pulse asked to quit\n")
				return
			case <-quit:
				fmt.Printf("pulse quitting after result returned\n")
				return
			case <-pulse:
				sendPulse()
			default:
			}
		}
	}()
	return resultStream, pulseStream
}

func mockResult(done <-chan interface{}, pc prospectcompany, resultStream chan<- result) {
	n := randInt(15)
	fmt.Printf("Sleeping %d seconds...\n", n)
	for {
		select {
		case <-done:
			println("function asked to quit...")
			return
		case <-time.After((time.Duration(n) * time.Second)):
			pc.isMatch = false
			resultStream <- result{pc: pc}
			return
		}
	}
}

func randInt(inrange int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(inrange)
	return n
}
