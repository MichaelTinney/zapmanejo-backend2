package routes

import (
	"github.com/gofiber/fiber/v2"
	"zapmanejo-cleanbackend/internal/database"
)

func Setup(app *fiber.App) {
	SetupAuthRoutes(app)
	SetupAnimalRoutes(app)
	SetupHealthRoutes(app)
	SetupPaymentRoutes(app)
	SetupWhatsAppRoutes(app)

	// Basic status endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ZapManejo backend live"})
	})

	// Health check endpoint with database connectivity validation
	// Used by DigitalOcean App Platform for health monitoring
	app.Get("/health", func(c *fiber.Ctx) error {
		// Check database connection
		sqlDB, err := database.DB.DB()
		if err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":   "unhealthy",
				"database": "error getting db instance",
				"error":    err.Error(),
			})
		}

		// Ping database to verify connectivity
		if err := sqlDB.Ping(); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":   "unhealthy",
				"database": "connection failed",
				"error":    err.Error(),
			})
		}

		// All checks passed
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"database": "connected",
			"app":      "ZapManejo v1.0",
		})
	})
}
