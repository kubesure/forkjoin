// The service implements a fork and join integration pattern.
// This service forks an event received on a kakfa topic on multiple goroutines.
// Each goroutine will invoke http service and merge to join point and publish the
// aggregated result to a Kafka topic.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type result struct {
	pc  prospectcompany
	err *myerror
}

type myerror struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func work(ctx context.Context, pc *prospectcompany) result {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(20)
	fmt.Printf("Sleeping %d seconds...\n", n)
	select {
	case <-ctx.Done():
		{
			return result{err: &myerror{Message: ctx.Err().Error()}}
		}
	case <-time.After(time.Duration(n) * time.Second):
	}
	return result{pc: prospectcompany{isMatch: true}}
}

func main() {
	fmt.Println("fork-join")
	var wg sync.WaitGroup
	var pc prospectcompany
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(10)
	for i := 1; i <= 10; i++ {
		go func(ctx context.Context, pc *prospectcompany) {
			defer wg.Done()
			result := work(ctx, pc)
			if result.err != nil {
				fmt.Println(result.err.Message)
			} else {
				fmt.Println(result.pc.isMatch)
			}

		}(ctx, &pc)
	}
	wg.Wait()
	fmt.Println("All go routines done")
}
