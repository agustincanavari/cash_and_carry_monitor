package main

import (
	"log"
	"os"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/joho/godotenv"
)

func loadEnv() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	return apiKey, secretKey
}

func createClients(apiKey, secretKey string) (*binance.Client, *delivery.Client) {
	spotClient := binance.NewClient(apiKey, secretKey)
	deliveryClient := binance.NewDeliveryClient(apiKey, secretKey)
	return spotClient, deliveryClient
}

func main() {
	apiKey, secretKey := loadEnv()
	spotClient, deliveryClient := createClients(apiKey, secretKey)

	deliverableFutures := fetchFutures(deliveryClient)
	calculators := generateRateCalculators(deliverableFutures)

	// Perform an initial synchronous update for each calculator
	for _, calc := range calculators {
		calc.updateSpotPrice(spotClient)
		calc.updateFuturePrices(deliveryClient)
		calc.print()
	}

	for _, calc := range calculators {
		startCalculatorUpdate(calc, spotClient, deliveryClient, time.Second*30)
	}

	printCalculators(calculators, time.Second*20)

	select {}
}

func printCalculators(calculators map[string]*rateCalculator, printInterval time.Duration) {
	ticker := time.NewTicker(printInterval)
	time.Sleep(printInterval)
	go func() {
		//lint:ignore S1000 for ticker-based loop
		for {
			select {
			case <-ticker.C:
				for _, calc := range calculators {
					calc.print()
				}
			}
		}
	}()
}
