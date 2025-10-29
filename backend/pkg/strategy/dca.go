package strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/ratheeshkumar25/forex_bot/backend/pkg/exchange"
)

// DCA implements Dollar-Cost Averaging strategy
type DCA struct {
	Interval time.Duration // Interval between purchases
	Amount   float64       // Amount to invest each time
}

// Execute runs the DCA logic
func (d *DCA) Execute(ctx context.Context, ex exchange.Exchange, symbol string) error {
	// Simple DCA: buy at regular intervals
	ticker := time.NewTicker(d.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			currentPrice, err := ex.GetPrice(ctx, symbol)
			if err != nil {
				return err
			}

			quantity := d.Amount / currentPrice
			err = ex.PlaceOrder(ctx, symbol, "BUY", quantity, currentPrice)
			if err != nil {
				return fmt.Errorf("failed to place DCA buy order: %w", err)
			}
		}
	}
}

// GetProfitPrediction predicts profit for DCA strategy
func (d *DCA) GetProfitPrediction(symbol string, investment float64, timeframe string) (float64, float64, error) {
	// Simple prediction: assume average return over time
	// Placeholder; real prediction would use historical data
	avgReturn := 0.05 // 5% average return
	totalProfit := investment * avgReturn
	percentage := avgReturn * 100

	if timeframe == "short" {
		totalProfit *= 0.3 // Lower for short term
		percentage *= 0.3
	}

	return totalProfit, percentage, nil
}

// GetSignals returns buy signals for DCA strategy
func (d *DCA) GetSignals(symbol string, currentPrice float64) ([]Signal, error) {
	// DCA typically buys at regular intervals, but for signals we can show buy levels
	// For simplicity, show buy signals at current price and slightly below with take profit levels
	signals := []Signal{
		{
			Type:          "BUY",
			Price:         currentPrice * 0.98, // Buy 2% below current
			TakeProfit:    currentPrice * 1.05, // Take profit at 5% above current
			StopLoss:      currentPrice * 0.92, // Stop loss at 8% below buy price
			Timeframe:     "long",
		},
		{
			Type:          "BUY",
			Price:         currentPrice * 0.95, // Buy 5% below current
			TakeProfit:    currentPrice * 1.08, // Take profit at 8% above current
			StopLoss:      currentPrice * 0.89, // Stop loss at 6% below buy price
			Timeframe:     "long",
		},
	}

	return signals, nil
}
