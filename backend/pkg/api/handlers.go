package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/exchange"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/predictor"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/repository"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/strategy"
)

// Handler holds dependencies for API handlers
type Handler struct {
	Exchanges  map[string]exchange.Exchange
	Strategies map[string]strategy.Strategy
	Predictor  *predictor.Predictor
	TradeRepo  *repository.TradeRepository
}

// NewHandler creates a new handler
func NewHandler(exchanges map[string]exchange.Exchange, strategies map[string]strategy.Strategy, pred *predictor.Predictor, tradeRepo *repository.TradeRepository) *Handler {
	return &Handler{
		Exchanges:  exchanges,
		Strategies: strategies,
		Predictor:  pred,
		TradeRepo:  tradeRepo,
	}
}

// GetPrice handles getting price from an exchange
func (h *Handler) GetPrice(c *fiber.Ctx) error {
	exchangeName := c.Params("exchange")
	symbol := c.Query("symbol")

	ex, ok := h.Exchanges[exchangeName]
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Exchange not found"})
	}

	price, err := ex.GetPrice(c.Context(), symbol)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"symbol":    symbol,
		"price":     price,
		"timestamp": c.Context().Time().Unix(),
	})
}

// PredictProfit handles profit prediction
func (h *Handler) PredictProfit(c *fiber.Ctx) error {
	strategyName := c.Params("strategy")
	symbol := c.Query("symbol")
	investment := c.QueryFloat("investment", 0)
	timeframe := c.Query("timeframe", "long")

	strat, ok := h.Strategies[strategyName]
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Strategy not found"})
	}

	profit, percentage, err := h.Predictor.PredictProfit(strat, symbol, investment, timeframe)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"strategy":         strategyName,
		"symbol":           symbol,
		"investment":       investment,
		"timeframe":        timeframe,
		"predictedProfit":  profit,
		"profitPercentage": percentage,
	})
}

// GetSignals handles getting trading signals for a strategy
func (h *Handler) GetSignals(c *fiber.Ctx) error {
	strategyName := c.Params("strategy")
	symbol := c.Query("symbol")

	strat, ok := h.Strategies[strategyName]
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Strategy not found"})
	}

	// Get current price to base signals on
	ex, ok := h.Exchanges["binance"] // Assuming binance for now
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Exchange not found"})
	}

	currentPrice, err := ex.GetPrice(c.Context(), symbol)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	signals, err := strat.GetSignals(symbol, currentPrice)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"strategy": strategyName,
		"symbol":   symbol,
		"signals":  signals,
	})
}
