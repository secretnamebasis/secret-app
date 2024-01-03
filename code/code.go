package code

import (
	"errors"
	"fmt"

	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions"
)

var (
	db_name string
	sale    []byte

	LoopActivated bool
)

func RunApp() error {

	functions.Logger()

	exports.Logs.Info(
		functions.Echo(
			"Logger has started",
		),
	)

	if functions.Connection() == false {
		err := errors.New("Wallet Connection Failure")
		exports.Logs.Error(err, "Error")
		return fmt.Errorf(
			"Failed to establish wallet connection",
		)
	}

	// Let's make a database
	db_name = fmt.Sprintf(
		"%s_%s.bbolt.db",
		exports.APP_NAME,
		functions.Sha1Sum(functions.Address()),
	)

	exports.Logs.Info(
		functions.Echo(
			"ID has been created",
		),
	)

	db, err := functions.CreateDB(db_name)

	if err != nil {
		exports.Logs.Error(err, err.Error())
	}

	exports.Logs.Info(
		functions.Echo(
			"Database has been created",
		),
	)

	// Let's make a bucket
	sale = []byte("SALE")
	functions.CreateBucket(db, sale)

	exports.Logs.Info(
		functions.Echo(
			"Sale's list initiated",
		),
	)
	exports.Logs.Info(
		functions.Echo(
			"Integrated Address with Expected Arguments: " +
				functions.CreateServiceAddress(
					functions.Address(),
				),
		),
	)

	exports.Logs.Info(
		functions.Echo(
			"Integrated Address with Expected Arguments minus Hardcoded Value: " +
				functions.CreateServiceAddressWithoutHardcodedValue(
					functions.Address(),
				),
		),
	)

	err = functions.HandleIncomingTransfers(db)
	if err != nil {
		return err
	}
	return nil // Stop the loop and return nil

}
