package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type result struct {
	id  int
	pc  prospectcompany
	err *myerror
}

type myerror struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

type prospectcompany struct {
	id                 int
	companyName        string
	tradeLicenseNumber string
	shareHolders       []shareholder
	isMatch            bool
}

type shareholder struct {
	firstName     string
	lastName      string
	accountNumber string
	cif           string
}

//Checker interace defines the behaviour or prospect checking implementation
type Checker interface {
	check(i input, multiplexdResult chan<- result) <-chan heartbeat
}

type input struct {
	id  int
	ctx context.Context
	pc  *prospectcompany
	wg  *sync.WaitGroup
}

type heartbeat struct {
	id int
}

type policechecker struct{}
type centralbankchecker struct{}
type creditratingchecker struct{}

func (p *policechecker) check(i input, multiplexdResult chan<- result) <-chan heartbeat {

	ctx, cancel := context.WithTimeout(i.ctx, 10*time.Second)
	pulse := time.Tick(2 * time.Second)
	work := make(chan result)
	beat := make(chan heartbeat)

	doWork := func(pc *prospectcompany) {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(15)
		fmt.Printf("Sleeping %d seconds...\n", n)
		time.Sleep(time.Duration(n) * time.Second)
		work <- result{id: i.id, pc: prospectcompany{isMatch: true}}
		return
	}

	go func(doWork func(pc *prospectcompany)) {
		defer i.wg.Done()
		defer close(work)
		defer close(beat)
		defer cancel()

		go doWork(i.pc)

		sendPulse := func() {
			fmt.Println("pulse sent")
			select {
			//case beat <- heartbeat{id: i.id}:
			default:
			}
		}

		sendResult := func(r result) {
			multiplexdResult <- r
		}

		for {
			select {
			case <-pulse:
				sendPulse()
			case r := <-work:
				sendResult(r)
				return
			case <-ctx.Done():
				r := result{id: i.id, err: &myerror{Message: ctx.Err().Error()}}
				sendResult(r)
				return
			}
		}
	}(doWork)
	return beat
}
