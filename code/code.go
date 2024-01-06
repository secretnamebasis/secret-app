package code

import (
	"errors"
	"fmt"

	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/crypto"
	"github.com/secretnamebasis/secret-app/functions/database"
	"github.com/secretnamebasis/secret-app/functions/handlers"
	logger "github.com/secretnamebasis/secret-app/functions/logger"
	"github.com/secretnamebasis/secret-app/functions/wallet"
)

var (
	db_name string
	create  []byte

	LoopActivated bool
)

func RunApp() error {

	logger.Logger()

	exports.Logs.Info(
		wallet.Echo(
			"Logger has started",
		),
	)

	if wallet.Connection() == false {
		err := errors.New("Wallet Connection Failure")
		exports.Logs.Error(err, wallet.Echo("Error"))
		return fmt.Errorf(
			"Failed to establish wallet connection",
		)
	}

	// Let's make a database
	db_name = fmt.Sprintf(
		"%s_%s.bbolt.db",
		exports.APP_NAME,
		crypto.Sha1Sum(
			wallet.Address(),
		),
	)

	exports.Logs.Info(
		wallet.Echo(
			"ID Created",
		),
	)

	db, err := database.CreateDB(db_name)

	if err != nil {
		exports.Logs.Error(err, err.Error())
	}

	// Let's make a bucket
	create = []byte("create")
	database.
		CreateBucket(db, create)

	exports.Logs.Info(
		wallet.Echo(
			"Integrated Address with Expected Arguments: " +
				wallet.CreateServiceAddress(
					wallet.Address(),
				),
		),
	)

	exports.Logs.Info(
		wallet.Echo(
			"Integrated Address with Expected Arguments minus Hardcoded Value: " +
				wallet.CreateServiceAddressWithoutHardcodedValue(
					wallet.Address(),
				),
		),
	)

	handlers.
		IncomingTransfers(db)

	return nil // Stop the loop and return nil
}
