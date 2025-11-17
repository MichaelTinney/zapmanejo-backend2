package routes

import (
	"github.com/gofiber/fiber/v2"
	"zapmanejo-backend2/internal/database"
	"zapmanejo-backend2/internal/middleware"
	"zapmanejo-backend2/internal/models"
)

func SetupHealthRoutes(app *fiber.App) {
	health := app.Group("/api/health", middleware.JWTProtected())

	health.Get("/", getHealthRecords)
	health.Post("/", createHealthRecord)
}

func getHealthRecords(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var records []models.HealthRecord
	database.DB.Where("user_id = ?", userID).Find(&records)
	return c.JSON(records)
}

func createHealthRecord(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var record models.HealthRecord
	if err := c.BodyParser(&record); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid data"})
	}
	record.UserID = userID
	database.DB.Create(&record)
	return c.JSON(record)
}
