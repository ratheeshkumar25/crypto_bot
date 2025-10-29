package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/exchange"
)

// PriceStreamer handles real-time price streaming via WebSocket
type PriceStreamer struct {
	exchanges map[string]exchange.Exchange
	clients   map[*websocket.Conn]bool
	mu        sync.RWMutex
}

// NewPriceStreamer creates a new price streamer
func NewPriceStreamer(exchanges map[string]exchange.Exchange) *PriceStreamer {
	return &PriceStreamer{
		exchanges: exchanges,
		clients:   make(map[*websocket.Conn]bool),
	}
}

// AddClient adds a new WebSocket client
func (ps *PriceStreamer) AddClient(conn *websocket.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.clients[conn] = true
	log.Printf("New WebSocket client connected. Total clients: %d", len(ps.clients))
}

// RemoveClient removes a WebSocket client
func (ps *PriceStreamer) RemoveClient(conn *websocket.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	delete(ps.clients, conn)
	conn.Close()
	log.Printf("WebSocket client disconnected. Total clients: %d", len(ps.clients))
}

// StartStreaming starts the price streaming for a symbol
func (ps *PriceStreamer) StartStreaming(ctx context.Context, symbol string, exchangeName string) {
	ex, ok := ps.exchanges[exchangeName]
	if !ok {
		log.Printf("Exchange %s not found", exchangeName)
		return
	}

	ticker := time.NewTicker(1 * time.Second) // Stream every second
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			price, err := ex.GetPrice(ctx, symbol)
			if err != nil {
				log.Printf("Error getting price for %s: %v", symbol, err)
				continue
			}

			ps.broadcastPrice(symbol, price)
		case <-ctx.Done():
			return
		}
	}
}

// broadcastPrice sends price update to all connected clients
func (ps *PriceStreamer) broadcastPrice(symbol string, price float64) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	message := map[string]interface{}{
		"symbol":    symbol,
		"price":     price,
		"timestamp": time.Now().Unix(),
	}

	for conn := range ps.clients {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("Error sending message to client: %v", err)
			// Remove broken connection
			go func(c *websocket.Conn) {
				ps.RemoveClient(c)
			}(conn)
		}
	}
}
