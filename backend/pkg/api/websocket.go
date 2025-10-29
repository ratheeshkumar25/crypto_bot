package api

import (
	"context"
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/service"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	priceStreamer *service.PriceStreamer
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(priceStreamer *service.PriceStreamer) *WebSocketHandler {
	return &WebSocketHandler{
		priceStreamer: priceStreamer,
	}
}

// HandlePriceStream handles WebSocket connections for price streaming
func (h *WebSocketHandler) HandlePriceStream(c *websocket.Conn) {
	// Get symbol from query parameters
	symbol := c.Query("symbol", "BTCUSDT")
	exchange := c.Query("exchange", "binance")

	log.Printf("WebSocket connection established for %s on %s", symbol, exchange)

	// Add client to streamer
	h.priceStreamer.AddClient(c)

	// Start streaming in a goroutine
	go h.priceStreamer.StartStreaming(context.Background(), symbol, exchange)

	// Keep connection alive and handle disconnections
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Printf("WebSocket error: %v", err)
			h.priceStreamer.RemoveClient(c)
			break
		}
	}
}
