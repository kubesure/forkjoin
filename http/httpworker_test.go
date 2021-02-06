package http

import (
	"context"
	"log"
	"os"
	"testing"

	f "github.com/kubesure/forkjoin"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)
}

func TestInvalidHttpMethod(t *testing.T) {
	msg := f.HTTPMessage{URL: "https://httpbin.org/anything"}
	reqMsg := f.HTTPRequest{Message: msg}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := f.NewMultiplexer()
	m.AddWorker(&DispatchWorker{})
	resultStream := m.Multiplex(ctx, reqMsg)
	for r := range resultStream {
		if r.Err == nil {
			t.Errorf("should have given dispatch config error")
		}
	}
}

func TestEmptyHttpURL(t *testing.T) {
	msg := f.HTTPMessage{Method: f.GET}
	reqMsg := f.HTTPRequest{Message: msg}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := f.NewMultiplexer()
	m.AddWorker(&DispatchWorker{})
	resultStream := m.Multiplex(ctx, reqMsg)
	for r := range resultStream {
		if r.Err == nil {
			t.Errorf("should have given dispatch config error")
		}
	}
}

func TestHttpGETDispatch(t *testing.T) {
	msg := f.HTTPMessage{Method: f.GET, URL: "https://httpbin.org/anything"}
	msg.Add("header1", "value1")
	reqMsg := f.HTTPRequest{Message: msg}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := f.NewMultiplexer()
	m.AddWorker(&DispatchWorker{})
	resultStream := m.Multiplex(ctx, reqMsg)
	for r := range resultStream {
		if r.Err != nil {
			log.Printf("Error for id: %v %v\n", r.ID, r.Err.Message)
		} else {
			res, ok := r.X.(f.HTTPResponse)
			if !ok {
				t.Errorf("type assertion err http.Response not found")
			} else {
				if res.Message.URL != "https://httpbin.org/anything" {
					t.Errorf("URL not found in response")
				}
				if len(res.Message.Method) == 0 {
					t.Errorf("response method is empty")
				}
				if len(res.Message.Payload) == 0 {
					t.Errorf("response payload is empty")
				}
				if len(res.Message.Headers) == 0 {
					t.Errorf("response headers is empty")
				}
				//log.Println(res)
			}
		}
	}
}

func TestHttpPOSTDispatch(t *testing.T) {
	msg := f.HTTPMessage{Method: f.POST, URL: "https://httpbin.org/post"}
	msg.Add("accept", "application/json")
	reqMsg := f.HTTPRequest{Message: msg}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := f.NewMultiplexer()
	m.AddWorker(&DispatchWorker{})
	resultStream := m.Multiplex(ctx, reqMsg)
	for r := range resultStream {
		if r.Err != nil {
			log.Printf("Error for id: %v %v\n", r.ID, r.Err.Message)
		} else {
			res, ok := r.X.(f.HTTPResponse)
			if !ok {
				t.Errorf("type assertion err http.Response not found")
			} else {
				if res.Message.URL != "https://httpbin.org/post" {
					t.Errorf("URL not found in response")
				}
				if len(res.Message.Method) == 0 {
					t.Errorf("response method is empty")
				}
				if len(res.Message.Payload) == 0 {
					t.Errorf("response payload is empty")
				}
				if len(res.Message.Headers) == 0 {
					t.Errorf("response headers is empty")
				}
				//log.Println(res)
			}
		}
	}
}

func TestHttpGETDeplyedResponse(t *testing.T) {
	msg := f.HTTPMessage{Method: f.GET, URL: "http://localhost:8000/healthz"}
	reqMsg := f.HTTPRequest{Message: msg}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := f.NewMultiplexer()
	m.AddWorker(&DispatchWorker{})
	resultStream := m.Multiplex(ctx, reqMsg)
	for r := range resultStream {
		if r.Err != nil {
			log.Printf("Error for id: %v %v\n", r.ID, r.Err.Message)
		} else {
			t.Errorf("response not expected check test")
		}
	}
}

func TestHttpURLError(t *testing.T) {
	msg := f.HTTPMessage{Method: f.GET, URL: "http://unknown:8000/healthz"}
	reqMsg := f.HTTPRequest{Message: msg}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := f.NewMultiplexer()
	m.AddWorker(&DispatchWorker{})
	resultStream := m.Multiplex(ctx, reqMsg)
	for r := range resultStream {
		if r.Err == nil {
			t.Errorf("Should have failed to connect to unknown host")
		}
	}
}