package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	f "github.com/kubesure/forkjoin"
)

//DispatchWorker dispatches to the configured URL
type DispatchWorker struct {
	Request f.HTTPRequest
}

//Work dispatches http request and stream a response back
//func (hdw *DispatchWorker) Work(done <-chan interface{}, x interface{}) <-chan f.Result {
func (hdw *DispatchWorker) Work(ctx context.Context, x interface{}) <-chan f.Result {
	resultStream := make(chan f.Result)

	go func() {
		defer close(resultStream)

		if len(hdw.Request.Message.Method) == 0 || len(hdw.Request.Message.URL) == 0 {
			resultStream <- f.Result{Err: &f.FJError{Code: f.RequestError, Message: "http dispatch configuration not passed"}}
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
			resultStream <- f.Result{Err: &f.FJError{Code: f.RequestError, Message: "http method not found"}}
			return
		}
		//httpDispatch(done, hdw.Request, resultStream)
		httpDispatch(ctx, hdw.Request, resultStream)
	}()
	return resultStream
}

//func httpDispatch(done <-chan interface{}, reqMsg f.HTTPRequest, resultStream chan<- f.Result) {
func httpDispatch(ctx context.Context, reqMsg f.HTTPRequest, resultStream chan<- f.Result) {

	ctxReq, cancelReq := context.WithCancel(context.Background())
	defer cancelReq()
	responseStream := make(chan f.Result)

	go func() {
		defer close(responseStream)
		req, _ := http.NewRequestWithContext(
			ctxReq, string(reqMsg.Message.Method),
			reqMsg.Message.URL,
			strings.NewReader(reqMsg.Message.Payload))

		for k, v := range reqMsg.Message.Headers {
			req.Header.Add(k, v)
		}

		client := &http.Client{}
		res, err := client.Do(req)

		if ctx.Err() != nil && res == nil {
			log.Printf("Context cancelled")
			responseStream <- f.Result{Err: &f.FJError{Code: f.RequestAborted, Message: fmt.Sprintf("request aborted %v", err)}}
		} else if err != nil {
			responseStream <- f.Result{Err: &f.FJError{Code: f.ConnectionError, Message: fmt.Sprintf("error in request call %v", err)}}
		} else {
			bb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				responseStream <- f.Result{Err: &f.FJError{Code: f.ResponseError, Message: "error reading http body"}}
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
		case <-ctx.Done():
			log.Println("done: http request taking too long aborting request")
			cancelReq()
			resultStream <- f.Result{Err: &f.FJError{Code: f.RequestAborted, Message: fmt.Sprintf("request aborted")}}
			return
		case r := <-responseStream:
			resultStream <- r
			return
		default:
		}
	}
}
