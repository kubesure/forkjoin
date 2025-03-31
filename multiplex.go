package forkjoin

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// NewMultiplexer creates new basic multiplexer
// NewMultiplexer creates and returns a new instance of the Multiplexer.
func NewMultiplexer() Multiplexer {
	return Multiplexer{}
}

// AddWorker adds a new Worker to the Multiplexer. If the workers slice is nil,
// it initializes the slice before appending the new Worker.
//
// Parameters:
//   - w: The Worker instance to be added to the Multiplexer.
func (m *Multiplexer) AddWorker(w Worker) {
	if m.workers == nil {
		m.workers = make([]Worker, 0)
	}
	m.workers = append(m.workers, w)
}

// Multiplex is a method of the Multiplexer struct that distributes work among
// multiple workers and collects their results into a single channel. It takes
// a context and an input value, and returns a read-only channel of Result.
//
// Parameters:
// - ctx: A context.Context used to manage the lifecycle of the goroutines.
// - x: An input value of any type to be processed by the workers.
//
// Returns:
//   - A read-only channel (<-chan Result) that streams the results produced by
//     the workers.
//
// Behavior:
// - If no workers are added to the Multiplexer, the function will panic.
// - The method launches a goroutine to manage the multiplexing process.
// - Each worker is assigned a unique input and runs concurrently.
// - The results from all workers are sent to the multiplexed result stream.
// - The channel is closed once all workers have completed their tasks.
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

// manage is responsible for managing the lifecycle of a worker goroutine and
// handling its results and heartbeat signals. It ensures that the worker
// operates within a specified timeout and restarts the worker if heartbeat
// signals become inconsistent.
//
// Parameters:
//   - ctx: The context used to manage the worker's lifecycle and timeout.
//   - i: An input struct containing data to be processed by the worker.
//   - multiplexdResultStream: A channel to send the worker's results
//
// Behavior:
//   - Creates a context with a timeout based on the worker's ActiveDeadLineSeconds.
//   - Monitors heartbeat signals from the worker via the pulseStream.
//   - Restarts the worker if heartbeat signals are inconsistent (e.g., delayed).
//   - Sends the worker's results to the multiplexdResultStream channel.
//   - Cleans up resources and signals completion via the input's WaitGroup.
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

// dispatch starts two goroutines to handle worker execution and heartbeat signaling.
// It takes the following parameters:
// - ctx: The context to manage cancellation and deadlines.
// - i: An input struct containing data to be processed by the worker.
// - w: A Worker interface that performs the actual work.
// - pulseInterval: The interval at which heartbeat signals are sent.
//
// It returns two channels:
// - A Result channel that streams the results of the worker's execution.
// - A Heartbeat channel that emits periodic heartbeat signals.
//
// The function ensures proper cleanup of resources by closing channels when the
// worker finishes execution or the context is canceled.
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

// SendPulse sends periodic heartbeat signals through the provided pulseStream channel
// at the specified pulseInterval. The function listens for context cancellation to
// gracefully terminate the operation.
//
// Parameters:
//   - ctx: The context used to manage the lifecycle of the function. When the context
//     is canceled, the function stops sending pulses and exits.
//   - pulseStream: A channel of type Heartbeat through which heartbeat signals are sent.
//   - pulseInterval: The duration between consecutive heartbeat signals.
//
// Behavior:
//   - Sends a Heartbeat{} struct through the pulseStream channel at regular intervals
//     defined by pulseInterval.
//   - Closes the pulseStream channel when the context is canceled or the function exits.
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
