package main

import (
	"fmt"
	"math"
	"time"
)

type rateCalculator struct {
	spotPrice      float64
	futurePrice    float64
	settlementDate time.Time
	tradeDate      time.Time
	quantity       float64
}

func (r rateCalculator) print() {
	fmt.Printf("%+v", r)
}

func (r rateCalculator) dayDifference() int {
	duration := r.settlementDate.Sub(r.tradeDate)
	return int(duration.Hours() / 24)
}

func (r rateCalculator) underlyingAmount() float64 {
	return r.quantity / r.spotPrice
}

func (r rateCalculator) earnings() float64 {
	return r.underlyingAmount() * (r.futurePrice - r.spotPrice)
}

func (r rateCalculator) APY() float64 {
	return 100 * (math.Pow(r.futurePrice/r.spotPrice, 365/float64(r.dayDifference())) - 1)
}

func (r rateCalculator) APR() float64 {
	return 365 * (math.Pow((r.APY()/100)+1, 1.0/365.0) - 1) * 100
}
