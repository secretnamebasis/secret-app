package logger

import (
	"os"

	"github.com/deroproject/derohe/globals"
	"github.com/secretnamebasis/secret-app/exports"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Logger() error {
	// parse arguments and setup logging, print basic information
	globals.Arguments["--debug"] = true
	exename, err := os.Executable()
	if err != nil {
		return err
	}

	globals.InitializeLog(os.Stdout, &lumberjack.Logger{
		Filename:   exename + ".log",
		MaxSize:    100, // megabytes
		MaxBackups: 2,
	})
	exports.Logs = globals.Logger

	return err

}
