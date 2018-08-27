package main

import (
	"fmt"
)

type Payment struct {
	Currency string
	Address  string
	Amount   string
	Hash     string
}

func (payment *Payment) String() string {
	return fmt.Sprintf("currency: %v\naddress: %v\namount: %v\nhash: %v", payment.Currency, payment.Address, payment.Amount, payment.Hash)
}
