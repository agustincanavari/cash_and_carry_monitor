package main

import "time"

type FutureData struct {
	FutureSymbol   string    `json:"futureSymbol"`
	FuturePrice    float64   `json:"futurePrice"`
	SettlementDate string    `json:"settlementDate"`
	APR            float64   `json:"apr"`
	APY            float64   `json:"apy"`
	Yield          float64   `json:"yield"`
	LastUpdated    time.Time `json:"lastUpdated"`
}

type CalculatorData struct {
	SpotSymbol  string       `json:"spotSymbol"`
	SpotPrice   float64      `json:"spotPrice"`
	TradeDate   string       `json:"tradeDate"`
	Futures     []FutureData `json:"futures"`
	LastUpdated time.Time    `json:"lastUpdated"`
}
