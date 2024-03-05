// views/views.go

package views

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/deroproject/derohe/rpc"
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/site/models"
)

func Home(c *fiber.Ctx) error {
	addr, _ := dero.Address()
	// Define data for rendering the template
	data := models.HomeData{
		Title:  exports.DEVELOPER_NAME,
		Dev:    addr,
		Height: dero.Height(),
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

func About(c *fiber.Ctx) error {
	addr, _ := dero.Address()
	// Define data for rendering the template
	data := models.HomeData{
		Title:  exports.DEVELOPER_NAME,
		Dev:    addr,
		Height: dero.Height(),
	}

	tmpl, err := template.ParseFiles("./site/public/about.html")
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

func EntryPage(c *fiber.Ctx) error {
	addr, _ := dero.Address()
	// Define data for rendering the template
	data := models.HomeData{
		Title: exports.APP_NAME,
		Dev:   addr,
	}

	tmpl, err := template.ParseFiles("./site/public/entry.html")
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

func SubmitEntryForm(c *fiber.Ctx) error {
	entry := c.FormValue("entry")

	// Check if the entry exceeds the character limit
	if len(entry) > 120 {
		// Split the entry into chunks of 120 characters
		chunks := chunkString(entry, 120)

		// Create a list to store transfers
		var transfers []rpc.Transfer

		// Iterate over the chunks and create a transfer for each
		for _, chunk := range chunks {
			payload := rpc.Arguments{
				{
					Name:     rpc.RPC_DESTINATION_PORT,
					DataType: rpc.DataUint64,
					Value:    uint64(0),
				},
				{
					Name:     rpc.RPC_COMMENT,
					DataType: rpc.DataString,
					Value:    chunk,
				},
			}
			transfer := rpc.Transfer{
				Destination: exports.CAPTAIN_WALLET_ADDRESS,
				Amount:      uint64(0),
				Payload_RPC: payload,
			}
			transfers = append(transfers, transfer)
		}

		// Create Transfer_Params with the list of transfers
		transferParams := rpc.Transfer_Params{
			Transfers: transfers,
		}

		// Send the transfers
		result := dero.SendTransfer(transferParams)

		// Handle the result as needed
		fmt.Println("Transfer result:", result)

		// Redirect to the success page or handle the result accordingly
		return c.Redirect("/success?result=" + result)
	}

	// If the entry is within the character limit, proceed as before
	payload := rpc.Arguments{
		{
			Name:     rpc.RPC_DESTINATION_PORT,
			DataType: rpc.DataUint64,
			Value:    uint64(0),
		},
		{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    entry,
		},
	}
	transfer := rpc.Transfer{
		Destination: exports.DEVELOPER_ADDRESS,
		Amount:      uint64(0),
		Payload_RPC: payload,
	}
	transferParams := rpc.Transfer_Params{
		Transfers: []rpc.Transfer{transfer},
	}
	result := dero.SendTransfer(transferParams)

	// Redirect to the success page with the form data
	return c.Redirect("/success?result=" + result)
}

func EntriesPage(c *fiber.Ctx) error {
	entries, err := dero.GetOutgoingTransfers()
	if err != nil {
		// Return a 500 Internal Server Error with a meaningful message
		return c.Status(http.StatusInternalServerError).SendString("Error fetching entries: " + err.Error())
	}
	addr, _ := dero.Address()
	// Define data for rendering the template
	data := models.EntriesData{
		Title:   exports.APP_NAME,
		Dev:     addr,
		Entries: entries.Entries,
	}

	tmpl, err := template.ParseFiles("./site/public/entries.html")
	if err != nil {
		// Return a 500 Internal Server Error with a meaningful message
		return c.Status(http.StatusInternalServerError).SendString("Error parsing template: " + err.Error())
	}

	// Execute the template with the provided data
	err = tmpl.Execute(c.Response().BodyWriter(), data)
	if err != nil {
		// Return a 500 Internal Server Error with a meaningful message
		return c.Status(http.StatusInternalServerError).SendString("Error executing template: " + err.Error())
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

func MessagePage(c *fiber.Ctx) error {
	addr, _ := dero.Address()
	// Define data for rendering the template
	data := models.HomeData{
		Title: exports.APP_NAME,
		Dev:   addr,
	}

	tmpl, err := template.ParseFiles("./site/public/message.html")
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

func SubmitMessageForm(c *fiber.Ctx) error {
	entry := c.FormValue("entry")
	address := c.FormValue("address")

	// Check if the entry exceeds the character limit
	if len(entry) > 120 {
		// Split the entry into chunks of 120 characters
		chunks := chunkString(entry, 120)

		// Create a list to store transfers
		var transfers []rpc.Transfer

		// Iterate over the chunks and create a transfer for each
		for _, chunk := range chunks {
			payload := rpc.Arguments{
				{
					Name:     rpc.RPC_DESTINATION_PORT,
					DataType: rpc.DataUint64,
					Value:    uint64(0),
				},
				{
					Name:     rpc.RPC_COMMENT,
					DataType: rpc.DataString,
					Value:    chunk,
				},
			}
			transfer := rpc.Transfer{
				Destination: address,
				Amount:      uint64(1),
				Payload_RPC: payload,
			}
			transfers = append(transfers, transfer)
		}

		// Create Transfer_Params with the list of transfers
		transferParams := rpc.Transfer_Params{
			Transfers: transfers,
		}

		// Send the transfers
		result := dero.SendTransfer(transferParams)

		// Handle the result as needed
		fmt.Println("Transfer result:", result)

		// Redirect to the success page or handle the result accordingly
		return c.Redirect("/success?result=" + result)
	}

	// If the entry is within the character limit, proceed as before
	payload := rpc.Arguments{
		{
			Name:     rpc.RPC_DESTINATION_PORT,
			DataType: rpc.DataUint64,
			Value:    uint64(0),
		},
		{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    entry,
		},
	}
	transfer := rpc.Transfer{
		Destination: address,
		Amount:      uint64(1),
		Payload_RPC: payload,
	}
	transferParams := rpc.Transfer_Params{
		Transfers: []rpc.Transfer{transfer},
	}
	result := dero.SendTransfer(transferParams)

	// Redirect to the success page with the form data
	return c.Redirect("/success?result=" + result)
}

func MessagesPage(c *fiber.Ctx) error {
	entries, err := dero.GetAllTransfers()
	if err != nil {
		// Return a 500 Internal Server Error with a meaningful message
		return c.Status(http.StatusInternalServerError).SendString("Error fetching entries: " + err.Error())
	}
	addr, _ := dero.Address()
	// Define data for rendering the template
	data := models.MessagesData{
		Title:    exports.APP_NAME,
		Dev:      addr,
		Messages: entries.Entries,
	}

	tmpl, err := template.ParseFiles("./site/public/messages.html")
	if err != nil {
		// Return a 500 Internal Server Error with a meaningful message
		return c.Status(http.StatusInternalServerError).SendString("Error parsing template: " + err.Error())
	}

	// Execute the template with the provided data
	err = tmpl.Execute(c.Response().BodyWriter(), data)
	if err != nil {
		// Return a 500 Internal Server Error with a meaningful message
		return c.Status(http.StatusInternalServerError).SendString("Error executing template: " + err.Error())
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

// Function to split a string into chunks of a specified size
func chunkString(s string, chunkSize int) []string {
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func OrderPage(c *fiber.Ctx) error {
	addr, _ := dero.Address()
	// Define data for rendering the template
	data := models.HomeData{
		Title: exports.APP_NAME,
		Dev:   addr,
	}

	tmpl, err := template.ParseFiles("./site/public/order.html")
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

func SubmitOrderForm(c *fiber.Ctx) error {
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
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Execute the success template
	err = tmpl.Execute(c.Response().BodyWriter(), data)
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
