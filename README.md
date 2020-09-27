# forkjoin
The service implements a fork and join integration pattern using goroutines.

# design of each goroutine
1. Each goroutine should timeout using the context - done 
2. Result and errors should be returned by channels - done 
3. Respond with heat beat to notify WIP and indicate livelock
4. Heal and restart unhealthy goroutine    

