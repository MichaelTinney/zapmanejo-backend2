package routes

import "github.com/gofiber/fiber/v2"

func Setup(app *fiber.App) {
	SetupAuthRoutes(app)
	SetupAnimalRoutes(app)
	SetupHealthRoutes(app)
	SetupPaymentRoutes(app)
	SetupWhatsAppRoutes(app)

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ZapManejo backend live"})
	})
}
