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
	cfg := new(HTTPDispatchCfg)
	cfg.method = GET
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := NewMultiplexer()
	m.AddWorker(&HTTPDispatchWorker{})
	resultStream := m.Multiplex(ctx, cfg)
	for r := range resultStream {
		if r.err != nil {
			log.Printf("Error for id: %v %v\n", r.id, r.err.Message)
		} else {
		}
	}
}
