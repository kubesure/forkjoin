# forkjoin
The library implements a fork(fanout) and join(fanin) pattern using goroutines


## Fork Join Design

![Fork Join Diagram](forkjoin_library.png)

1. Multiplexer spawns N goroutines for N worker added through addWorker method on the Multiplexer. 
2. Multiplexer's model stream response from workers to client.
3. Each worker needs to, 
    * Implement Worker interface and return on result channel. Heartbeat is managed for the worker.
    * Exit its work on a signal from Multiplexer on the done channel  
	* Worker only need to implement the actual work.
4. The worker (goroutine) is considered unhealthy if the heartbeat is delayed by more than two seconds and is restarted. 

## Using fork join

1. [As go library](./usage_go_library.md)
2. [GRPC streaming service](./usage_turnout.md)