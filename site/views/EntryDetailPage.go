package views

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/deroproject/derohe/rpc"
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/site/models"
)

func EntryDetailPage(c *fiber.Ctx) error {
	// Get the title parameter from the route
	title := c.Params("title")

	// Retrieve the entries from the database or wherever they are stored
	transfers, err := dero.GetOutgoingTransfers()
	if err != nil {
		// Return a 500 Internal Server Error with a meaningful message
		return c.Status(http.StatusInternalServerError).SendString("Error fetching entries: " + err.Error())
	}

	// Define a slice to store transfer comments
	timeEntries := buildTimeEntries(transfers.Entries)

	// Define variables to hold the content of the blog entry
	var blogTitle, blogContent, blogTime string

	// Iterate over timeEntries to create blog content
	for _, timeEntry := range timeEntries {
		// Combine comments for the same time using "__DELIMITER__"
		commentString := strings.Join(timeEntry.Comments, "")

		// Deduce title from comments using the delimiter
		parts := strings.SplitN(commentString, "__DELIMITER__", 2)
		if len(parts) >= 2 && parts[0] == title {

			// Extract title and content
			blogTime = timeEntry.Time.String()
			blogTitle = strings.TrimSpace(parts[0])
			blogContent = strings.TrimSpace(strings.TrimSuffix(parts[1], "__DELIMITER__"))

			break // Stop iteration after finding the matching title
		}
	}

	addr, err := dero.Address()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error retrieving address: " + err.Error())
	}

	// Render the entry detail page template with the entry content
	data := struct {
		App   string
		Dev   string
		Entry models.BlogEntry
	}{
		App: exports.DEVELOPER_NAME,
		Dev: addr, // Assuming addr is defined somewhere in your code
		Entry: models.BlogEntry{
			Time:    blogTime, // Assuming 'time' should be 'entry.Time', update with actual value
			Title:   blogTitle,
			Content: blogContent,
		},
	}

	tmpl, err := template.New("entry_detail.html").Funcs(
		template.FuncMap{
			"URLDecode":       URLDecode,
			"ReplaceNewlines": ReplaceNewlines,
		},
	).ParseFiles("./site/public/entry_detail.html")
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error parsing template: " + err.Error())
	}
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error parsing template: " + err.Error())
	}

	if err := tmpl.Execute(c.Response().BodyWriter(), data); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error executing template: " + err.Error())
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

func buildTimeEntries(entriesData []rpc.Entry) []TimeEntry {
	var transferEntries []TimeEntry
	for _, e := range entriesData {
		if e.DestinationPort != exports.DEST_PORT {
			continue
		}
		// Check if the entry contains comments
		if e.Amount != 0 {
			continue
		}
		if !e.Payload_RPC.Has(rpc.RPC_COMMENT, rpc.DataString) {
			continue
		}
		var existingEntry *TimeEntry
		for i := range transferEntries {
			if transferEntries[i].Time.Equal(e.Time) {
				existingEntry = &transferEntries[i]
				break
			}
		}
		if existingEntry != nil {
			// Preserve string formatting
			commentValue, ok := e.Payload_RPC.Value(rpc.RPC_COMMENT, rpc.DataString).(string)
			if !ok {
				// Handle error, log, or skip this entry
				continue
			}
			existingEntry.Comments = append(existingEntry.Comments, commentValue)
		} else {
			// Preserve string formatting
			commentValue, ok := e.Payload_RPC.Value(rpc.RPC_COMMENT, rpc.DataString).(string)
			if !ok {
				// Handle error, log, or skip this entry
				continue
			}
			newEntry := TimeEntry{
				Time:     e.Time,
				Comments: []string{commentValue},
			}
			transferEntries = append(transferEntries, newEntry)
		}
	}
	return transferEntries
}
