package code

import (
	"errors"
	"fmt"

	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/crypto"
	"github.com/secretnamebasis/secret-app/functions/database"
	"github.com/secretnamebasis/secret-app/functions/logger"
	"github.com/secretnamebasis/secret-app/functions/wallet"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/functions/wallet/monero"
)

var (
	db_name string
	create  []byte

	LoopActivated bool
)

func RunApp() error {

	logger.Logger()

	exports.Logs.Info(
		dero.Echo(
			"Logger has started",
		),
	)

	if dero.Connection() == false {
		err := errors.New("Wallet Connection Failure")
		exports.Logs.Error(err, "Error")
		return fmt.Errorf(
			"Failed to establish wallet connection",
		)
	}

	if monero.Height() <= 0 {
		err := errors.New("Wallet Connection Failure")
		exports.Logs.Error(err, dero.Echo("Error"))
		return fmt.Errorf(
			"Failed to establish wallet connection")
	}

	// Let's make a database
	db_name = fmt.Sprintf(
		"%s_%s.bbolt.db",
		exports.APP_NAME,
		crypto.Sha1Sum(
			dero.Address(),
		),
	)

	exports.Logs.Info(
		dero.Echo(
			"ID Created",
		),
	)

	db, err := database.CreateDB(db_name) // database management might become a serious issue
	/*
		I mean, I made the database because capt demonstrates it, but when I play the thought out...
		the wallet is the freaking database, why the hell do I need to make another one?
		I think that having the ability to interact with the wallet database would be nice...
		but I am too smoothe brained to figure that out right now.
	*/

	if err != nil {
		exports.Logs.Error(err, err.Error())
	}

	// Let's make a bucket
	create = []byte("create")
	database.
		CreateBucket(db, create)

	exports.Logs.Info(
		dero.Echo(
			"Integrated Address with Expected Arguments: " +
				dero.CreateServiceAddress(
					dero.Address(),
				),
		),
	)

	exports.Logs.Info(
		dero.Echo(
			"Integrated Address with Expected Arguments minus Hardcoded Value: " +
				dero.CreateServiceAddressWithoutHardcodedValue(
					dero.Address(),
				),
		),
	)

	wallet.
		ProcessIncomingTransfers(db, LoopActivated)

	return nil // Stop the loop and return nil
}
