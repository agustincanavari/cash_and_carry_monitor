package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
)

type underlyingFuture struct {
	lastUpdate     time.Time
	futureSymbol   string
	futurePrice    float64
	settlementDate time.Time
}

type rateCalculator struct {
	lastUpdate time.Time
	spotSymbol string
	spotPrice  float64
	futures    []underlyingFuture
	tradeDate  time.Time
}

func (f underlyingFuture) dayDifference(date time.Time) int {
	duration := f.settlementDate.Sub(date)
	return int(duration.Hours() / 24)
}

func (f underlyingFuture) APY(spotPrice float64, tradeDate time.Time) float64 {
	return 100 * (math.Pow(f.futurePrice/spotPrice, 365/float64(f.dayDifference(tradeDate))) - 1)
}

func (f underlyingFuture) APR(spotPrice float64, tradeDate time.Time) float64 {
	return 365 * (math.Pow((f.APY(spotPrice, tradeDate)/100)+1, 1.0/365.0) - 1) * 100
}

func (f underlyingFuture) yield(spotPrice float64) float64 {
	return 100 * (f.futurePrice - spotPrice) / spotPrice
}

func (r *rateCalculator) updateSpotPrice(spotClient *binance.Client) {
	spotSymbolPrice, err := spotClient.NewListPricesService().Symbol(r.spotSymbol).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	spotPrice, err := strconv.ParseFloat(spotSymbolPrice[0].Price, 64)
	if err != nil {
		fmt.Println("error converting", spotSymbolPrice[0].Price, "to float64")
		os.Exit(1)
	}
	r.spotPrice = spotPrice
	r.lastUpdate = time.Now()
}

func (r *rateCalculator) updateFuturePrices(deliveryClient *delivery.Client) {
	for i := range r.futures {
		f := &r.futures[i]
		deliverySymbolPrice, err := deliveryClient.NewListPricesService().Symbol(f.futureSymbol).Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		futurePrice, err := strconv.ParseFloat(deliverySymbolPrice[0].Price, 64)
		if err != nil {
			fmt.Println("error converting", deliverySymbolPrice[0].Price, "to float64")
			os.Exit(1)
		}
		f.futurePrice = futurePrice
		f.lastUpdate = time.Now()
	}
}

func startCalculatorUpdate(calc *rateCalculator, spotClient *binance.Client, deliveryClient *delivery.Client, updateInterval time.Duration) {
	ticker := time.NewTicker(updateInterval)
	go func() {
		//lint:ignore S1000 for ticker-based loop
		for {
			select {
			case <-ticker.C:
				calc.updateSpotPrice(spotClient)
				calc.updateFuturePrices(deliveryClient)
			}
		}
	}()
}
