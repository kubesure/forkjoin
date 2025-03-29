package forkjoin

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// NewMultiplexer creates new basic multiplexer
func NewMultiplexer() Multiplexer {
	return Multiplexer{}
}

// AddWorker adds workers to multiplex on N worker
func (m *Multiplexer) AddWorker(w Worker) {
	if m.workers == nil {
		m.workers = make([]Worker, 0)
	}
	m.workers = append(m.workers, w)
}

// Multiplex starts N goroutines configured in []config
func (m *Multiplexer) Multiplex(ctx context.Context, x interface{}) <-chan Result {

	if len(m.workers) == 0 {
		panic("no worker added")
	}

	multiplexdResultStream := make(chan Result)

	go func() {
		defer close(multiplexdResultStream)
		var wg sync.WaitGroup
		for i, worker := range m.workers {
			wg.Add(1)
			w := worker
			in := input{id: i + 1, x: x, wg: &wg, worker: w}
			go manage(ctx, in, multiplexdResultStream)
		}
		wg.Wait()
	}()
	return multiplexdResultStream
}

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

// implements the template worker algo for all workers and dispatches work to worker
func dispatch(ctx context.Context, i input, w Worker, pulseInterval time.Duration) (<-chan Result, <-chan Heartbeat) {

	pulseStream := make(chan Heartbeat)
	result := make(<-chan Result)
	resultStream := make(chan Result)
	quitPulseStream := make(chan interface{})

	go func() {
		defer close(quitPulseStream)
		defer close(resultStream)
		result = w.Work(ctx, i.x)
		resultStream <- <-result
		quitPulseStream <- struct{}{}
	}()

	go func() {
		SendPulse(ctx, pulseStream, pulseInterval)
	}()

	return resultStream, pulseStream
}

func SendPulse(ctx context.Context, pulseStream chan Heartbeat, pulseInterval time.Duration) {
	defer close(pulseStream)
	pulse := time.NewTicker(pulseInterval)
	log := NewLogger()
	sendPulse := func() {
		select {
		case pulseStream <- Heartbeat{}:
			log.LogInfo(fmt.Sprint(pulseStream), "Sending pulse")
		default:
		}
	}
	//send pulse at a pulseInterval
	for {
		select {
		case <-ctx.Done():
			return
		case <-pulse.C:
			sendPulse()
		}
	}
}
