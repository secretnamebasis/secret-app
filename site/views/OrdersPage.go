package views

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/site/models"
)

func OrderPage(c *fiber.Ctx) error {
	addr, _ := dero.Address()
	// Define data for rendering the template
	data := models.HomeData{
		App: exports.APP_NAME,
		Dev: addr,
	}

	tmpl, err := template.ParseFiles("./site/public/order.html")
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

func SubmitOrderForm(c *fiber.Ctx) error {
	// Get form data from the request
	name := c.FormValue("name")
	address := c.FormValue("address")

	// Manually encode the address to preserve newline characters
	encodedAddress := url.QueryEscape(address)

	// Redirect to the pay page with the form data
	return c.Redirect("/pay?name=" + name + "&address=" + encodedAddress)
}
