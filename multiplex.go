package forkjoin

import (
	"context"
	"sync"
)

func init() {
	//log.SetOutput(os.Stdout)
	//log.SetFlags(log.Ltime)

}

//NewMultiplexer creates new basic multiplexer
func NewMultiplexer() Multiplexer {
	return Multiplexer{}
}

//AddWorker adds workers to multiplex on N worker
func (m *Multiplexer) AddWorker(w Worker) {
	if m.workers == nil {
		m.workers = make([]Worker, 0)
	}
	m.workers = append(m.workers, w)
}

//Multiplex starts N goroutines configured in []config
func (m *Multiplexer) Multiplex(ctx context.Context, x interface{}) <-chan Result {

	log := NewLogger()

	if len(m.workers) == 0 {
		panic("no worker added")
	}
	log.LogInfo(RequestID(ctx), "", "forked..")
	multiplexdResultStream := make(chan Result)

	go func() {
		defer close(multiplexdResultStream)
		var wg sync.WaitGroup
		for i, worker := range m.workers {
			wg.Add(1)
			w := worker
			in := input{id: i, x: x, wg: &wg, worker: w}
			go manage(ctx, in, multiplexdResultStream)
		}
		wg.Wait()
		log.LogInfo(RequestID(ctx), "", "joined..")
	}()
	return multiplexdResultStream
}

func RequestID(ctx context.Context) string {
	return ctx.Value(CtxRequestID).(string)
}
