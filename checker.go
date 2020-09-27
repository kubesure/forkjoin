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
	{
		c: &creditratingchecker{},
	},
}

//Mock code. Need transform and call extenal api
func (p *policechecker) check(ctx context.Context, pc prospectcompany, work chan<- result) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(15)
	fmt.Printf("Sleeping %d seconds...\n", n)
	time.Sleep(time.Duration(n) * time.Second)
	pc.isMatch = true
	work <- result{pc: pc}
}

//Mock code. Need transform and call extenal api
func (p *centralbankchecker) check(ctx context.Context, pc prospectcompany, work chan<- result) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(15)
	fmt.Printf("Sleeping %d seconds...\n", n)
	time.Sleep(time.Duration(n) * time.Second)
	pc.isMatch = false
	work <- result{pc: pc}
}

//Mock code. Need transform and call extenal api
func (p *creditratingchecker) check(ctx context.Context, pc prospectcompany, work chan<- result) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(15)
	fmt.Printf("Sleeping %d seconds...\n", n)
	time.Sleep(time.Duration(n) * time.Second)
	pc.isMatch = false
	work <- result{pc: pc}
}
