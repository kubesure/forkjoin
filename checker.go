package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

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
	check(ctx context.Context, pc *prospectcompany) (<-chan interface{}, <-chan result)
}

type policechecker struct{}
type centralbankchecker struct{}
type creditratingchecker struct{}

func heartBeat() <-chan interface{} {
	heartbeat := make(chan interface{})
	defer close(heartbeat)
	pulse := time.Tick(2 * time.Second)

	sendPulse := func() {
		select {
		case heartbeat <- struct{}{}:
		default:
		}
	}

	for {
		select {
		case <-pulse:
			sendPulse()
		}
	}
}

func (p *policechecker) check(ctx context.Context, pc *prospectcompany) (<-chan interface{}, <-chan result) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	pulse := time.Tick(2 * time.Second)
	response := make(chan result)
	work := make(chan result)
	heartbeat := make(chan interface{})

	go func() {
		defer close(work)
		defer close(heartbeat)
		defer close(response)
		defer cancel()

		doWork := func(pc *prospectcompany) {
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(20)
			fmt.Printf("Sleeping %d seconds...\n", n)
			time.Sleep(time.Duration(n) * time.Second)
			work <- result{pc: prospectcompany{isMatch: true}}
		}

		go doWork(pc)

		sendPulse := func() {
			fmt.Println("pulse sent")
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}

		sendResult := func(r result) {
			select {
			case response <- r:
				return
			}
		}

		for {
			select {
			case <-pulse:
				sendPulse()
			case r := <-work:
				sendResult(r)
			case <-ctx.Done():
				r := result{err: &myerror{Message: ctx.Err().Error()}}
				sendResult(r)
			}
		}

	}()
	return heartbeat, response
}

/*func (p *centralbankchecker) check(ctx context.Context, pc *prospectcompany) *result {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(20)
	fmt.Printf("Sleeping %d seconds...\n", n)

	select {
	case <-ctx.Done():
		{
			return &result{err: &myerror{Message: ctx.Err().Error()}}
		}
	case <-time.After(time.Duration(n) * time.Second):
	}
	return &result{pc: prospectcompany{isMatch: true}}
}

func (p *creditratingchecker) check(ctx context.Context, pc *prospectcompany) *result {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(20)
	fmt.Printf("Sleeping %d seconds...\n", n)

	select {
	case <-ctx.Done():
		{
			return &result{pc: *pc, err: &myerror{Message: ctx.Err().Error()}}
		}
	case <-time.After(time.Duration(n) * time.Second):
	}
	return &result{pc: prospectcompany{isMatch: true}, err: nil}
}*/
