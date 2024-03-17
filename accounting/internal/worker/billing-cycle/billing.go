package billing_cycle

import (
	"context"
	"fmt"
)

type billingCycleRepository interface {
	CreateNextCycleAndClosePrevious(ctx context.Context) error
}

type transactionRepository interface {
}

type Worker struct {
	Schedule string
}

func New(schedule string) Worker {
	return Worker{Schedule: schedule}
}

func (w Worker) Do(ctx context.Context) error {
	// create new cycle
	// close old cycle
	// add record with payments for all accounts
	//
	fmt.Println("hello")
	return nil
}
