// views/views.go

package views

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/site/controllers"
	"github.com/secretnamebasis/secret-app/site/models"
)

type HomeData struct {
	Title   string
	Address string
	Items   []models.Item // Add this field
}

func Home(c *fiber.Ctx) error {
	// Retrieve blog posts
	items, err := controllers.DisplayItems(c)
	if err != nil {
		return err
	}

	// Define data for rendering the template
	data := HomeData{
		Title:   exports.APP_NAME,
		Address: dero.Address(),
		Items:   items,
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
