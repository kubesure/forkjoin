package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	f "github.com/kubesure/forkjoin"
)

//DispatchWorker dispatches to the configured URL
type DispatchWorker struct {
	Request f.HTTPRequest
}

//Work dispatches http request and stream a response back
func (hdw *DispatchWorker) Work(done <-chan interface{}, x interface{}, resultStream chan<- f.Result) {

	if len(hdw.Request.Message.Method) == 0 || len(hdw.Request.Message.URL) == 0 {
		resultStream <- f.Result{Err: &f.FJerror{Code: f.RequestError, Message: "http dispatch configuration not passed"}}
		return
	}

	var validMethod bool = false
	var methods = []f.METHOD{f.GET, f.POST, f.PUT, f.PATCH}

	for _, method := range methods {
		if hdw.Request.Message.Method == method {
			validMethod = true
			break
		}
	}

	if !validMethod {
		resultStream <- f.Result{Err: &f.FJerror{Code: f.RequestError, Message: "http method not set in configuration"}}
		return
	}

	go httpDispatch(done, hdw.Request, resultStream)

}

func httpDispatch(done <-chan interface{}, reqMsg f.HTTPRequest, resultStream chan<- f.Result) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	responseStream := make(chan f.Result)
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
			responseStream <- f.Result{Err: &f.FJerror{Code: f.ConnectionError, Message: fmt.Sprintf("error in request call %v", err)}}
		} else {
			bb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				responseStream <- f.Result{Err: &f.FJerror{Code: f.ResponseError, Message: "error reading http body"}}
			}

			hr := f.HTTPResponse{
				Message: f.HTTPMessage{
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
			responseStream <- f.Result{X: hr}
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
