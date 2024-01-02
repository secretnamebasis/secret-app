package main

import (
	"log"

	"github.com/secretnamebasis/secret-app/code"
)

func main() {
	err := code.RunApp()
	if err != nil {
		log.Fatal(err)
	}
}
