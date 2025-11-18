package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/gofiber/fiber/v2"
	"zapmanejo-cleanbackend/internal/models"
	"zapmanejo-cleanbackend/internal/database"
)

func SetupWhatsAppRoutes(app *fiber.App) {
	app.Get("/webhook", verifyWebhook)
	app.Post("/webhook", handleWhatsAppMessage)
}

func verifyWebhook(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == os.Getenv("WHATSAPP_VERIFY_TOKEN") {
		return c.SendString(challenge)
	}
	return c.Status(403).SendString("Forbidden")
}

func handleWhatsAppMessage(c *fiber.Ctx) error {
	// Simple parser for inbound messages from vaqueros
	var payload struct {
		Entry []struct {
			Changes []struct {
				Value struct {
					Messages []struct {
						From string `json:"from"`
						Text struct {
							Body string `json:"body"`
						} `json:"text"`
					} `json:"messages"`
				} `json:"value"`
			} `json:"changes"`
		} `json:"entry"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(400)
	}

	// Extract message and parse (example: "#brinco 22145 #vacinada #aftosa")
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			for _, msg := range change.Value.Messages {
				body := msg.Text.Body
				// Future: real parser here → create Animal or HealthRecord
				// For now: just log
				_ = body
			}
		}
	}

	return c.SendStatus(200)
}
