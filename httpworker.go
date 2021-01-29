package forkjoin

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

func (hdw *HTTPDispatchWorker) work(done <-chan interface{}, x interface{}, resultStream chan<- Result) {

	req, ok := x.(HTTPRequest)
	if !ok {
		resultStream <- Result{err: &FJerror{Message: "type assertion err HTTPRequest not found"}}
		return
	}

	if len(req.Message.Method) == 0 || len(req.Message.URL) == 0 {
		resultStream <- Result{err: &FJerror{Message: "http dispatch configuration not passed"}}
		return
	}

	var validMethod bool = false
	var methods = []METHOD{GET, POST, PUT, PATCH}

	for _, method := range methods {
		if req.Message.Method == method {
			validMethod = true
			break
		}
	}

	if !validMethod {
		resultStream <- Result{err: &FJerror{Message: "http method not set in configuration"}}
		return
	}

	go httpDispatch(done, req, resultStream)

}

func httpDispatch(done <-chan interface{}, reqMsg HTTPRequest, resultStream chan<- Result) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	responseStream := make(chan Result)
	defer close(responseStream)
	var isClosed bool
	var lock sync.Mutex

	go func() {
		req, _ := http.NewRequestWithContext(
			ctx, string(reqMsg.Message.Method),
			reqMsg.Message.URL,
			strings.NewReader(reqMsg.Message.Payload))

		for k, v := range reqMsg.Message.Headers {
			req.Header.Add(k, v)
		}

		client := &http.Client{}
		res, err := client.Do(req)

		lock.Lock()
		if isClosed {
			log.Println("responseStream is closed aborting goroutine")
			return
		}
		lock.Unlock()

		if err != nil {
			log.Printf("error in request call %v", err)

			responseStream <- Result{err: &FJerror{Message: fmt.Sprintf("error in request call %v", err)}}
		} else {
			bb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				responseStream <- Result{err: &FJerror{Message: "error reading http body"}}
			}

			hr := HTTPResponse{
				Message: HTTPMessage{
					ID:         reqMsg.Message.ID,
					StatusCode: res.StatusCode,
					Method:     reqMsg.Message.Method,
					URL:        reqMsg.Message.URL,
					Payload:    string(bb),
				},
			}

			for k, values := range res.Header {
				for _, value := range values {
					hr.Message.Add(k, value)
				}
			}
			defer res.Body.Close()
			responseStream <- Result{x: hr}
		}
	}()

	for {
		select {
		case <-done:
			log.Println("http request taking too long cancelling the request")
			cancel()
			lock.Lock()
			isClosed = true
			lock.Unlock()
			return
		case r := <-responseStream:
			resultStream <- r
			return
		default:
		}
	}
}
