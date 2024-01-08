package code

import (
	"errors"
	"fmt"

	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/crypto"
	"github.com/secretnamebasis/secret-app/functions/database"
	"github.com/secretnamebasis/secret-app/functions/wallet"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/functions/wallet/monero"
	"go.etcd.io/bbolt"
)

var (
	db_name string
	create  []byte

	LoopActivated bool
)

func RunApp() error {

	dero.Logger()

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

	dero_DB := make_db_name(dero.Address())
	exports.Logs.Info(
		dero.Echo(
			"DERO ID Created: " + dero_DB,
		),
	)

	exports.Logs.Info(
		dero.Echo("DERO Address: " +
			dero.Address(),
		),
	)

	exports.Logs.Info(
		dero.Echo("Monero Address: " +
			monero.Address(0),
		),
	)

	exports.Logs.Info(
		dero.Echo("DERO Integrated Address with Expected Arguments: " +
			dero.CreateServiceAddress(
				dero.Address(),
			),
		),
	)

	exports.Logs.Info(
		dero.Echo("DERO Integrated Address with Expected Arguments minus Hardcoded Value: " +
			dero.CreateServiceAddressWithoutHardcodedValue(
				dero.Address(),
			),
		),
	)

	wallet.IncomingTransfers(make_db(dero_DB))

	return nil // Stop the loop and return nil
}

func make_db_name(s string) string {

	// Let's make a database
	db_name = fmt.Sprintf(
		"%s_%s.bbolt.db",
		exports.APP_NAME,
		crypto.Sha1Sum(
			s,
		),
	)
	return db_name
}

func make_db(s string) *bbolt.DB {

	db, err := database.CreateDB(db_name)

	if err != nil {
		exports.Logs.Error(err, err.Error())
	}

	// Let's make a bucket
	create = []byte("create")
	database.CreateBucket(db, create)
	return db
}
