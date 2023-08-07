package forkjoin

import (
	"context"
	"fmt"
	"time"
)

// goroutine waits for response from worker on work channel.
func manage(ctx context.Context, i input, multiplexdResultStream chan<- Result) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(i.worker.ActiveDeadLineSeconds())*time.Second)
	defer cancel()
	defer i.wg.Done()

	log := NewLogger()

	sendResult := func(r Result) {
		multiplexdResultStream <- r
	}

	const pulseInterval = 1 * time.Second
	resultStream, pulseStream := dispatch(ctx, i, i.worker, pulseInterval)
	lastPulseT := time.Now()

	//worker moniter loop
	for {
		select {
		case hb := <-pulseStream:
			currPulseT := time.Now()
			diff := int32(currPulseT.Sub(lastPulseT).Seconds())
			lastPulseT = currPulseT
			if diff > 2 {
				log.LogInfo(fmt.Sprint(hb.id), "Heartbeat inconsistent spawning new worker goroutine")
				resultStream, pulseStream = dispatch(ctx, i, i.worker, pulseInterval)
			}
		case r := <-resultStream:
			sendResult(r)
			return
		}
	}
}

// TODO how to make worker whom this lib has no control implement send pulse?
// SendPulse to be called by worker to send heart beat back to manager
func SendPulse(ctx context.Context, pulseStream chan Heartbeat, pulseInterval time.Duration) {
	defer close(pulseStream)
	pulse := time.NewTicker(pulseInterval)

	sendPulse := func() {
		select {
		case pulseStream <- Heartbeat{}:
		default:
		}
	}
	//send pulse at a interval or 1 second
	for {
		select {
		case <-ctx.Done():
			return
		case <-pulse.C:
			sendPulse()
		}
	}
}
