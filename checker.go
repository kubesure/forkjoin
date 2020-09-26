package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

var configuration = []config{
	{
		c: &policechecker{},
	},
	{
		c: &centralbankchecker{},
	},
	/*{
		checker: &creditratingchecker{},
	},*/
}

func (p *policechecker) check(ctx context.Context, pc *prospectcompany) <-chan result {
	work := make(chan result)
	go func() {
		defer close(work)
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(15)
		fmt.Printf("Sleeping %d seconds...\n", n)
		time.Sleep(time.Duration(n) * time.Second)
		work <- result{pc: prospectcompany{isMatch: true}}
	}()
	return work
}

func (p *centralbankchecker) check(ctx context.Context, pc *prospectcompany) <-chan result {
	work := make(chan result)
	go func() {
		defer close(work)
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(15)
		fmt.Printf("Sleeping %d seconds...\n", n)
		time.Sleep(time.Duration(n) * time.Second)
		work <- result{pc: prospectcompany{isMatch: false}}
	}()
	return work
}
