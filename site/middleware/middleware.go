// middleware/middleware.go

package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/secretnamebasis/secret-app/site/config"
)

func LogRequests(c *fiber.Ctx) error {
	log.Printf("Request: %s %s", c.Method(), c.OriginalURL())
	return c.Next()
}

// AuthReq middleware
func AuthReq() func(*fiber.Ctx) error {
	cfg := basicauth.Config{
		Users: map[string]string{
			config.Config("API_USERNAME"): config.Config("API_PASSWORD"),
		},
	}

	authMiddleware := basicauth.New(cfg)
	if authMiddleware == nil {
		// Handle error appropriately (e.g., log it)
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
	}

	return authMiddleware
}
