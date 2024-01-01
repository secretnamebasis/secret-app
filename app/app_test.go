package app_test

import (
	"testing"

	"github.com/secretnamebasis/secret-app/app"
)

func TestRunApp(t *testing.T) {
	if app.RunApp() != nil {
		t.Errorf("App is not running when trying to run app")
	}
}

func TestSayVar(t *testing.T) {
	name := "Alixander"
	if app.Echo(name) != name {
		t.Errorf("App is not returning strings")
	}
}

func TestSayHelloVar(t *testing.T) {
	given := "secret"
	got := app.SayHello(given)
	want := "Hello, secret"
	assertCorrectMessage(t, got, want)
}

func TestPing(t *testing.T) {
	got := Ping()
	if got != true {
		t.Errorf("App is not returning pinging")
	}
}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
