package forkjoin

import "time"

//implements the template worker algo for all workers and dispatches work to worker
func dispatch(done <-chan interface{}, i input, w Worker, pulseInterval time.Duration) (<-chan Result, <-chan heartbeat) {

	pulseStream := make(chan heartbeat)
	resultStream := make(chan Result)
	c := make(chan Result)
	quit := make(chan interface{})

	go func() {
		defer close(quit)
		defer close(c)
		defer close(resultStream)
		//dispatches work to worker
		go w.Work(done, i.x, c)
		for {
			select {
			case <-done:
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
			select {
			case pulseStream <- heartbeat{}:
			default:
			}
		}
		//send pulse at a interval or 1 second
		for {
			select {
			case <-done:
				return
			case <-quit:
				return
			case <-pulse:
				sendPulse()
			default:
			}
		}
	}()
	return resultStream, pulseStream
}
