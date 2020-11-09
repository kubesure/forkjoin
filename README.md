# forkjoin
The library implements a fork(fanout) and join(fanin) pattern using goroutines

## Design

1. Multiplexer spawns N goroutines for N worker added through addWorker method on the Multiplexer 
2. Multiplexer's model is request/response, it return only one response form the worker
3. Each worker needs to 
    * Implement Worker interface and return on result channel. Heartbeat is managed for the worker.
    * Exit its work on a signal from Manager on the done channel  
	* Worker only need to implement the actual work. 
4. The worker (goroutine) is considered unhealthy if the heartbeat is delayed by more than two seconds and   	is restarted 
   
## TODO

1. Bindings for Kafka and Java 
2. HTTP worker for simple HTTP dispatches 

## Usage & Test

refer to forkjoin_test.go for complete code

```
func TestChecker(t *testing.T) {
	//client can cancel entire processing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var pc prospectcompany = prospectcompany{}
	m := NewMultiplexer()
	m.AddWorker(&centralbankchecker{})
	m.AddWorker(&policechecker{})
	resultStream := m.Multiplex(ctx, pc)
	for r := range resultStream {
		if r.err != nil {
			log.Printf("Error for id: %v %v\n", r.id, r.err.Message)
		} else {
			pc, ok := r.x.(prospectcompany)
			if !ok {
				log.Println("type assertion err prospectcompany not found")
			} else {
				log.Printf("Result for id %v is %v\n", r.id, pc.isMatch)
			}
		}
	}
}

func (c *centralbankchecker) work(done <-chan interface{}, x interface{}, resultStream chan<- Result) {
	pc, ok := x.(prospectcompany)
	if !ok {
		resultStream <- Result{err: &FJerror{Message: "type assertion err prospectcompany not found"}}
		return
	}
	n := randInt(15)
	log.Printf("Sleeping %d seconds...\n", n)
	for {
		select {
		case <-done:
			return
		case <-time.After((time.Duration(n) * time.Second)):
			pc.isMatch = false
			resultStream <- Result{x: pc}
			return
		}
	}
}
```
