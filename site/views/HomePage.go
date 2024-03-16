package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/site/models"
)

func Home(c *fiber.Ctx) error {
	addr, _ := dero.Address()
	// Define data for rendering the template
	data := models.HomeData{
		App:    exports.DEVELOPER_NAME,
		Dev:    addr,
		Height: dero.Height(),
	}

	tmpl, err := template.ParseFiles("./site/public/index.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return c.Status(
			http.StatusInternalServerError,
		).SendString("Internal Server Error")
	}

	// Execute the template with the provided data
	err = tmpl.Execute(c.Response().BodyWriter(), data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return c.Status(
			http.StatusInternalServerError,
		).SendString("Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
