package main

import (
	"fmt"
	"strings"

	"github.com/mattn/anko/vm"
)

// ConditionFunc is a func pointer to a condition eval func
type Expression func(price float64) (bool, string)

// PriceCondition used for simple price compares
type PriceCondition struct {
	Statement string
	Env       *vm.Env
}

func (cond *PriceCondition) evaluate(price float64) (bool, string) {
	p := fmt.Sprintf("%.8f", price)
	// replace all occurences of price in the statement
	c := strings.Replace(cond.Statement, "price", p, -1)

	result, err := cond.Env.Execute(c)
	if err != nil {
		return false, fmt.Sprintf("could not execute: %s", err.Error())
	}

	if result == true {
		return true, c
	}

	return false, "evaluated as false"
}

// TrailingStopPoint based upon difference in points
type TrailingStopPoint struct {
	Top    float64
	Points float64
}

func (cond *TrailingStopPoint) evaluate(price float64) (bool, string) {
	if cond.Top <= 0.0 {
		cond.Top = price
		return false, fmt.Sprintf("new top: %.8f", price)
	}
	// we have a new top when the price
	// is greater than the top
	if price > cond.Top {
		cond.Top = price
	}

	// trailing stop is true when the price is less than the
	// condition top minus the pts
	if (cond.Top - cond.Points) > price {
		return true, fmt.Sprintf("{condition: TrailingStopPoint, top:%.8f, points:%.8f, price:%.8f}",
			cond.Top, cond.Points, price)
	}
	return false, "evaluated as false"
}

// TrailingStopPercent based upon difference in percent
type TrailingStopPercent struct {
	Top     float64
	Percent float64
}

func (cond *TrailingStopPercent) evaluate(price float64) (bool, string) {
	if cond.Top <= 0.0 {
		cond.Top = price
		return false, fmt.Sprintf("new top: %.8f", price)
	}

	// we have a new top when the price
	// is greater than the top
	if price > cond.Top {
		cond.Top = price
	}

	// trailing stop is true when the price is less than percent from top
	if (cond.Top * (1 - cond.Percent)) > price {
		return true, fmt.Sprintf("{condition: TrailingStopPercent, top:%.8f, percent:%.8f, price:%.8f}",
			cond.Top, cond.Percent, price)
	}

	return false, "evaluated as false"
}

// This trigger will execute on the next price
type Immediate struct {
}

func (cond *Immediate) evaluate(price float64) (bool, string) {
	return true, fmt.Sprintf("%.8f", price)
}
