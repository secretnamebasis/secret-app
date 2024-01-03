package src

import (
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions"
	"github.com/secretnamebasis/secret-app/handlers"
)

var (
	db_name string
	sale    []byte
	logger  logr.Logger = logr.Discard()
)

func RunApp() error {

	functions.Logger()

	if functions.Connection() == false {
		err := errors.New("Wallet Connection Failure")
		logger.Error(err, "Error")
		return fmt.Errorf(
			"Failed to establish wallet connection",
		)
	}
	logger.Info(
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

	logger.Info(
		functions.Echo(
			"ID has been created",
		),
	)

	db, err := functions.CreateDB(db_name)

	if err != nil {
		logger.Error(err, err.Error())
	}

	logger.Info(
		functions.Echo(
			"Database has been created",
		),
	)

	// Let's make a bucket
	sale = []byte("SALE")
	functions.CreateBucket(db, sale)

	logger.Info(
		functions.Echo(
			"Sale's list initiated",
		),
	)
	logger.Info(
		functions.Echo(
			"Integrated Address with Expected Arguments: " +
				functions.CreateServiceAddress(
					functions.Address(),
				),
		),
	)

	logger.Info(
		functions.Echo(
			"Integrated Address with Expected Arguments minus Hardcoded Value: " +
				functions.CreateServiceAddressWithoutHardcodedValue(
					functions.Address(),
				),
		),
	)

	handlers.HandleIncomingTransfers(db)

	return nil // Stop the loop and return nil

}
