package forkjoin

import (
	"context"
	"net/http"
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

	httpDispatch(done, cfg, resultStream)

}

func httpDispatch(done <-chan interface{}, cfg HTTPDispatchCfg, resultStream chan<- Result) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, string(cfg.method), cfg.url, nil)
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		resultStream <- Result{err: &FJerror{Message: "http method not set in configuration"}}
		return
	}

	resultStream <- Result{id: 1, x: res}
	return
}
