package forkjoin

import (
	"sync"
)

//Result returned by checks with the result
type Result struct {
	id  int
	x   interface{}
	err *FJerror
}

//FJerror error reported by ForkJoin
type FJerror struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

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

//HTTPDispatchWorker dispatches to the configured URL
type HTTPDispatchWorker struct {
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
	Message HTTPMessage
}

//HTTPResponse URL and method to be dispatched too
type HTTPResponse struct {
	Message HTTPMessage
}

//HTTPMessage URL and method to be dispatched too
type HTTPMessage struct {
	ID      int
	URL     string
	Method  METHOD
	Payload string
	Headers map[string]string
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
	work(done <-chan interface{}, x interface{}, resultStream chan<- Result)
}
