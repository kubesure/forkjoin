## Generate Go code for protocol buffer and grpc 

1. download protoc version 3.14.0 
2. download protoc plugin for protobuff grpc (https://grpc.io/docs/languages/go/quickstart/) 
3. execute 
```
protoc -I . --go_out=./http --go-grpc_out=./http ./api/httpforkjoin.proto
```

## Usage & Test

1. As go library

check [forkjoin_test.go](./forkjoin_test.go) forkjoin_test.go for complete code

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
```