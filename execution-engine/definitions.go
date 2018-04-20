package main

import (
	"fmt"
	"strings"

	"github.com/mattn/anko/vm"
)

type ConditionFunc func(price float64) bool

type Condition interface {
	evaluate(price float64) bool
}

type PriceCondition struct {
	Statement string
	Env       *vm.Env
}

func (cond *PriceCondition) evaluate(price float64) bool {
	p := fmt.Sprintf("%.18f", price)
	// replace all occurences of price in the statement
	c := strings.Replace(cond.Statement, "price", p, -1)

	result, err := cond.Env.Execute(c)
	if err != nil {
		return false
	}

	return result == true
}

type TrailingStopPoint struct {
	Top    float64
	Points float64
}

func (cond *TrailingStopPoint) evaluate(price float64) bool {
	if cond.Top <= 0.0 {
		cond.Top = price
		fmt.Printf("setting initial top: %.8f\n", cond.Top)
		return false
	} else {
		fmt.Printf("top: %.8f price: %.8f\n", cond.Top, price)
	}
	// we have a new top when the price
	// is greater than the top
	if price > cond.Top {
		fmt.Printf("setting a new top: %.8f\n", price)
		cond.Top = price
	}

	// trailing stop is true when the price is less than the
	// condition top minus the pts
	return (cond.Top - cond.Points) > price
}

type TrailingStopPercent struct {
	Top     float64
	Percent float64
}

func (cond *TrailingStopPercent) evaluate(price float64) bool {
	// we have a new top when the price
	// is greater than the top
	if price > cond.Top {
		cond.Top = price
	}

	// trailing stop is true when the price is less than the
	// condition top minus the pts
	return (cond.Top - (cond.Top * cond.Percent)) > price
}
