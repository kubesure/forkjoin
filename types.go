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

//HTTPDispatchCfg URL and method to be dispatched too
type HTTPDispatchCfg struct {
	url    string
	method METHOD
}

//Dispatch will be implemented by the worker to dispatch the request
//type Dispatch func(done <-chan interface{}, x interface{}, result <-chan result)

//Worker will be implement the work to be done and exit on the done channel
type Worker interface {
	work(done <-chan interface{}, x interface{}, resultStream chan<- Result)
}
