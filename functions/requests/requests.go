package requests

import (
	"net/http"

	"github.com/secretnamebasis/secret-app/exports"
)

func PerformRequestWithBasicAuth(givenUsername, givenPassword, url string) error {
	// Create a TransportWithBasicAuth instance
	authTransport := &exports.TransportWithBasicAuth{
		Username: givenUsername,
		Password: givenPassword,
		Base:     http.DefaultTransport,
	}

	// Create a request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Execute RoundTrip method
	_, err = authTransport.RoundTrip(req)
	if err != nil {
		return err
	}

	return nil
}
