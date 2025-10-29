package service

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ratheeshkumar25/forex_bot/backend/pkg/model"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/repository"
)

// WorkerService orchestrates the continuous analysis of market data.
type WorkerService struct {
	fetcherSvc *FetcherService
	predSvc    *PredictionService
	tradeRepo  *repository.TradeRepository
	signalRepo *repository.SignalRepository
	trades     []*model.Trade
	mu         sync.Mutex
}

// NewWorkerService creates a new automated worker.
func NewWorkerService(fetcher *FetcherService, predictor *PredictionService, tradeRepo *repository.TradeRepository, signalRepo *repository.SignalRepository) *WorkerService {
	return &WorkerService{
		fetcherSvc: fetcher,
		predSvc:    predictor,
		tradeRepo:  tradeRepo,
		signalRepo: signalRepo,
	}
}

// AnalysisResult holds the prediction for a single timeframe.
type AnalysisResult struct {
	Timeframe  string
	Prediction model.Prediction
	Error      error
}

var timeframes = []string{"1m", "5m", "1d"}

var timeframeIntervals = map[string]time.Duration{
	"1m": 10 * time.Second,
	"5m": 10 * time.Second,
	"1d": 1 * time.Hour, // or 24 * time.Hour for true daily
}

func (s *WorkerService) Start(ctx context.Context, symbol string) {
	log.Printf("Starting automated analysis worker for %s...", symbol)
	log.Println("--- Bot is now running. Press Ctrl+C to stop. ---")

	var wg sync.WaitGroup
	for _, tf := range timeframes {
		wg.Add(1)
		go func(tf string) {
			defer wg.Done()
			ticker := time.NewTicker(timeframeIntervals[tf])
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					s.runAnalysisForTimeframe(symbol, tf)
				case <-ctx.Done():
					return
				}
			}
		}(tf)
	}
	wg.Wait()
}

// Add this helper function:
func (s *WorkerService) runAnalysisForTimeframe(symbol, tf string) {
	log.Printf("Running analysis for %s [%s]...", symbol, tf)
	var wg sync.WaitGroup
	results := make(chan AnalysisResult, 1)
	wg.Add(1)
	go s.analyzeTimeframe(symbol, tf, &wg, results)
	wg.Wait()
	close(results)
	for result := range results {
		if result.Error != nil {
			log.Printf("  | %-4s -> Error: %v", result.Timeframe, result.Error)
			continue
		}
		p := result.Prediction
		log.Printf("  | %-4s -> Signal: %-4s | Price: %-12.4f | Confidence: %.2f%%",
			result.Timeframe,
			strings.ToUpper(p.Signal),
			p.Price,
			p.Confidence*100,
		)
		s.CheckTrades(p.Price)
	}
}

// analyzeTimeframe is a helper to run prediction for a single timeframe concurrently.
func (s *WorkerService) analyzeTimeframe(symbol, timeframe string, wg *sync.WaitGroup, results chan<- AnalysisResult) {
	defer wg.Done()

	klines, err := s.fetcherSvc.FetchKlines(symbol, timeframe)
	if err != nil {
		results <- AnalysisResult{Timeframe: timeframe, Error: err}
		return
	}

	params := s.predSvc.DefaultPredictionParams()
	prediction := s.predSvc.AdvancedPredictBuySell(symbol, klines, params)
	results <- AnalysisResult{Timeframe: timeframe, Prediction: prediction}
}

func (s *WorkerService) AddTrade(trade *model.Trade) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.trades = append(s.trades, trade)
}

func (s *WorkerService) CheckTrades(latestPrice float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for _, trade := range s.trades {
		if !trade.IsOpen {
			continue
		}
		if latestPrice >= trade.TakeProfit {
			trade.IsOpen = false
			log.Printf("Trade for %s closed at TP: %.4f", trade.Symbol, latestPrice)
		} else if latestPrice <= trade.StopLoss {
			trade.IsOpen = false
			log.Printf("Trade for %s closed at SL: %.4f", trade.Symbol, latestPrice)
		} else if now.Sub(trade.OpenTime) >= trade.MaxDuration {
			trade.IsOpen = false
			log.Printf("Trade for %s closed by duration at: %.4f", trade.Symbol, latestPrice)
		}
	}
}
