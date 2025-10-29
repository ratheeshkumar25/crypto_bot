package predictor

import (
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/strategy"
)

// Predictor handles profit predictions
type Predictor struct{}

// NewPredictor creates a new predictor
func NewPredictor() *Predictor {
	return &Predictor{}
}

// PredictProfit predicts profit based on strategy, symbol, investment, and timeframe
func (p *Predictor) PredictProfit(strat strategy.Strategy, symbol string, investment float64, timeframe string) (float64, float64, error) {
	return strat.GetProfitPrediction(symbol, investment, timeframe)
}
