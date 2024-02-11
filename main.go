package main

import (
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

	deliverableFutures := fetchFutures(deliveryClient)
	calculators := generateRateCalculators(deliverableFutures)

	for _, calc := range calculators {
		calc.updateSpotPrice(spotClient)
		calc.updateFuturePrices(deliveryClient)
		calc.print()
	}
}
