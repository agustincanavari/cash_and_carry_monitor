package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	spotClient := binance.NewClient(apiKey, secretKey)
	deliveryClient := binance.NewDeliveryClient(apiKey, secretKey)

	// Get the current price of BTCUSDT
	spotSymbolPrice, err := spotClient.NewListPricesService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	var spotPrice float64
	for _, p := range spotSymbolPrice {
		var err error
		spotPrice, err = strconv.ParseFloat(p.Price, 64)
		if err != nil {
			fmt.Println("error converting", p.Price, "to float64")
			os.Exit(1)
		}
		fmt.Println(p.Symbol, p.Price)
	}

	// Get the current price of BTCUSDT COIN-M Delivery future
	deliverySymbolPrice, err := deliveryClient.NewListPricesService().Symbol("BTCUSD_240329").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	var futurePrice float64
	for _, p := range deliverySymbolPrice {
		futurePrice, err = strconv.ParseFloat(p.Price, 64)
		if err != nil {
			fmt.Println("error converting", p.Price, "to float64")
			os.Exit(1)
		}
		fmt.Println(p.Symbol, p.Price)
	}

	now := time.Now()
	year, month, date := now.Date()
	calculator := rateCalculator{
		spotPrice:      spotPrice,
		futurePrice:    futurePrice,
		settlementDate: time.Date(2024, time.March, 29, 0, 0, 0, 0, time.UTC),
		tradeDate:      time.Date(year, month, date, 0, 0, 0, 0, time.UTC),
		quantity:       1000,
	}

	fmt.Println("day difference:", calculator.dayDifference())

	fmt.Println("BTC amout:", calculator.underlyingAmount())
	fmt.Println("Earnings:", calculator.earnings())
	fmt.Println("APY: ", calculator.APY(), "%")
	fmt.Println("APR: ", calculator.APR(), "%")

}
