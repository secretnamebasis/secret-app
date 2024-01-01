package app

import (
	"fmt"
)

var Name = "secret"

func RunApp() error {
	fmt.Println(SayHello(Name))
	return nil
}

func Echo(name string) string {
	return name
}

func SayHello(name string) string {
	return "Hello, " + Echo(name)
}
