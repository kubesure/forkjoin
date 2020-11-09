package forkjoin

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

//prospect company to be checked by workers
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

//Worker checks prospectcompany against police records
type policechecker struct{}

//Worker checks prospectcompany against central bank records
type centralbankchecker struct{}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)
}

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

//example worker
func (c *centralbankchecker) work(done <-chan interface{}, x interface{}, resultStream chan<- Result) {
	pc, ok := x.(prospectcompany)
	if !ok {
		resultStream <- Result{err: &FJerror{Message: "type assertion err prospectcompany not found"}}
		return
	}
	n := randInt(15)
	log.Printf("Sleeping %d seconds...\n", n)
	for {
		select {
		case <-done:
			return
		case <-time.After((time.Duration(n) * time.Second)):
			pc.isMatch = false
			resultStream <- Result{x: pc}
			return
		}
	}
}

//example worker
func (c *policechecker) work(done <-chan interface{}, x interface{}, resultStream chan<- Result) {
	pc, ok := x.(prospectcompany)
	if !ok {
		resultStream <- Result{err: &FJerror{Message: "type assertion err prospectcompany not found"}}
		return
	}
	n := randInt(15)
	log.Printf("Sleeping %d seconds...\n", n)
	for {
		select {
		case <-done:
			return
		case <-time.After((time.Duration(n) * time.Second)):
			pc.isMatch = false
			resultStream <- Result{x: pc}
			return
		}
	}
}

func randInt(inrange int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(inrange)
	return n
}
