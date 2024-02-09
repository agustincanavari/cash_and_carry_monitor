package main

import (
	"fmt"
	"time"
)

func main() {

	calculator := rateCalculator{
		spotPrice:      47170.08,
		futurePrice:    47832,
		settlementDate: time.Date(2024, time.March, 29, 0, 0, 0, 0, time.UTC),
		tradeDate:      time.Date(2024, time.February, 9, 0, 0, 0, 0, time.UTC),
		quantity:       1000,
	}
	fmt.Println("day difference:", calculator.dayDifference())
	fmt.Println("BTC amout:", calculator.underlyingAmount())
	fmt.Println("Earnings:", calculator.earnings())
	fmt.Println("APY: ", calculator.APY(), "%")
	fmt.Println("APR: ", calculator.APR(), "%")

}
