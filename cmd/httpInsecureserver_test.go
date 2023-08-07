package main

import (
	"context"
	"io"
	"io/ioutil"
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

//run ../test/test_server.go
//run httpfjserver.go
//TODO: Create mock GRPC tests

func TestHTTPForkJoin(t *testing.T) {
	conn := makeGrpcConn()
	defer conn.Close()
	c := h.NewHTTPForkJoinServiceClient(conn)

	req := h.Request{Id: "BIN1", Messages: makeValidRequests()}
	stream, err := c.FanoutFanin(context.Background(), &req)
	if err != nil {
		t.Errorf("GRPC error call should have not failed with %v", err)
	}

	res := []*h.Response{}

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

		if response.Message.URL != "http://localhost/anything" {
			t.Errorf("URL not found in response")
		}
		if response.Message.Method == h.Message_NIL {
			t.Errorf("response method is empty")
		}
		if len(response.Message.Payload) == 0 {
			t.Errorf("response payload is empty")
		}
		if len(response.Message.Headers) == 0 {
			t.Errorf("response headers is empty")
		}
		res = append(res, response)
	}
	if len(res) != 2 {
		t.Error("there should be a 2 responses")
	}
}

func TestMutulAuth(t *testing.T) {
	conn := makeGrpcConn()
	defer conn.Close()
	c := h.NewHTTPForkJoinServiceClient(conn)

	req := h.Request{Id: "BIN1", Messages: makeMutualAuthRequestsCfg()}
	stream, err := c.FanoutFanin(context.Background(), &req)
	if err != nil {
		t.Errorf("GRPC error call should have not failed with %v", err)
	}

	res := []*h.Response{}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("%v.FanoutFanin = _, %v", c, err)
		}

		if response.Message.StatusCode != 502 {
			t.Errorf("error code is not 200 but %v", response.Message.StatusCode)
		}

		if response.Message.URL != "https://localhost:8000/mutual" {
			t.Errorf("URL not found in response")
		}
		res = append(res, response)
	}
	if len(res) != 1 {
		t.Error("there should be a 1 responses")
	}
}

func TestInvalidMutulAuth(t *testing.T) {
	conn := makeGrpcConn()
	defer conn.Close()
	c := h.NewHTTPForkJoinServiceClient(conn)

	req := h.Request{Id: "BIN1", Messages: makeAuthErrMutualReq()}
	stream, err := c.FanoutFanin(context.Background(), &req)
	if err != nil {
		t.Errorf("GRPC error call should have not failed with %v", err)
	}

	res := []*h.Response{}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("%v.FanoutFanin = _, %v", c, err)
		}
		res = append(res, response)
	}
	if len(res) != 1 {
		t.Error("there should be a 1 responses")
	}

	if res[0].Errors[0].Code != h.ErrorCode_AuthenticationError {
		t.Errorf("there should be error code: %v", h.ErrorCode_AuthenticationError)
	}

}

func TestInvalidHTTPURLForkJoin(t *testing.T) {
	conn := makeGrpcConn()
	defer conn.Close()
	c := h.NewHTTPForkJoinServiceClient(conn)

	req := h.Request{Messages: makeInValidURLRequests()}
	stream, err := c.FanoutFanin(context.Background(), &req)
	if err != nil {
		t.Errorf("GRPC error call should have not failed with %v", err)
	}

	res := []*h.Response{}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("%v.FanoutFanin = _, %v", c, err)
		}
		res = append(res, response)
	}

	if len(res) != 2 {
		t.Error("there should be a 2 errors")
	}

	if res[0].Errors[0].Code != h.ErrorCode_ConnectionError {
		t.Errorf("there should be error code: %v", h.ErrorCode_ConnectionError)
	}

	if res[1].Errors[0].Code != h.ErrorCode_ConcurrencyContextError {
		t.Errorf("there should be error code: %v", h.ErrorCode_ConnectionError)
	}
}

func TestSlowHTTPResponse(t *testing.T) {
	conn := makeGrpcConn()
	defer conn.Close()
	c := h.NewHTTPForkJoinServiceClient(conn)

	req := h.Request{Messages: makeSlowResponseRequests()}
	stream, err := c.FanoutFanin(context.Background(), &req)
	if err != nil {
		t.Errorf("GRPC error call should have not failed with %v", err)
	}

	res := []*h.Response{}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("%v.FanoutFanin = _, %v", c, err)
		}
		res = append(res, response)
	}

	if len(res) != 1 {
		t.Error("there should be a 1 errors")
	}

	log.Printf("code %v", res[0].Errors[0].Code)

	if res[0].Errors[0].Code != h.ErrorCode_ConnectionError {
		t.Errorf("there should be error code: %v", h.ErrorCode_ConnectionError)
	}
}

