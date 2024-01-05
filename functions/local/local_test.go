package local_test

import (
	"testing"

	asserts_tests "github.com/secretnamebasis/secret-app/asserts"
	"github.com/secretnamebasis/secret-app/functions/local"
)

func TestSayVar(t *testing.T) {
	username := "secret"
	if local.SayEcho(username) != username {
		t.Errorf("App is not returning strings")
	}
}

func TestSayHelloVar(t *testing.T) {
	given := "secret"
	got := local.SayHello(given)
	want := "Hello, secret"
	asserts_tests.CorrectMessage(t, got, want)
}

func TestPing(t *testing.T) {
	got := local.Ping()
	if got != true {
		t.Errorf("App is not returning pinging")
	}
}
