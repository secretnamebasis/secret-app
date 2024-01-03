package functions_test

import (
	"testing"

	"github.com/secretnamebasis/secret-app/asserts"
	"github.com/secretnamebasis/secret-app/functions"
)

func TestSayVar(t *testing.T) {
	username := "secret"
	if functions.SayEcho(username) != username {
		t.Errorf("App is not returning strings")
	}
}

func TestSayHelloVar(t *testing.T) {
	given := "secret"
	got := functions.SayHello(given)
	want := "Hello, secret"
	asserts.CorrectMessage(t, got, want)
}

func TestPing(t *testing.T) {
	got := functions.Ping()
	if got != true {
		t.Errorf("App is not returning pinging")
	}
}

func TestLogger(t *testing.T) {
	got := functions.Logger()
	if got != nil {
		t.Errorf("got %q", got)
	}
}
