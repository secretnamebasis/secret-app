package views

import (
	"html/template"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/site/models"
)

// TimeEntry represents a single transfer entry with comments
type TimeEntry struct {
	Time     time.Time
	Comments []string
}

func BlogPage(c *fiber.Ctx) error {
	transfers, err := dero.GetOutgoingTransfers()
	if err != nil {
		// Return a 500 Internal Server Error with a meaningful message
		return c.Status(http.StatusInternalServerError).SendString("Error fetching entries: " + err.Error())
	}
	addr, _ := dero.Address()

	// Define a slice to store transfer entries
	var transferEntries []TimeEntry

	// Iterate over transfers to populate the transferEntries slice
	for _, transfer := range transfers.Entries {
		if transfer.DestinationPort != exports.DEST_PORT {
			continue
		}
		// Convert transfer time to string
		transferTime := transfer.Time

		commentValue := transfer.Payload_RPC.Value(rpc.RPC_COMMENT, rpc.DataString)
		if commentValue != nil {
			// Check if transfer time already exists in transferEntries
			var existingEntry *TimeEntry
			for i := range transferEntries {
				if transferEntries[i].Time.Equal(transferTime) {
					existingEntry = &transferEntries[i]
					break
				}
			}

			// If transfer time exists, append comment to existing entry, else create a new entry
			if existingEntry != nil {
				existingEntry.Comments = append(existingEntry.Comments, commentValue.(string))
			} else {
				newEntry := TimeEntry{
					Time:     transferTime,
					Comments: []string{commentValue.(string)},
				}
				transferEntries = append(transferEntries, newEntry)
			}
		}
	}

	// Define data for rendering the template
	var blogEntries []models.BlogEntry

	// Iterate over transferEntries to create blog entries
	for _, entry := range transferEntries {
		// Combine comments for the same time using "__DELIMITER__"
		comment := strings.Join(entry.Comments, "__DELIMITER__")

		// Deduce title from comments using the delimiter
		titleParts := strings.SplitN(comment, "__DELIMITER__", 2)
		if len(titleParts) >= 1 {
			title := titleParts[0] // First part of comment as title
			blogEntry := models.BlogEntry{
				Title: title,
			}
			blogEntries = append(blogEntries, blogEntry)
		}
	}

	// Sort the blog entries based on their titles
	sort.Slice(blogEntries, func(i, j int) bool {
		return blogEntries[i].Title < blogEntries[j].Title
	})

	data := models.EntriesData{
		App:     exports.DEVELOPER_NAME,
		Dev:     addr,
		Entries: blogEntries,
	}

	tmpl, err := template.New("entries.html").Funcs(
		template.FuncMap{
			"split":     SplitString,
			"urldecode": URLDecode,
		},
	).ParseFiles(
		"./site/public/entries.html",
	)
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
