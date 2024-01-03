package local

func SayEcho(username string) string {
	return username
}

func SayHello(username string) string {
	return "Hello, " + SayEcho(username)
}

func Ping() bool {
	return true
}
