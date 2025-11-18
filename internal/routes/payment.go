package routes

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"zapmanejo-cleanbackend/internal/database"
	"zapmanejo-cleanbackend/internal/middleware"
	"zapmanejo-cleanbackend/internal/models"
	"github.com/paypal/paypal-checkout-sdk/v2/payments"
	"github.com/paypal/paypal-checkout-sdk/v2/payments/capture"
)

func SetupPaymentRoutes(app *fiber.App) {
	pay := app.Group("/api/payment", middleware.JWTProtected())

	pay.Post("/create", createPayPalOrder)
	pay.Post("/capture", capturePayPalOrder)
}

func createPayPalOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var input struct {
		Type string `json:"type"` // "monthly" or "lifetime"
	}
	c.BodyParser(&input)

	amount := "250.00"
	if input.Type == "lifetime" {
		amount = "2500.00"
	}

	// Simple PayPal order creation (real implementation uses PayPal SDK)
	orderID := "ORDER_" + strconv.Itoa(int(userID)) + "_" + input.Type // placeholder
	return c.JSON(fiber.Map{
		"orderID": orderID,
		"amount":  amount,
		"type":    input.Type,
	})
}

func capturePayPalOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var input struct {
		OrderID string `json:"orderID"`
		Type    string `json:"type"`
	}
	c.BodyParser(&input)

	// Simulate capture success
	payment := models.Payment{
		UserID:        userID,
		PayPalOrderID: input.OrderID,
		Amount:        250.00,
		Type:          input.Type,
		Status:        "COMPLETED",
	}
	if input.Type == "lifetime" {
		payment.Amount = 2500.00
		var slot models.LifetimeSlot
		database.DB.Where("occupied = false").First(&slot)
		if slot.ID == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "No lifetime slots left"})
		}
		slot.Occupied = true
		slot.UserID = &userID
		database.DB.Save(&slot)

		database.DB.Model(&models.User{}).Where("id = ?", userID).Update("plan", "lifetime")
	}

	database.DB.Create(&payment)
	return c.JSON(fiber.Map{"success": true, "plan": input.Type == "lifetime" ? "lifetime" : "monthly"})
}
