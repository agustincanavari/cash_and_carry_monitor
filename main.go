package main

import (
	"encoding/json"
	"log"
	"net/http"
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
	// Serve static files (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	apiKey, secretKey := loadEnv()
	spotClient, deliveryClient := createClients(apiKey, secretKey)

	deliverableFutures := fetchFutures(deliveryClient)
	calculators := generateRateCalculators(deliverableFutures)

	// Perform an initial synchronous update for each calculator
	for _, calc := range calculators {
		calc.updateSpotPrice(spotClient)
		calc.updateFuturePrices(deliveryClient)
	}

	for _, calc := range calculators {
		startCalculatorUpdate(calc, spotClient, deliveryClient, time.Second*5)
	}

	// API endpoint to provide data
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		data := fetchData(calculators)
		json.NewEncoder(w).Encode(data)
	})

	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}

func fetchData(calculators map[string]*rateCalculator) []CalculatorData {
	var data []CalculatorData
	for _, calc := range calculators {
		calcData := CalculatorData{
			SpotSymbol: calc.spotSymbol,
			SpotPrice:  calc.spotPrice,
			TradeDate:  calc.tradeDate.Format("2006-01-02"),
			Futures:    []FutureData{},
		}

		for _, f := range calc.futures {
			futureData := FutureData{
				FutureSymbol:   f.futureSymbol,
				FuturePrice:    f.futurePrice,
				SettlementDate: f.settlementDate.Format("2006-01-02"),
				APR:            f.APR(calc.spotPrice, calc.tradeDate),
				APY:            f.APY(calc.spotPrice, calc.tradeDate),
			}
			calcData.Futures = append(calcData.Futures, futureData)
		}

		data = append(data, calcData)
	}
	return data
}
