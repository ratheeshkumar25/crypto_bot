package di

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ratheeshkumar25/forex_bot/backend/config"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/api"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/database"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/repository"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/service"
	"go.uber.org/fx"
)

func Init() {
	app := fx.New(
		fx.Provide(
			// -- Configuration --
			config.NewConfig,

			// -- Database --
			database.NewDatabase,

			// -- Repositories --
			func(db *database.DB) *repository.UserRepository {
				return repository.NewUserRepository(db.DB)
			},
			func(db *database.DB) *repository.TradeRepository {
				return repository.NewTradeRepository(db.DB)
			},
			func(db *database.DB) *repository.SignalRepository {
				return repository.NewSignalRepository(db.DB)
			},

			// -- Services --
			service.NewFetcherService,
			service.NewPredictionService,
			service.NewWorkerService,
			service.NewWorkerManager, // The new manager for our workers

			// -- API Layer --
			api.NewPredictionHandler,
			api.NewWorkerHandler, // The new handler to control the workers

			// -- Web Server --
			func() *fiber.App {
				return fiber.New()
			},
		),
		// Invoke the function to register all our API routes
		fx.Invoke(registerHooks),
	)

	// Run the application. It will start the web server and wait for API calls.
	app.Run()
}

// registerHooks configures the web server and registers all API handlers.
func registerHooks(lifecycle fx.Lifecycle, app *fiber.App, cfg *config.Config, predHandler *api.PredictionHandler, workerHandler *api.WorkerHandler) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				// Register routes from both handlers
				predHandler.RegisterRoutes(app)
				workerHandler.RegisterRoutes(app)

				go func() {
					log.Printf("API server listening on [::]:%s. Use the /api/worker/start endpoint to begin analysis.", cfg.Port)
					if err := app.Listen("[::]:" + cfg.Port); err != nil {
						log.Fatalf("Failed to start API server: %v", err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Println("Shutting down API server...")
				return app.Shutdown()
			},
		},
	)
}

// startWorker starts our new automated analysis bot in the background.
// func startWorker(lifecycle fx.Lifecycle, worker *service.WorkerService) {
// 	lifecycle.Append(fx.Hook{
// 		OnStart: func(ctx context.Context) error {
// 			// The symbol to monitor can be configured later.
// 			go worker.Start(ctx, "BTCUSDT")
// 			return nil
// 		},
// 	})
// }
