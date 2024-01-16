package main

import (
	"log"

	"github.com/secretnamebasis/secret-app/code"
)

func main() {
	err := code.Run()
	if err != nil {
		log.Fatal(err)
	}
}
