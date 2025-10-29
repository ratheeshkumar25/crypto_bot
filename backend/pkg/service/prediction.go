package service

import (
	"math"

	"github.com/ratheeshkumar25/forex_bot/backend/pkg/model"
)

// PredictionService encapsulates the logic for making trading predictions.
type PredictionService struct {
	// Dependencies like a logger could be added here later.
}

// NewPredictionService creates a new PredictionService.
func NewPredictionService() *PredictionService {
	return &PredictionService{}
}

// PredictionParameters holds the configuration for the prediction indicators.
type PredictionParameters struct {
	RSI_Period           int
	MACD_Fast_Period     int
	MACD_Slow_Period     int
	MACD_Signal_Period   int
	BBands_Period        int
	BBands_StdDev_Factor float64
}

// DefaultPredictionParams returns a default set of parameters.
func (s *PredictionService) DefaultPredictionParams() *PredictionParameters {
	return &PredictionParameters{
		RSI_Period:           14,
		MACD_Fast_Period:     12,
		MACD_Slow_Period:     26,
		MACD_Signal_Period:   9,
		BBands_Period:        20,
		BBands_StdDev_Factor: 2.0,
	}
}

// AdvancedPredictBuySell uses a combination of technical indicators to generate a trading signal.
func (s *PredictionService) AdvancedPredictBuySell(pair string, data []model.ForexData, params *PredictionParameters) model.Prediction {
	currentPrice := 0.0
	if len(data) > 0 {
		currentPrice = data[0].Price
	}

	requiredDataPoints := max(params.RSI_Period, params.MACD_Slow_Period, params.BBands_Period) + params.MACD_Signal_Period
	if len(data) < requiredDataPoints {
		return model.Prediction{Pair: pair, Signal: "hold", Confidence: 0, Price: currentPrice, Reason: "not enough historical data for a reliable prediction"}
	}

	prices := make([]float64, len(data))
	for i, d := range data {
		prices[len(data)-1-i] = d.Price // Reverse for chronological calculations
	}

	var buyScore, sellScore int

	// 1. RSI
	rsi := calculateRSI(prices, params.RSI_Period)
	if rsi < 30 {
		buyScore++
	} else if rsi > 70 {
		sellScore++
	}

	// 2. MACD
	macdLine, signalLine := calculateMACD(prices, params.MACD_Fast_Period, params.MACD_Slow_Period, params.MACD_Signal_Period)
	if len(macdLine) > 1 && len(signalLine) > 1 {
		if macdLine[len(macdLine)-1] > signalLine[len(signalLine)-1] && macdLine[len(macdLine)-2] <= signalLine[len(signalLine)-2] {
			buyScore++ // Bullish crossover
		} else if macdLine[len(macdLine)-1] < signalLine[len(signalLine)-1] && macdLine[len(macdLine)-2] >= signalLine[len(signalLine)-2] {
			sellScore++ // Bearish crossover
		}
	}

	// 3. Bollinger Bands
	_, upperBand, lowerBand := calculateBollingerBands(prices, params.BBands_Period, params.BBands_StdDev_Factor)
	latestPrice := prices[len(prices)-1]
	if latestPrice < lowerBand[len(lowerBand)-1] {
		buyScore += 2 // Strong buy signal
	} else if latestPrice > upperBand[len(upperBand)-1] {
		sellScore += 2 // Strong sell signal
	}

	var signal string
	var confidence float64
	totalPossibleScore := 4.0

	if buyScore > sellScore {
		signal = "buy"
		confidence = float64(buyScore-sellScore) / totalPossibleScore
	} else if sellScore > buyScore {
		signal = "sell"
		confidence = float64(sellScore-buyScore) / totalPossibleScore
	} else {
		signal = "hold"
		confidence = 0
	}

	return model.Prediction{
		Pair:       pair,
		Signal:     signal,
		Confidence: confidence,
		Price:      latestPrice,
		Reason:     "Advanced model prediction",
	}
}

// --- Calculation Helper Functions (unchanged from your version) ---

