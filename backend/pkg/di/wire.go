package di

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ratheeshkumar25/forex_bot/backend/config"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/api"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/database"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/exchange"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/predictor"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/repository"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/service"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/strategy"
	"go.uber.org/fx"
)

// Module provides the dependency injection module
var Module = fx.Options(
	config.Module,
	fx.Provide(database.NewDatabase),
	fx.Provide(func(db *database.DB) *repository.UserRepository { return repository.NewUserRepository(db.DB) }),
	fx.Provide(func(db *database.DB) *repository.TradeRepository { return repository.NewTradeRepository(db.DB) }),
	fx.Provide(func(db *database.DB) *repository.SignalRepository { return repository.NewSignalRepository(db.DB) }),
	fx.Provide(NewExchanges),
	fx.Provide(NewStrategies),
	fx.Provide(predictor.NewPredictor),
	fx.Provide(service.NewPriceStreamer),
	fx.Provide(func(exchanges map[string]exchange.Exchange, strategies map[string]strategy.Strategy, pred *predictor.Predictor, tradeRepo *repository.TradeRepository, signalRepo *repository.SignalRepository) *api.Handler {
		return api.NewHandler(exchanges, strategies, pred, tradeRepo, signalRepo)
	}),
	fx.Provide(api.NewAuthHandler),
	fx.Provide(func(cfg *config.Config) string { return cfg.JWTSecret }),
	fx.Provide(api.NewWebSocketHandler),
	fx.Provide(NewApp),
	fx.Invoke(SetupRoutes),
	fx.Invoke(StartServer),
)

// NewExchanges provides exchange instances
func NewExchanges(cfg *config.Config) map[string]exchange.Exchange {
	exchanges := make(map[string]exchange.Exchange)
	if cfg.BinanceAPIKey != "" && cfg.BinanceSecret != "" {
		exchanges["binance"] = exchange.NewBinanceExchange(cfg.BinanceAPIKey, cfg.BinanceSecret)
	}
	// Add Solana exchange
	exchanges["solana"] = exchange.NewSolanaExchange()
	return exchanges
}

// NewStrategies provides strategy instances
func NewStrategies() map[string]strategy.Strategy {
	strategies := make(map[string]strategy.Strategy)
	strategies["grid"] = &strategy.GridStrategy{GridLevels: 5, GridSize: 1.0}
	strategies["dca"] = &strategy.DCA{Interval: 86400000, Amount: 100} // 1 day in milliseconds, $100
	return strategies
}

// NewApp creates the Fiber app
func NewApp() *fiber.App {
	app := fiber.New()
    app.Use(cors.New())
    return app
}

// SetupRoutes sets up the routes
func SetupRoutes(app *fiber.App, handler *api.Handler, authHandler *api.AuthHandler, wsHandler *api.WebSocketHandler, cfg *config.Config) {
	api.SetupRoutes(app, handler, authHandler, wsHandler, cfg.JWTSecret)
}

// StartServer starts the server with fx lifecycle
func StartServer(lc fx.Lifecycle, app *fiber.App, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen("[::]:" + cfg.Port); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
}
