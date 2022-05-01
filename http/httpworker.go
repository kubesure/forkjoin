package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	f "github.com/kubesure/forkjoin"
)

//DispatchWorker dispatches to the configured URL
type DispatchWorker struct {
	Request f.HTTPRequest
}

//Work dispatches http request and stream a response back
func (hdw *DispatchWorker) Work(ctx context.Context, x interface{}) <-chan f.Result {
	resultStream := make(chan f.Result)
	log := f.NewLogger()

	go func() {
		defer close(resultStream)

		if len(hdw.Request.Message.Method) == 0 || len(hdw.Request.Message.URL) == 0 {
			log.LogInvalidRequest(RequestID(ctx), hdw.Request.Message.ID, "Method or URL is empty")
			resultStream <- f.Result{Err: &f.FJError{Code: f.RequestError, Message: "Method or URL is empty"}}
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
			log.LogInvalidRequest(RequestID(ctx), hdw.Request.Message.ID, fmt.Sprintf("Method %v is invalid", hdw.Request.Message.Method))
			resultStream <- f.Result{Err: &f.FJError{Code: f.RequestError, Message: fmt.Sprintf("Method %v is invalid", hdw.Request.Message.Method)}}
			return
		}
		//httpDispatch(done, hdw.Request, resultStream)
		httpDispatch(ctx, hdw.Request, resultStream)
	}()
	return resultStream
}

func httpDispatch(ctx context.Context, reqMsg f.HTTPRequest, resultStream chan<- f.Result) {
	ctxReq, cancelReq := context.WithCancel(context.Background())
	defer cancelReq()
	responseStream := make(chan f.Result)
	log := f.NewLogger()

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
			log.LogAbortedRequest(RequestID(ctx), reqMsg.Message.ID, "Aborted")
			responseStream <- f.Result{Err: &f.FJError{Code: f.RequestAborted, Message: fmt.Sprintf("Aborted %v", err)}}
		} else if err != nil {
			log.LogRequestDispatchError(RequestID(ctx), reqMsg.Message.ID, err.Error())
			responseStream <- f.Result{Err: &f.FJError{Code: f.ConnectionError, Message: fmt.Sprintf("Error in dispatching request: %v", err)}}
		} else {
			bb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.LogResponseError(RequestID(ctx), reqMsg.Message.ID, err.Error())
				responseStream <- f.Result{Err: &f.FJError{Code: f.ResponseError, Message: fmt.Sprintf("Error reading http response: %v", err)}}
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
			responseStream <- f.Result{ID: reqMsg.ID, X: hr}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.LogAbortedRequest(RequestID(ctx), reqMsg.Message.ID, "Signal received")
			cancelReq()
			resultStream <- f.Result{Err: &f.FJError{Code: f.RequestAborted, Message: fmt.Sprintf("Request aborted took longer than expected")}}
			return
		case r := <-responseStream:
			resultStream <- r
			return
		default:
		}
	}
}