func calculateRSI(prices []float64, period int) float64 {
	if len(prices) <= period {
		return 50 // Neutral RSI
	}
	var gains, losses []float64
	for i := 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains = append(gains, change)
			losses = append(losses, 0)
		} else {
			gains = append(gains, 0)
			losses = append(losses, -change)
		}
	}
	avgGain := simpleMovingAverage(gains, period)
	avgLoss := simpleMovingAverage(losses, period)
	if len(avgGain) == 0 || len(avgLoss) == 0 || avgLoss[len(avgLoss)-1] == 0 {
		return 50 // Neutral RSI
	}
	rs := avgGain[len(avgGain)-1] / avgLoss[len(avgLoss)-1]
	return 100 - (100 / (1 + rs))
}

func calculateMACD(prices []float64, fastPeriod, slowPeriod, signalPeriod int) (macdLine, signalLine []float64) {
	emaFast := exponentialMovingAverage(prices, fastPeriod)
	emaSlow := exponentialMovingAverage(prices, slowPeriod)
	if len(emaSlow) > len(emaFast) {
		emaSlow = emaSlow[len(emaSlow)-len(emaFast):]
	}
	offset := len(emaFast) - len(emaSlow)
	macdLine = make([]float64, len(emaSlow))
	for i := range emaSlow {
		macdLine[i] = emaFast[i+offset] - emaSlow[i]
	}
	signalLine = exponentialMovingAverage(macdLine, signalPeriod)
	return
}

func calculateBollingerBands(prices []float64, period int, stdDevFactor float64) (middle, upper, lower []float64) {
	sma := simpleMovingAverage(prices, period)
	stdDev := standardDeviation(prices, period)
	if len(sma) > len(stdDev) {
		sma = sma[len(sma)-len(stdDev):]
	} else if len(stdDev) > len(sma) {
		stdDev = stdDev[len(stdDev)-len(sma):]
	}
	upper = make([]float64, len(sma))
	lower = make([]float64, len(sma))
	for i := range sma {
		upper[i] = sma[i] + (stdDev[i] * stdDevFactor)
		lower[i] = sma[i] - (stdDev[i] * stdDevFactor)
	}
	return sma, upper, lower
}

func exponentialMovingAverage(data []float64, period int) []float64 {
	if len(data) < period {
		return nil
	}
	emas := make([]float64, len(data)-period+1)
	multiplier := 2.0 / (float64(period) + 1.0)
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += data[i]
	}
	emas[0] = sum / float64(period)
	for i := 1; i <= len(data)-period; i++ {
		emas[i] = (data[i+period-1]-emas[i-1])*multiplier + emas[i-1]
	}
	return emas
}

func simpleMovingAverage(data []float64, period int) []float64 {
	if len(data) < period {
		return nil
	}
	smas := make([]float64, len(data)-period+1)
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += data[i]
	}
	smas[0] = sum / float64(period)
	for i := 1; i <= len(data)-period; i++ {
		sum = sum - data[i-1] + data[i+period-1]
		smas[i] = sum / float64(period)
	}
	return smas
}

func standardDeviation(data []float64, period int) []float64 {
	if len(data) < period {
		return nil
	}
	stdDevs := make([]float64, len(data)-period+1)
	for i := 0; i <= len(data)-period; i++ {
		slice := data[i : i+period]
		sum := 0.0
		for _, val := range slice {
			sum += val
		}
		mean := sum / float64(period)
		varianceSum := 0.0
		for _, val := range slice {
			varianceSum += math.Pow(val-mean, 2)
		}
		variance := varianceSum / float64(period)
		stdDevs[i] = math.Sqrt(variance)
	}
	return stdDevs
}

func max(vars ...int) int {
	maxVal := 0
	if len(vars) > 0 {
		maxVal = vars[0]
	}
	for _, i := range vars {
		if i > maxVal {
			maxVal = i
		}
	}
	return maxVal
}

// CalculatePositionSize returns the size of the position based on risk management.
func CalculatePositionSize(accountBalance, riskPercent, entryPrice, stopLoss float64) float64 {
	riskAmount := accountBalance * (riskPercent / 100)
	stopLossDistance := math.Abs(entryPrice - stopLoss)
	if stopLossDistance == 0 {
		return 0 // Prevent division by zero
	}
	size := riskAmount / stopLossDistance
	return size
}
