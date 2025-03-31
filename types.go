package forkjoin

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

// Result represents the result returned by a worker after processing.
type Result struct {
	ID  string      // Unique identifier for the result
	X   interface{} // The actual result data
	Err *FJError    // Error details, if any, encountered during processing
}

// FJError represents an error reported by the ForkJoin library.
type FJError struct {
	Code       EventCode              // Error code categorizing the type of error
	Inner      error                  // The underlying error, if any
	Message    string                 // Human-readable error message
	StackTrace string                 // Stack trace for debugging purposes
	Misc       map[string]interface{} // Additional metadata related to the error
}

// EventCode represents the type of event or error in the ForkJoin library.
type EventCode int32

// Predefined event codes for categorizing errors and events.
const (
	InternalError           EventCode = iota // Internal library error
	RequestError                             // Error related to the request
	ResponseError                            // Error related to the response
	ConnectionError                          // Error related to connection issues
	ConcurrencyContextError                  // Error related to concurrency or context
	RequestAborted                           // Request was aborted
	AuthenticationError                      // Authentication failure
	RequestInfo                              // Informational event related to a request
	Info                                     // General informational event
	HeartBeatInfo                            // Heartbeat-related informational event
)

// input is a composite object used to hold data for multiplexed goroutines.
type input struct {
	id     int             // Unique identifier for the input
	x      interface{}     // The actual input data
	wg     *sync.WaitGroup // WaitGroup to synchronize goroutines
	worker Worker          // The worker responsible for processing this input
}

// Heartbeat represents a heartbeat signal sent by a worker to indicate it is alive.
type Heartbeat struct {
	id int // Unique identifier for the heartbeat
}

// Multiplexer manages multiple workers and coordinates their execution.
type Multiplexer struct {
	workers []Worker // List of workers added to the multiplexer
}

// Worker defines the interface that all workers must implement.
type Worker interface {
	// Work processes the input and returns a channel for results.
	Work(ctx context.Context, x interface{}) <-chan Result

	// ActiveDeadLineSeconds specifies the maximum time (in seconds) a worker can remain active.
	ActiveDeadLineSeconds() uint32
}

// LogEvent represents a structured log message.
type LogEvent struct {
	id      EventCode // The event code associated with the log
	message string    // The log message
}

// StandardLogger wraps the logrus.Logger to enforce specific log message formats.
type StandardLogger struct {
	*logrus.Logger // Embedded logrus.Logger
}
