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

func generateRateCalculators(futureList []string) map[string]*rateCalculator {
	rateCalculatorsMap := make(map[string]*rateCalculator)
	for _, s := range futureList {
		spotSymbol, futureDate := strings.Split(s, "_")[0]+"T", strings.Split(s, "_")[1]
		rc, exists := rateCalculatorsMap[spotSymbol]
		if exists {
			rc.futures = append(rc.futures, underlyingFuture{
				futureSymbol:   s,
				settlementDate: parseFutureDate(futureDate),
			})
			rateCalculatorsMap[spotSymbol] = rc
		} else {
			rateCalculatorsMap[spotSymbol] = &rateCalculator{
				spotSymbol: spotSymbol,
				tradeDate:  getTodayDate(),
				futures: []underlyingFuture{
					{
						futureSymbol:   s,
						settlementDate: parseFutureDate(futureDate),
					},
				},
			}
		}
	}
	return rateCalculatorsMap
}

func parseFutureDate(futureDate string) time.Time {
	year, _ := strconv.Atoi("20" + futureDate[:2])
	month, _ := strconv.Atoi(futureDate[2:4])
	day, _ := strconv.Atoi(futureDate[4:])
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return date
}

func getTodayDate() time.Time {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
