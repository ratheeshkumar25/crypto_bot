package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/model"
)

// FetcherService is responsible for fetching data from external APIs.
// This version is configured to fetch kline data from Binance.
type FetcherService struct {
	binanceClient *binance.Client
}

// NewFetcherService creates a new FetcherService for Binance.
// For public data endpoints, API keys are not required.
func NewFetcherService() *FetcherService {
	return &FetcherService{
		binanceClient: binance.NewClient("", ""),
	}
}

// FetchKlines retrieves and sorts the latest kline (candlestick) data for a given crypto symbol.
func (s *FetcherService) FetchKlines(symbol, interval string) ([]model.ForexData, error) {
	// 1. Fetch klines from Binance API.
	klines, err := s.binanceClient.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(100). // The number of candles to retrieve, sufficient for our indicators
		Do(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to fetch klines from Binance for symbol %s: %w", symbol, err)
	}

	if len(klines) == 0 {
		return nil, fmt.Errorf("no kline data returned from Binance for symbol %s (is the symbol valid?)", symbol)
	}

	// 2. Transform the Binance response into our internal model.ForexData format.
	var data []model.ForexData
	for _, k := range klines {
		// The prediction engine uses the closing price of each candle.
		price, err := strconv.ParseFloat(k.Close, 64)
		if err != nil {
			// Skip this kline if the price is not a valid number, though this is unlikely with Binance.
			continue
		}
		data = append(data, model.ForexData{
			// The timestamp from Binance is the *closing* time of the candle.
			Timestamp: time.UnixMilli(k.CloseTime).Format("2006-01-02 15:04:05"),
			Price:     price,
		})
	}

	// 3. The Binance API returns data sorted oldest to newest.
	// Our prediction engine expects data sorted newest to oldest. So, we reverse the slice.
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	return data, nil
}
