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
	spotApiKey := os.Getenv("SPOT_API_KEY")
	spotSecretKey := os.Getenv("SPOT_SECRET_KEY")

	client := binance.NewClient(spotApiKey, spotSecretKey)

	// Get the current price of BTCUSDT
	symbolPrice, err := client.NewListPricesService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, p := range symbolPrice {
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
