# forkjoin
The service implements a fork and join integration pattern using goroutines.

# design of each goroutine
1. Should timeout using the context
2. Result and errors should be returned by channels
3. Respond with heat beat to notify WIP and indicate livelock
4. Participate as ward and stewards be monitered and restarted in case unhealthy   

