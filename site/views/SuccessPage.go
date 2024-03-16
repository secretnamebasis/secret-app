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

func SuccessPage(c *fiber.Ctx) error {
	result := c.Query("result")
	addr, _ := dero.Address()
	data := models.SuccessData{
		Title: exports.APP_NAME,
		Txid:  result,
		Dev:   addr,
	}
	// Log template path (for debugging purposes)
	templatePath := "./site/public/success.html"

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println("Error parsing success template:", err)
		return c.Status(
			http.StatusInternalServerError,
		).SendString("Internal Server Error")
	}

	// Execute the success template
	err = tmpl.Execute(c.Response().BodyWriter(), data)
	if err != nil {
		fmt.Println("Error executing success template:", err)
		return c.Status(
			http.StatusInternalServerError,
		).SendString("Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
