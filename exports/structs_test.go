package exports_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/secretnamebasis/secret-app/functions/requests"
)

func TestRoundTrip(t *testing.T) {
	// Given
	givenUsername := "testuser"
	givenPassword := "testpassword"

	// Mock HTTP server for testing
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, Password, ok := r.BasicAuth()
		if !ok {
			t.Error("Basic authentication not set")
			return
		}

		// Got
		gotUsername := username
		gotPassword := Password

		// Want
		wantUsername := givenUsername
		wantPassword := givenPassword

		// Check if Got matches Want
		if gotUsername != wantUsername || gotPassword != wantPassword {
			t.Errorf("Unexpected basic auth headers. Got username: %s, Want username: %s, Got Password: %s,  Want Password: %s",
				gotUsername, wantUsername, gotPassword, wantPassword)
		}
	}))
	defer ts.Close()

	requests.PerformRequestWithBasicAuth(givenUsername, givenPassword, ts.URL)

}
