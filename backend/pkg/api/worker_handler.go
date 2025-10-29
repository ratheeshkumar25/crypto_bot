package api

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/service"
)

// WorkerHandler handles API requests related to controlling the analysis workers.
type WorkerHandler struct {
	manager *service.WorkerManager
}

// NewWorkerHandler creates a new handler for the worker manager.
// It depends on the WorkerManager, which fx will provide.
func NewWorkerHandler(manager *service.WorkerManager) *WorkerHandler {
	return &WorkerHandler{
		manager: manager,
	}
}

// StartWorker handles the POST /api/worker/start endpoint.
func (h *WorkerHandler) StartWorker(c *fiber.Ctx) error {
	symbol := strings.ToUpper(c.Query("pair"))
	if symbol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query parameter 'pair' is required."})
	}

	err := h.manager.StartWorker(symbol)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("API request to START worker for %s", symbol)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Worker started successfully for " + symbol})
}

// StopWorker handles the POST /api/worker/stop endpoint.
func (h *WorkerHandler) StopWorker(c *fiber.Ctx) error {
	symbol := strings.ToUpper(c.Query("pair"))
	if symbol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query parameter 'pair' is required."})
	}

	err := h.manager.StopWorker(symbol)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("API request to STOP worker for %s", symbol)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Worker stopped successfully for " + symbol})
}

// GetStatus handles the GET /api/worker/status endpoint.
func (h *WorkerHandler) GetStatus(c *fiber.Ctx) error {
	activeWorkers := h.manager.GetStatus()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"active_workers": len(activeWorkers),
		"monitoring":     activeWorkers,
	})
}

// RegisterRoutes sets up all the routes for the worker control API.
func (h *WorkerHandler) RegisterRoutes(app *fiber.App) {
	workerAPI := app.Group("/api/worker")
	workerAPI.Post("/start", h.StartWorker)
	workerAPI.Post("/stop", h.StopWorker)
	workerAPI.Get("/status", h.GetStatus)
}
