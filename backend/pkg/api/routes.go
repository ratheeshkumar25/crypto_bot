package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/middleware"
)

// SetupRoutes sets up the API routes
func SetupRoutes(app *fiber.App, handler *Handler, authHandler *AuthHandler, wsHandler *WebSocketHandler, jwtSecret string) {
	api := app.Group("/api")

	// Public routes
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)

	// WebSocket route
	app.Get("/api/ws/price", websocket.New(wsHandler.HandlePriceStream))

	// Protected routes
	protected := api.Group("", middleware.JWTMiddleware(jwtSecret))
	protected.Get("/predict/:strategy", handler.PredictProfit)
	protected.Get("/auth/profile", authHandler.GetProfile)
	protected.Put("/auth/exchange-keys", authHandler.UpdateExchangeKeys)
	protected.Get("/trades", handler.GetUserTrades)

	// Public routes (no auth required)
	api.Get("/price/:exchange", handler.GetPrice)
	api.Get("/signals/:strategy", handler.GetSignals)
}
