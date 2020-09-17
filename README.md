# forkjoin
The service implements a fork and join integration pattern using goroutines.

# design of each goroutine
1. Should Timeout using the context
2. Errors should be returned by channels
3. Respond with heat beat to notify WIP and indicate livelock
4. Should behave a ward to mointoring stewards to restart it in case of it becoming unhealthy   

