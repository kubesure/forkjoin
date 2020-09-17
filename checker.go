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
	check(ctx context.Context, pc *prospectcompany) *result
}

type policechecker struct{}
type centralbankchecker struct{}
type creditratingchecker struct{}

func (p *policechecker) check(ctx context.Context, pc *prospectcompany) *result {
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

func (p *centralbankchecker) check(ctx context.Context, pc *prospectcompany) *result {
	return nil
}

func (p *creditratingchecker) check(ctx context.Context, pc *prospectcompany) *result {
	return nil
}
