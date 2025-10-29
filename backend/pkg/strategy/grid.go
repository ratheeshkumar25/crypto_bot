package strategy

import (
	"context"
	"fmt"

	"github.com/ratheeshkumar25/forex_bot/backend/pkg/exchange"
)

// GridStrategy implements grid trading strategy
type GridStrategy struct {
	GridLevels int     // Number of grid levels
	GridSize   float64 // Percentage size of each grid
}

// Execute runs the grid trading logic
func (g *GridStrategy) Execute(ctx context.Context, ex exchange.Exchange, symbol string) error {
	// Simple grid logic: buy at lower levels, sell at higher levels
	currentPrice, err := ex.GetPrice(ctx, symbol)
	if err != nil {
		return err
	}

	for i := 1; i <= g.GridLevels; i++ {
		buyPrice := currentPrice * (1 - float64(i)*g.GridSize/100)
		sellPrice := currentPrice * (1 + float64(i)*g.GridSize/100)

		// Place buy order
		err = ex.PlaceOrder(ctx, symbol, "BUY", 0.01, buyPrice) // Example quantity
		if err != nil {
			return fmt.Errorf("failed to place buy order: %w", err)
		}

		// Place sell order
		err = ex.PlaceOrder(ctx, symbol, "SELL", 0.01, sellPrice) // Example quantity
		if err != nil {
			return fmt.Errorf("failed to place sell order: %w", err)
		}
	}

	return nil
}

// GetProfitPrediction predicts profit for grid strategy
func (g *GridStrategy) GetProfitPrediction(symbol string, investment float64, timeframe string) (float64, float64, error) {
	// Simple prediction: assume average profit per grid level
	// This is a placeholder; real prediction would use historical data, ML, etc.
	avgProfitPerLevel := 0.02 // 2% per level
	totalProfit := float64(g.GridLevels) * avgProfitPerLevel * investment
	percentage := totalProfit / investment * 100

	if timeframe == "short" {
		totalProfit *= 0.5 // Adjust for short term
		percentage *= 0.5
	}

	return totalProfit, percentage, nil
}

// GetSignals returns buy/sell signals for grid strategy
func (g *GridStrategy) GetSignals(symbol string, currentPrice float64) ([]Signal, error) {
	var signals []Signal

	for i := 1; i <= g.GridLevels; i++ {
		buyPrice := currentPrice * (1 - float64(i)*g.GridSize/100)
		sellPrice := currentPrice * (1 + float64(i)*g.GridSize/100)

		// Calculate profit percentage and take profit levels for each level
		// buyProfit := (currentPrice - buyPrice) / buyPrice * 100
		buyTakeProfit := buyPrice * (1 + g.GridSize/100*2) // Take profit at 2x grid size above buy
		buyStopLoss := buyPrice * (1 - g.GridSize/100*1.5) // Stop loss at 1.5x grid size below buy

		// sellProfit := (sellPrice - currentPrice) / currentPrice * 100
		sellTakeProfit := sellPrice * (1 - g.GridSize/100*2) // Take profit at 2x grid size below sell
		sellStopLoss := sellPrice * (1 + g.GridSize/100*1.5) // Stop loss at 1.5x grid size above sell

		signals = append(signals, Signal{
			Type:       "BUY",
			Price:      buyPrice,
			TakeProfit: buyTakeProfit,
			StopLoss:   buyStopLoss,
			Timeframe:  "long", // Grid works better for longer term
		})
		signals = append(signals, Signal{
			Type:       "SELL",
			Price:      sellPrice,
			TakeProfit: sellTakeProfit,
			StopLoss:   sellStopLoss,
			Timeframe:  "long",
		})
	}

	return signals, nil
}
