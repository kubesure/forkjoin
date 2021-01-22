package forkjoin

import (
	"context"
	"log"
	"os"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)
}

func TestHttpGETDispatch(t *testing.T) {
	msg := HTTPMessage{Method: GET, URL: "https://httpbin.org/anything"}
	msg.Add("header1", "value1")
	reqMsg := HTTPRequest{Message: msg}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := NewMultiplexer()
	m.AddWorker(&HTTPDispatchWorker{})
	resultStream := m.Multiplex(ctx, reqMsg)
	for r := range resultStream {
		if r.err != nil {
			log.Printf("Error for id: %v %v\n", r.id, r.err.Message)
		} else {
			res, ok := r.x.(HTTPResponse)
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
	msg := HTTPMessage{Method: POST, URL: "https://httpbin.org/post"}
	msg.Add("accept", "application/json")
	reqMsg := HTTPRequest{Message: msg}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := NewMultiplexer()
	m.AddWorker(&HTTPDispatchWorker{})
	resultStream := m.Multiplex(ctx, reqMsg)
	for r := range resultStream {
		if r.err != nil {
			log.Printf("Error for id: %v %v\n", r.id, r.err.Message)
		} else {
			res, ok := r.x.(HTTPResponse)
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
				log.Println(res)
			}
		}
	}
}