func TestInvalidRequestsCfgHTTPForkJoin(t *testing.T) {
	conn := makeGrpcConn()
	defer conn.Close()
	c := h.NewHTTPForkJoinServiceClient(conn)

	req := h.Request{Messages: makeInValidRequestsCfg()}
	stream, err := c.FanoutFanin(context.Background(), &req)
	if err != nil {
		t.Errorf("GRPC error call should have not failed with %v", err)
	}

	res := []*h.Response{}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("%v.FanoutFanin = _, %v", c, err)
		}
		res = append(res, response)
	}

	if len(res) == 0 {
		t.Error("there should be a 2 errors")
	}

	if res[0].Errors[0].Code != h.ErrorCode_RequestError {
		t.Errorf("there should be error code: %v", h.ErrorCode_RequestError)
	}

	if res[1].Errors[0].Code != h.ErrorCode_ConcurrencyContextError {
		t.Errorf("there should be error code: %v", h.ErrorCode_ConnectionError)
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

	m1 := &h.Message{Method: h.Message_GET, URL: "http://localhost/anything"}
	m1.Headers = headers
	m1.Payload = "{body:'body'}"
	//m1.ActiveDeadLineSeconds = 10

	m2 := &h.Message{Method: h.Message_POST, URL: "http://localhost/anything"}
	m2.Headers = headers
	m2.Payload = "{body:'body'}"
	//ßm2.ActiveDeadLineSeconds = 10

	msgs := []*h.Message{}
	msgs = append(msgs, m1)
	msgs = append(msgs, m2)

	return msgs
}

func makeInValidURLRequests() []*h.Message {
	invalidURLMsg := &h.Message{Method: h.Message_GET, URL: "https://unknown/anything"}
	delayedMesg := &h.Message{Method: h.Message_POST, URL: "http://localhost:8000/healthz"}
	msgs := []*h.Message{}
	msgs = append(msgs, invalidURLMsg)
	msgs = append(msgs, delayedMesg)
	return msgs
}

func makeSlowResponseRequests() []*h.Message {
	delayedMesg := &h.Message{Method: h.Message_POST, URL: "http://localhost:8000/healthz", ActiveDeadLineSeconds: 10}
	msgs := []*h.Message{}
	msgs = append(msgs, delayedMesg)
	return msgs
}

func makeInValidRequestsCfg() []*h.Message {
	invalidURLMsg := &h.Message{Method: h.Message_GET}
	delayedMesg := &h.Message{Method: h.Message_GET, URL: "http://localhost:8000/healthz"}
	msgs := []*h.Message{}
	msgs = append(msgs, invalidURLMsg)
	msgs = append(msgs, delayedMesg)
	return msgs
}

func makeMutualAuthRequestsCfg() []*h.Message {
	msg := &h.Message{
		Method:                h.Message_GET,
		URL:                   "https://localhost:8000/mutual",
		Authentication:        h.Message_MUTUAL,
		ActiveDeadLineSeconds: 10, Id: "001"}

	clientcrt, _ := ioutil.ReadFile("..//certs//client.crt")
	clientkey, _ := ioutil.ReadFile("..//certs//client.key")
	ca, _ := ioutil.ReadFile("..//certs//ca.crt")

	mcreds := h.Message_MutualAuthCredentials{
		ClientCertificate: string(clientcrt),
		ClientKey:         string(clientkey),
		CACertificate:     string(ca)}

	msg.MutualAuthCredentials = &mcreds
	msgs := []*h.Message{}
	msgs = append(msgs, msg)
	return msgs
}

func makeAuthErrMutualReq() []*h.Message {
	msg := &h.Message{
		Method:                h.Message_GET,
		URL:                   "https://localhost:8000/mutual",
		Authentication:        h.Message_MUTUAL,
		ActiveDeadLineSeconds: 10, Id: "001"}

	clientcrt, _ := ioutil.ReadFile("..//certs//client.crt")
	clientkey, _ := ioutil.ReadFile("..//certs//client.key")
	ca, _ := ioutil.ReadFile("..//certs//ca.crt")

	mcreds := h.Message_MutualAuthCredentials{
		ClientCertificate: string(clientcrt),
		ClientKey:         string(clientkey),
		CACertificate:     string(ca)}

	msg.MutualAuthCredentials = &mcreds
	msgs := []*h.Message{}
	msgs = append(msgs, msg)
	return msgs
}
