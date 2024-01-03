package functions

func echo(username string) string {
	return username
}

func sayHello(username string) string {
	return "Hello, " + Echo(username)
}

func Ping() bool {
	return true
}
