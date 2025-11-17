package utils

import "github.com/gofiber/fiber/v2/middleware/cors"

func CORS() cors.Config {
	return cors.Config{
		AllowOrigins:     "https://yourdomain.com, http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowCredentials: true,
	}
}
