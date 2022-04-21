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
func (hdw *DispatchWorker) Work(done <-chan interface{}, x interface{}) <-chan f.Result {
	resultStream := make(chan f.Result)

	go func() {
		ctx := context.Background()
		ctx, cancelReq := context.WithCancel(ctx)
		defer cancelReq()
		defer close(resultStream)
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
		go httpDispatch(ctx, hdw.Request, resultStream)

		for {
			select {
			case <-done:
				log.Println("http request taking too long cancelling the request")
				cancelReq()
				return
			default:
			}
		}
	}()
	return resultStream
}

func httpDispatch(ctx context.Context, reqMsg f.HTTPRequest, resultStream chan<- f.Result) {
	req, _ := http.NewRequestWithContext(
		ctx, string(reqMsg.Message.Method),
		reqMsg.Message.URL,
		strings.NewReader(reqMsg.Message.Payload))

	for k, v := range reqMsg.Message.Headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if context.Canceled != nil {
		log.Println(fmt.Sprintf("error in request call %v", err))
		log.Println("responseStream is closed aborting goroutine")
		return
	}

	if err != nil {
		resultStream <- f.Result{Err: &f.FJerror{Code: f.ConnectionError, Message: fmt.Sprintf("error in request call %v", err)}}
	} else {
		bb, err := ioutil.ReadAll(res.Body)
		if err != nil {
			resultStream <- f.Result{Err: &f.FJerror{Code: f.ResponseError, Message: "error reading http body"}}
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
		resultStream <- f.Result{X: hr}
	}
}
