# forkjoin
The library implements a fork and join integration pattern using goroutines.

## Design

1. Spawns N goroutines for  N worker added through addWorker method on the Multiplexer 
2. Multiplexer is a request response, it return only one response for the worker
3. Multiplexer's manager manages the heart beat, worker only need to implement the actual work
4. Each worker needs to 
    * Implement the Worker interface and return on result channel
    * Exit its work on a signal from Manager on the done channel  


## Usage & Test

```
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var pc prospectcompany = prospectcompany{}
	m := NewMultiplexer()
	m.addWorker(&centralbankchecker{})
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
