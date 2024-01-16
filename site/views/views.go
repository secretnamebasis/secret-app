// views/views.go

package views

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
)

func Home(c *fiber.Ctx) error {
	message := fmt.Sprintf("Welcome to secret-swap!\n%s", dero.CreateServiceAddress(dero.Address()))
	return c.SendString(message)
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendString("404 Not Found")
}
