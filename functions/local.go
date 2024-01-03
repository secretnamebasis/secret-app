package functions

import (
	"os"

	"github.com/deroproject/derohe/globals"
	"github.com/go-logr/logr"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger logr.Logger = logr.Discard()

func SayEcho(username string) string {
	return username
}

func SayHello(username string) string {
	return "Hello, " + SayEcho(username)
}

func Ping() bool {
	return true
}

func Logger() error {
	// parse arguments and setup logging, print basic information
	globals.Arguments["--debug"] = true
	exename, err := os.Executable()

	globals.InitializeLog(os.Stdout, &lumberjack.Logger{
		Filename:   exename + ".log",
		MaxSize:    100, // megabytes
		MaxBackups: 2,
	})
	logger = globals.Logger

	return err

}
