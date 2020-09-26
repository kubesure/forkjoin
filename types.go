package main

import (
	"context"
	"sync"
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

//Config defines configuration for routing
type config struct {
	c checker
}

//Checker interace defines the behaviour or prospect checking implementation
type checker interface {
	check(ctx context.Context, pc *prospectcompany) <-chan result
}

//Multiplexer starts n checker goroutines and waits on multiplexResultStream for results
//type Multiplexer interface {
//multiplex(i input, multiplexResultStream chan<- result) <-chan heartbeat
//}

type input struct {
	id  int
	ctx context.Context
	pc  *prospectcompany
	wg  *sync.WaitGroup
	chk checker
}

type heartbeat struct {
	id int
}

type multiplexer struct {
}

type policechecker struct{}
type centralbankchecker struct{}
type creditratingchecker struct{}
