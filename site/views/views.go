// views/views.go

package views

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
)

type HomeData struct {
	Title string
	Dev   string
}

func Home(c *fiber.Ctx) error {

	// Define data for rendering the template
	data := HomeData{
		Title: exports.APP_NAME,
		Dev:   dero.Address(),
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

func SubmitForm(c *fiber.Ctx) error {
	// Get form data from the request
	name := c.FormValue("name")
	address := c.FormValue("address")

	// Manually encode the address to preserve newline characters
	encodedAddress := url.QueryEscape(address)

	// Redirect to the pay page with the form data
	return c.Redirect("/pay?name=" + name + "&address=" + encodedAddress)
}

func PayPage(c *fiber.Ctx) error {
	// Retrieve data from the query parameters
	name := c.Query("name")
	address := c.Query("address")
	order := name + "|" + address
	serviceAddr := dero.CreateServiceAddress(order)

	// Replace newline characters with HTML line breaks
	addressWithBreaks := strings.ReplaceAll(address, "\n", "<br>")
	fmt.Printf(addressWithBreaks)
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
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Execute the pay template with the provided data
	err = tmpl.Execute(c.Response().BodyWriter(), data)
	if err != nil {
		fmt.Println("Error executing pay template:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

// views/views.go
// ...

func SuccessPage(c *fiber.Ctx) error {
	// Log template path (for debugging purposes)
	templatePath := "./site/public/success.html"

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println("Error parsing success template:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Execute the success template
	err = tmpl.Execute(c.Response().BodyWriter(), nil)
	if err != nil {
		fmt.Println("Error executing success template:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendString("404 Not Found")
}
