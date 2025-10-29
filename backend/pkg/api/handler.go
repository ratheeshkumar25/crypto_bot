package api

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/service"
)

// PredictionHandler handles API requests for predictions.
type PredictionHandler struct {
	fetcherSvc *service.FetcherService
	predSvc    *service.PredictionService
	manager    *service.WorkerManager
}

// NewPredictionHandler creates a new handler.
// It depends on the FetcherService, PredictionService, and WorkerManager, which fx will provide.
func NewPredictionHandler(fetcher *service.FetcherService, predictor *service.PredictionService, manager *service.WorkerManager) *PredictionHandler {
	return &PredictionHandler{
		fetcherSvc: fetcher,
		predSvc:    predictor,
		manager:    manager,
	}
}

// GetPrediction is the handler for the /api/prediction endpoint.
func (h *PredictionHandler) GetPrediction(c *fiber.Ctx) error {
	// 1. Get and validate the crypto symbol from the query.
	// We'll use BTCUSDT as a robust default.
	symbol := strings.ToUpper(c.Query("pair", "BTCUSDT"))
	if symbol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "The 'pair' query parameter cannot be empty.",
		})
	}

	// 2. Fetch the latest market data using the new Binance FetcherService.
	// We'll hardcode "5m" as the interval for this specific API endpoint.
	data, err := h.fetcherSvc.FetchKlines(symbol, "5m")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch market data from Binance",
			"details": err.Error(),
		})
	}

	// 3. Generate a prediction using the PredictionService.
	// Notice this part of our code doesn't need to change at all!
	// Our abstraction works perfectly.
	params := h.predSvc.DefaultPredictionParams()
	prediction := h.predSvc.AdvancedPredictBuySell(symbol, data, params)

	// Log the signal to the terminal if it is a "buy" or "sell" event.
	if prediction.Signal == "buy" || prediction.Signal == "sell" {
		log.Printf(
			"[%s] Signal Triggered: %s | Price: %.4f | Confidence: %.2f%%",
			prediction.Pair,
			strings.ToUpper(prediction.Signal),
			prediction.Price,
			prediction.Confidence*100,
		)
	}

	// After providing the manual prediction, also start monitoring the symbol in the background.
	// We run this in a goroutine so it doesn't block the API response.
	go func() {
		err := h.manager.StartWorker(symbol)
		if err != nil {
			// This error is expected if the worker is already running. We can just log it for info.
			log.Printf("Info while auto-starting worker for %s: %v", symbol, err)
		}
	}()

	return c.JSON(prediction)
}

// RegisterRoutes sets up the routes for this handler on the Fiber app.
func (h *PredictionHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/api/prediction", h.GetPrediction)
}
