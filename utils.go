package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2/delivery"
)

func fetchFutures(client *delivery.Client) []string {
	var symbols []string
	futuresInfo, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return symbols
	}
	for _, f := range futuresInfo {
		if !strings.Contains(f.Symbol, "_PERP") {
			symbols = append(symbols, f.Symbol)
		}
	}
	return symbols
}

func generateRateCalculators(futureList []string) []rateCalculator {
	var calculators []rateCalculator
	for _, s := range futureList {
		spotSymbol, futureDate := strings.Split(s, "_")[0], strings.Split(s, "_")[1]
		calculators = append(calculators, rateCalculator{
			futureSymbol:   s,
			spotSymbol:     spotSymbol,
			settlementDate: time.Date(2024, time.March, 29, 0, 0, 0, 0, time.UTC),
			tradeDate:      parseFutureDate(futureDate),
		})
	}
	for _, c := range calculators {
		fmt.Println(c)
	}
	return calculators
}

func parseFutureDate(futureDate string) time.Time {
	year, _ := strconv.Atoi("20" + futureDate[:2])
	month, _ := strconv.Atoi(futureDate[2:4])
	day, _ := strconv.Atoi(futureDate[4:])
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return date
}
