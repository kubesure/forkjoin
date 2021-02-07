package main

import (
	"context"
	"io"
	"log"
	"os"
	"testing"

	h "github.com/kubesure/forkjoin/http"
	"google.golang.org/grpc"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)
}

const (
	address = "localhost:50051"
)

//TODO: Create mock GRPC tests
func TestHTTPForkJoin(t *testing.T) {
	conn := makeGrpcConn()
	defer conn.Close()
	c := h.NewHTTPForkJoinServiceClient(conn)

	req := h.HTTPRequest{Messages: makeValidRequests()}
	stream, err := c.FanoutFanin(context.Background(), &req)
	if err != nil {
		t.Errorf("GRPC error call should have not failed with %v", err)
	}
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("%v.FanoutFanin = _, %v", c, err)
		}

		if response.Message.StatusCode != 200 {
			t.Errorf("error code is not 200 but %v", response.Message.StatusCode)
		}
	}
}

func TestInvalidHTTPURLForkJoin(t *testing.T) {
	conn := makeGrpcConn()
	defer conn.Close()
	c := h.NewHTTPForkJoinServiceClient(conn)

	req := h.HTTPRequest{Messages: makeInValidURLRequests()}
	stream, err := c.FanoutFanin(context.Background(), &req)
	if err != nil {
		t.Errorf("GRPC error call should have not failed with %v", err)
	}
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("%v.FanoutFanin = _, %v", c, err)
		}

		//TODO: collect errors and test
		if len(response.Errors) == 0 {
			t.Errorf(" there should be %v errors", 2)
		}
	}
}

func TestInvalidRequestsCfgHTTPForkJoin(t *testing.T) {
	conn := makeGrpcConn()
	defer conn.Close()
	c := h.NewHTTPForkJoinServiceClient(conn)

	req := h.HTTPRequest{Messages: makeInValidRequestsCfg()}
	stream, err := c.FanoutFanin(context.Background(), &req)
	if err != nil {
		t.Errorf("GRPC error call should have not failed with %v", err)
	}
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("%v.FanoutFanin = _, %v", c, err)
		}

		//TODO: collect errors and test
		if response.Errors[0] == nil {
			t.Error(" there should be a errors")
		}
		log.Printf("code: %v error: %v", response.Errors[0].Code, response.Errors[0].Message)
	}
}

func makeGrpcConn() *grpc.ClientConn {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func makeValidRequests() []*h.Message {

	headers := make(map[string]string)
	headers["header1"] = "value1"

	m1 := &h.Message{Method: h.Message_GET, Url: "https://httpbin.org/anything"}
	m1.Headers = headers
	m1.Payload = "{body:'body'}"

	m2 := &h.Message{Method: h.Message_POST, Url: "https://httpbin.org/anything"}
	m2.Headers = headers
	m2.Payload = "{body:'body'}"

	msgs := []*h.Message{}
	msgs = append(msgs, m1)
	msgs = append(msgs, m2)
	return msgs
}

func makeInValidURLRequests() []*h.Message {
	invalidURLMsg := &h.Message{Method: h.Message_GET, Url: "https://unknown/anything"}
	delayedMesg := &h.Message{Method: h.Message_POST, Url: "http://localhost:8000/healthz"}
	msgs := []*h.Message{}
	msgs = append(msgs, invalidURLMsg)
	msgs = append(msgs, delayedMesg)
	return msgs
}

func makeInValidRequestsCfg() []*h.Message {
	//invalidURLMsg := &h.Message{Method: h.Message_GET}
	delayedMesg := &h.Message{Method: h.Message_GET, Url: "http://localhost:8000/healthz"}
	msgs := []*h.Message{}
	//msgs = append(msgs, invalidURLMsg)
	msgs = append(msgs, delayedMesg)
	return msgs
}
