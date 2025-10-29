package strategy

import (
	"context"

	"github.com/ratheeshkumar25/forex_bot/backend/pkg/exchange"
)

// Signal represents a trading signal
type Signal struct {
	Type          string  // "BUY" or "SELL"
	Price         float64 // Price level for the signal
	TakeProfit    float64 // Take profit price level
	StopLoss      float64 // Stop loss price level
	Timeframe     string  // "short" or "long"
}

// Strategy defines the interface for trading strategies
type Strategy interface {
	Execute(ctx context.Context, exchange exchange.Exchange, symbol string) error
	GetProfitPrediction(symbol string, investment float64, timeframe string) (float64, float64, error)
	GetSignals(symbol string, currentPrice float64) ([]Signal, error)
}

// GridStrategy is defined in grid.go

// DCA is defined in dca.go
