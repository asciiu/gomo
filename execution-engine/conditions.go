package main

import (
	"fmt"
	"strings"

	"github.com/mattn/anko/vm"
)

// ConditionFunc is a func pointer to a condition eval func
type ConditionFunc func(price float64) bool

// PriceCondition used for simple price compares
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

// TrailingStopPoint based upon difference in points
type TrailingStopPoint struct {
	Top    float64
	Points float64
}

func (cond *TrailingStopPoint) evaluate(price float64) bool {
	if cond.Top <= 0.0 {
		cond.Top = price
		return false
	}
	// we have a new top when the price
	// is greater than the top
	if price > cond.Top {
		cond.Top = price
	}

	// trailing stop is true when the price is less than the
	// condition top minus the pts
	return (cond.Top - cond.Points) > price
}

// TrailingStopPercent based upon difference in percent
type TrailingStopPercent struct {
	Top     float64
	Percent float64
}

func (cond *TrailingStopPercent) evaluate(price float64) bool {
	if cond.Top <= 0.0 {
		cond.Top = price
		return false
	}

	// we have a new top when the price
	// is greater than the top
	if price > cond.Top {
		cond.Top = price
	}

	// trailing stop is true when the price is less than percent from top
	return (cond.Top * (1 - cond.Percent)) > price
}
