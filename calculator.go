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
	futureSymbol   string
	futurePrice    float64
	settlementDate time.Time
}

type rateCalculator struct {
	spotSymbol string
	spotPrice  float64
	futures    []underlyingFuture
	tradeDate  time.Time
}

func (r *rateCalculator) print() {
	fmt.Printf("Spot Symbol: %s, Spot Price: %f, Trade Date: %s\n", r.spotSymbol, r.spotPrice, r.tradeDate.Format("2006-01-02"))
	fmt.Println("Futures:")
	for _, f := range r.futures {
		fmt.Printf("  Symbol: %s, Price: %f, Settlement Date: %s, APR: %f, APY: %f\n", f.futureSymbol, f.futurePrice, f.settlementDate.Format("2006-01-02"), f.APR(r.spotPrice, r.tradeDate), f.APY(r.spotPrice, r.tradeDate))
	}
	fmt.Println("")
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
