package forkjoin

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (hdw *HTTPDispatchWorker) work(done <-chan interface{}, x interface{}, resultStream chan<- Result) {

	cfg, ok := x.(HTTPDispatchCfg)
	if !ok {
		resultStream <- Result{err: &FJerror{Message: "type assertion err HTTPDispatchCfg not found"}}
		return
	}

	if len(cfg.method) == 0 || len(cfg.url) == 0 {
		resultStream <- Result{err: &FJerror{Message: "http dispatch configuration not passed"}}
		return
	}

	var validMethod bool = false
	var methods = []METHOD{GET, POST, PUT, PATCH}

	for _, method := range methods {
		if cfg.method == method {
			validMethod = true
			break
		}
	}

	if !validMethod {
		resultStream <- Result{err: &FJerror{Message: "http method not set in configuration"}}
		return
	}

	go httpDispatch(done, cfg, resultStream)

}

func httpDispatch(done <-chan interface{}, cfg HTTPDispatchCfg, resultStream chan<- Result) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	responseStream := make(chan Result)

	go func(rstream chan<- Result) {

		req, _ := http.NewRequestWithContext(
			ctx, string(cfg.method),
			cfg.url,
			strings.NewReader(cfg.payload))

		for k, v := range cfg.headers {
			req.Header.Add(k, v)
		}

		client := &http.Client{}
		res, err := client.Do(req)

		if err != nil {
			rstream <- Result{err: &FJerror{Message: fmt.Sprintf("error in request call %v", err)}}
		} else {
			bb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				rstream <- Result{err: &FJerror{Message: "error reading http body"}}
			}
			rstream <- Result{x: string(bb)}
			res.Body.Close()
		}
	}(responseStream)

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
