package views

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
)

func PayPage(c *fiber.Ctx) error {
	// Retrieve data from the query parameters
	name := c.Query("name")
	address := c.Query("address")
	order := name + "|" + address
	serviceAddr, err := dero.CreateServiceAddress(order)
	if err != nil {
		return err
	}

	// Replace newline characters with HTML line breaks
	addressWithBreaks := strings.ReplaceAll(address, "\n", "<br>")
	fmt.Println(addressWithBreaks)
	// Prepare data for rendering the pay page template
	data := struct {
		Name           string
		Address        template.HTML // Use template.HTML to interpret HTML as raw HTML
		ServiceAddress string
	}{
		Name:           name,
		Address:        template.HTML(addressWithBreaks),
		ServiceAddress: serviceAddr,
	}

	// Log template path (for debugging purposes)
	templatePath := "./site/public/pay.html"

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println("Error parsing pay template:", err)
		return c.Status(
			http.StatusInternalServerError,
		).SendString("Internal Server Error")
	}

	// Execute the pay template with the provided data
	err = tmpl.Execute(c.Response().BodyWriter(), data)
	if err != nil {
		fmt.Println("Error executing pay template:", err)
		return c.Status(
			http.StatusInternalServerError,
		).SendString("Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
