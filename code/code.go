package code

import (
	"errors"
	"fmt"
	"os"

	"github.com/deroproject/derohe/globals"
	"github.com/go-logr/logr"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	db_name       string
	sale          []byte
	Logs          logr.Logger = logr.Discard()
	LoopActivated bool
)

func RunApp() error {

	Logger()

	if functions.Connection() == false {
		err := errors.New("Wallet Connection Failure")
		Logs.Error(err, "Error")
		return fmt.Errorf(
			"Failed to establish wallet connection",
		)
	}
	Logs.Info(
		functions.Echo(
			"Logger has started",
		),
	)

	// Let's make a database
	db_name = fmt.Sprintf(
		"%s_%s.bbolt.db",
		exports.APP_NAME,
		functions.Sha1Sum(functions.Address()),
	)

	Logs.Info(
		functions.Echo(
			"ID has been created",
		),
	)

	db, err := functions.CreateDB(db_name)

	if err != nil {
		Logs.Error(err, err.Error())
	}

	Logs.Info(
		functions.Echo(
			"Database has been created",
		),
	)

	// Let's make a bucket
	sale = []byte("SALE")
	functions.CreateBucket(db, sale)

	Logs.Info(
		functions.Echo(
			"Sale's list initiated",
		),
	)
	Logs.Info(
		functions.Echo(
			"Integrated Address with Expected Arguments: " +
				functions.CreateServiceAddress(
					functions.Address(),
				),
		),
	)

	Logs.Info(
		functions.Echo(
			"Integrated Address with Expected Arguments minus Hardcoded Value: " +
				functions.CreateServiceAddressWithoutHardcodedValue(
					functions.Address(),
				),
		),
	)

	functions.HandleIncomingTransfers(db)

	return nil // Stop the loop and return nil

}

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
	Logs = globals.Logger

	return nil

}
