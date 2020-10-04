package main

import (
	"log"
	"os"
	"time"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)
}

//implements the worker algo for all workers and dispatches work to worker
func dispatch(done <-chan interface{}, x interface{}, w Worker, pulseInterval time.Duration) (<-chan Result, <-chan heartbeat) {

	pulseStream := make(chan heartbeat)
	resultStream := make(chan Result)
	c := make(chan Result)
	quit := make(chan interface{})

	go func() {
		defer close(c)
		defer close(resultStream)
		//dispatches work to worker
		go w.work(done, x, c)
		for {
			select {
			case <-done:
				//log.Printf("goroutine asked to quit\n")
				return
			case r := <-c:
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
			//log.Printf("pluse sent from %v\n")
			select {
			case pulseStream <- heartbeat{}:
			default:
			}
		}

		for {
			select {
			case <-done:
				//log.Printf("pulse asked to quit\n")
				return
			case <-quit:
				//log.Printf("pulse quitting after result returned\n")
				return
			case <-pulse:
				sendPulse()
			default:
			}
		}
	}()
	return resultStream, pulseStream
}
