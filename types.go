package forkjoin

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

//Result returned by checks with the result
type Result struct {
	ID  int
	X   interface{}
	Err *FJError
}

//FJError error reported by ForkJoin
type FJError struct {
	Code       EventCode
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

//ErrorCode for GRPC error responses
type EventCode int32

//Error codes for GRPC error responses
const (
	InternalError EventCode = iota
	RequestError
	ResponseError
	ConnectionError
	ConcurrencyContextError
	RequestAborted
	Info
)

//composite object to hold data for multiplexed go routines
type input struct {
	id     int
	x      interface{}
	wg     *sync.WaitGroup
	worker Worker
}

type heartbeat struct {
	id int
}

//Multiplexer starts N goroutine for N dispatchers
type Multiplexer struct {
	workers []Worker
}

//METHOD http methods supported by http dispatcher
type METHOD string

//METHOD http methods supported by http dispatcher
const (
	GET   METHOD = "GET"
	POST         = "POST"
	PUT          = "PUT"
	PATCH        = "PATCH"
)

//HTTPRequest URL and method to be dispatched too
type HTTPRequest struct {
	ID      string
	Message HTTPMessage
}

//HTTPResponse URL and method to be dispatched too
type HTTPResponse struct {
	Message HTTPMessage
}

//HTTPMessage URL and method to be dispatched too
type HTTPMessage struct {
	ID         string
	URL        string
	Method     METHOD
	Payload    string
	Headers    map[string]string
	StatusCode int
}

//Add adds headers to messsage
func (hm *HTTPMessage) Add(key, value string) {
	if hm.Headers == nil {
		hm.Headers = make(map[string]string)
	}
	hm.Headers[key] = value
}

//Worker will be implement the work to be done and exit on the done channel
type Worker interface {
	//Work(done <-chan interface{}, x interface{}) <-chan Result
	Work(ctx context.Context, x interface{}) <-chan Result
}

//LogEvent stores log message
type LogEvent struct {
	id      EventCode
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

//Request ID key passsed as data in context
type ctxKey int

const (
	CtxRequestID ctxKey = iota
)

/*
func CtxRequestID(ctx context.Context, value string) string {
	return ctx.Value(ctxRequestID).(string)
}*/
