// views/views.go

package views

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
)

func Home(c *fiber.Ctx) error {
	// Define data for rendering the template
	data := struct {
		Title   string
		Address string // Add this field
	}{
		Title:   "secret-swap",
		Address: dero.CreateServiceAddress(dero.Address()),
	}

	tmpl, err := template.ParseFiles("./site/public/index.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Execute the template with the provided data
	err = tmpl.Execute(c.Response().BodyWriter(), data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendString("404 Not Found")
}
