package main

import (
	"log"

	"github.com/secretnamebasis/secret-app/app"
)

func main() {
	err := app.RunApp()
	if err != nil {
		log.Fatal(err)
	}
}
