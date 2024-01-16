// middleware/middleware.go

package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func LogRequests(c *fiber.Ctx) error {
	log.Printf("Request: %s %s", c.Method(), c.OriginalURL())
	return c.Next()
}
