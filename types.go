package main

import (
	"context"
	"sync"
	"time"
)

//returned by checks with the result
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

//prospect company to be checked by checkers
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

//composite object to hold data for multiplexed go routines
type input struct {
	id  int
	ctx context.Context
	pc  prospectcompany
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

//Config defines configured checkers
type config struct {
	c checker
}

//Checker interace defines the behaviour or prospect checking implementation
type checker interface {
	check(done <-chan interface{}, pc prospectcompany, pulseInterval time.Duration) (<-chan result, <-chan heartbeat)
}
