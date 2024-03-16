// views/views.go

package views

import (
	"html/template"
	"net/url"
	"strings"
)

// SplitString splits a string based on the delimiter and returns a slice of strings
func SplitString(str, delim string) []string {
	return strings.Split(str, delim)
}

// URLDecode decodes a URL-encoded string.
func URLDecode(encodedString string) (string, error) {
	decodedString, err := url.QueryUnescape(encodedString)
	if err != nil {
		return "", err
	}
	return decodedString, nil
}

func ReplaceNewlines(content string) template.HTML {
	decodedContent, err := url.QueryUnescape(content)
	if err != nil {
		// Handle error
		return template.HTML(content) // Return original content if decoding fails
	}
	return template.HTML(strings.ReplaceAll(decodedContent, "\r\n", "<br>"))
}
