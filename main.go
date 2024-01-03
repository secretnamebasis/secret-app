package main

import (
	"log"

	"github.com/secretnamebasis/secret-app/src"
)

func main() {
	err := src.RunApp()
	if err != nil {
		log.Fatal(err)
	}
}
