package exports

import "net/http"

type RoundTripper interface {
	RoundTrip(req *http.Request) (*http.Response, error)
}

type TransportWithBasicAuth struct {
	Username string
	Password string
	Base     http.RoundTripper
}

// RoundTrip implements the RoundTripper interface
func (t *TransportWithBasicAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return t.Base.RoundTrip(req)
}
