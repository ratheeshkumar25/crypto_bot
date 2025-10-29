package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/ratheeshkumar25/forex_bot/backend/pkg/repository"
)

// WorkerManager oversees all active analysis workers.
// It is responsible for starting, stopping, and tracking workers for different symbols.
type WorkerManager struct {
	// workerFactory is a dependency that creates new worker instances.
	workerFactory func() *WorkerService

	// activeWorkers holds the cancellation function for each running worker.
	// The map is protected by a mutex to allow safe concurrent access.
	activeWorkers map[string]context.CancelFunc
	mu            sync.Mutex
}

// NewWorkerManager creates a new manager.
// It takes a factory function to create worker instances, which decouples it
// from the specific implementation of WorkerService.
func NewWorkerManager(fetcher *FetcherService, predictor *PredictionService, tradeRepo *repository.TradeRepository, signalRepo *repository.SignalRepository) *WorkerManager {
	return &WorkerManager{
		// This factory function captures the dependencies needed by a WorkerService.
		workerFactory: func() *WorkerService {
			return NewWorkerService(fetcher, predictor, tradeRepo, signalRepo)
		},
		activeWorkers: make(map[string]context.CancelFunc),
	}
}

// StartWorker starts a new analysis worker for the given symbol if one isn't already running.
func (m *WorkerManager) StartWorker(symbol string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.activeWorkers[symbol]; exists {
		return fmt.Errorf("a worker for symbol %s is already running", symbol)
	}

	// Create a new context with a cancellation function for this specific worker.
	ctx, cancel := context.WithCancel(context.Background())

	// Create a new worker instance using the factory.
	worker := m.workerFactory()

	// Start the worker in its own goroutine.
	go worker.Start(ctx, symbol)

	// Store the cancellation function so we can stop it later.
	m.activeWorkers[symbol] = cancel

	return nil
}

// StopWorker stops a running analysis worker for the given symbol.
func (m *WorkerManager) StopWorker(symbol string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cancel, exists := m.activeWorkers[symbol]
	if !exists {
		return fmt.Errorf("no worker is running for symbol %s", symbol)
	}

	// Call the worker's cancellation function to signal it to stop.
	cancel()

	// Remove the worker from the active list.
	delete(m.activeWorkers, symbol)

	return nil
}

// GetStatus returns a slice of strings containing the symbols of all currently running workers.
func (m *WorkerManager) GetStatus() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var symbols []string
	for symbol := range m.activeWorkers {
		symbols = append(symbols, symbol)
	}
	return symbols
}
