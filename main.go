package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
)

func loadKeys() (string, string) {
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

	apiKey, secretKey := loadKeys()
	spotClient, deliveryClient := createClients(apiKey, secretKey)

	deliverableFutures := fetchFutures(deliveryClient)
	calculators := generateRateCalculators(deliverableFutures)

	// Perform an initial synchronous update for each calculator
	for _, calc := range calculators {
		calc.updateSpotPrice(spotClient)
		calc.updateFuturePrices(deliveryClient)
	}

	for _, calc := range calculators {
		startCalculatorUpdate(calc, spotClient, deliveryClient, time.Second*60)
	}

	// API endpoint to provide data
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		data := fetchData(calculators)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Printf("Error encoding data: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	})

	port := getEnvWithDefault("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      nil, // use http.DefaultServeMux
	}
	log.Printf("Listening on port %s", port)
	err := server.ListenAndServe()
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
	sort.Slice(data, func(i, j int) bool {
		return data[i].SpotSymbol < data[j].SpotSymbol
	})
	return data
}
