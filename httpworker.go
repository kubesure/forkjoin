package forkjoin

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (hdw *HTTPDispatchWorker) work(done <-chan interface{}, x interface{}, resultStream chan<- Result) {

	req, ok := x.(HTTPRequest)
	if !ok {
		resultStream <- Result{err: &FJerror{Message: "type assertion err HTTPDispatchCfg not found"}}
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

	go func() {
		//TODO: why close is out of the this goroutine is giving an write on close error
		defer close(responseStream)
		req, _ := http.NewRequestWithContext(
			ctx, string(reqMsg.Message.Method),
			reqMsg.Message.URL,
			strings.NewReader(reqMsg.Message.Payload))
		//req.Header = reqMsg.Message.Headers

		for k, v := range reqMsg.Message.Headers {
			req.Header.Add(k, v)
		}

		client := &http.Client{}
		res, err := client.Do(req)

		if err != nil {
			responseStream <- Result{err: &FJerror{Message: fmt.Sprintf("error in request call %v", err)}}
		} else {
			bb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				responseStream <- Result{err: &FJerror{Message: "error reading http body"}}
			}

			hr := HTTPResponse{
				Message: HTTPMessage{
					ID:      reqMsg.Message.ID,
					Method:  reqMsg.Message.Method,
					URL:     reqMsg.Message.URL,
					Payload: string(bb),
				},
			}

			for k, values := range res.Header {
				for _, value := range values {
					hr.Message.Add(k, value)
				}
			}
			responseStream <- Result{x: hr}
			defer res.Body.Close()
		}
	}()

	for {
		select {
		case <-done:
			resultStream <- Result{err: &FJerror{Message: "http call taking too long breaking the call"}}
			cancel()
			return
		case r := <-responseStream:
			resultStream <- r
			return
		default:
		}
	}
}
