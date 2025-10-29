package exchange

import (
	"context"
	"time"
)

// Exchange defines the interface for interacting with cryptocurrency exchanges
type Exchange interface {
	GetPrice(ctx context.Context, symbol string) (float64, error)
	GetVolume(ctx context.Context, symbol string, timeframe string) (float64, error)
	PlaceOrder(ctx context.Context, symbol string, side string, quantity float64, price float64) error
	GetBalance(ctx context.Context, asset string) (float64, error)
	// Add more methods as needed, e.g., GetOrderBook, etc.
}

// PriceData represents price information
type PriceData struct {
	Symbol string
	Price  float64
	Time   time.Time
}

// VolumeData represents volume information
type VolumeData struct {
	Symbol    string
	Volume    float64
	Time      time.Time
	Timeframe string // e.g., "1m", "5m", "1h"
}
