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
	reqMsg := HTTPRequest{}
	msg := HTTPMessage{Method: GET, URL: "http://localhost/anything"}
	msg.Add("header1", "value1")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := NewMultiplexer()
	m.AddWorker(&HTTPDispatchWorker{})
	resultStream := m.Multiplex(ctx, reqMsg)
	for r := range resultStream {
		if r.err != nil {
			log.Printf("Error for id: %v %v\n", r.id, r.err.Message)
		} else {
			res, ok := r.x.(string)
			if !ok {
				log.Println("type assertion err http.Response not found")
			} else {
				//bb, err := ioutil.ReadAll(res.Body)
				//if err != nil {
				//log.Println("error", err)
				//}
				log.Println(res)
				//res.Body.Close()
			}
		}
	}
}
