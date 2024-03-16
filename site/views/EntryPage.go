package views

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/deroproject/derohe/rpc"
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/site/models"
)

func EntryPage(c *fiber.Ctx) error {
	addr, _ := dero.Address()
	// Define data for rendering the template
	data := models.HomeData{
		App:    exports.DEVELOPER_NAME,
		Dev:    addr,
		Height: dero.Height(),
	}

	tmpl, err := template.ParseFiles("./site/public/entry.html")
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
	c.Set(
		fiber.HeaderContentType,
		fiber.MIMETextHTML,
	)

	return nil
}
func SubmitEntryForm(c *fiber.Ctx) error {
	rawTitle := c.FormValue("title")
	fmt.Println("raw: " + rawTitle)

	entry := c.FormValue("entry")
	fmt.Print(entry)

	encodedTitle := url.PathEscape(rawTitle)
	fmt.Println("encoded: " + encodedTitle)

	// Encode entry to preserve newlines
	encodedEntry := url.QueryEscape(entry)
	fmt.Println("encoded entry: " + encodedEntry)

	// Concatenate encoded title and entry
	submission := fmt.Sprintf("%s__DELIMITER__%s__DELIMITER__", encodedTitle, encodedEntry)

	if len(submission) <= exports.CHUNK_SIZE_MAX {
		fmt.Print("single entry \n")
		result, err := processSingleEntry(submission)
		if err != nil {
			return err
		}
		return c.Redirect("/success?result=" + result.TXID)
	}
	fmt.Print("multiple entries \n")
	chunks := chunkString(submission, 100)

	var txs []rpc.Transfer

	for _, text := range chunks {
		tx := prepareTransfer(text)
		txs = append(txs, tx...)
	}

	result, err := dero.MakeBulkTransfer(txs)
	if err != nil {
		return err
	}

	return c.Redirect("/success?result=" + result.TXID)
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

func prepareTransfer(text string) []rpc.Transfer {
	payload := rpc.Arguments{
		{
			Name:     rpc.RPC_DESTINATION_PORT,
			DataType: rpc.DataUint64,
			Value:    exports.DEST_PORT,
		},
		{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    text,
		},
	}

	transfer := rpc.Transfer{
		Destination: exports.DEVELOPER_ADDRESS,
		Amount:      uint64(0),
		Payload_RPC: payload,
	}

	return []rpc.Transfer{transfer}
}

func prepareParams(transfers []rpc.Transfer) rpc.Transfer_Params {

	transferParams := rpc.Transfer_Params{
		Transfers: transfers,
	}

	return transferParams

}

func processSingleEntry(text string) (rpc.Transfer_Result, error) {
	// Send the transfers

	return dero.SendTransfer(prepareParams(prepareTransfer(text)))
}
