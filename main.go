package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	for _, p := range spotSymbolPrice {
		fmt.Println(p.Symbol, p.Price)
	}

	// Get the current price of BTCUSDT COIN-M Delivery future
	deliverySymbolPrice, err := deliveryClient.NewListPricesService().Symbol("BTCUSD_240329").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, p := range deliverySymbolPrice {
		fmt.Println(p.Symbol, p.Price)
	}

	/*
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
	*/

}
