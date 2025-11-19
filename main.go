package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"zapmanejo-cleanbackend/internal/database"
	"zapmanejo-cleanbackend/internal/routes"
)

func validateEnvVars() {
	requiredEnvs := map[string]string{
		"DATABASE_URL":          "PostgreSQL connection string",
		"JWT_SECRET":            "Secret key for JWT token signing",
		"WHATSAPP_VERIFY_TOKEN": "Token for WhatsApp webhook verification",
	}

	var missing []string
	for env, description := range requiredEnvs {
		if os.Getenv(env) == "" {
			missing = append(missing, env+" ("+description+")")
		}
	}

	if len(missing) > 0 {
		log.Fatal("FATAL: Required environment variables not set:\n  - " +
			missing[0] +
			func() string {
				result := ""
				for i := 1; i < len(missing); i++ {
					result += "\n  - " + missing[i]
				}
				return result
			}())
	}

	log.Println("✓ All required environment variables are set")
}

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// Validate required environment variables
	validateEnvVars()

	// Connect to PostgreSQL
	database.Connect()

	app := fiber.New(fiber.Config{
		AppName: "ZapManejo v1.0",
	})

	// Middleware - CORS with environment-based origin control
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000" // Default for local development
		log.Println("ALLOWED_ORIGINS not set, using default:", allowedOrigins)
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Routes
	routes.Setup(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("ZapManejo backend running on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
