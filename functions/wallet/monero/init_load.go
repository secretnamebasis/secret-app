package monero

import (
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/local"

	"go.etcd.io/bbolt"
)

var loaded bool

func MoneroLoad(db *bbolt.DB) error {
	if loaded {
		return nil
	}
	exports.Logs.Info(local.SayEcho("Initial Loading of Wallet Entries"))

	transfers, err := GetIncomingTransfers()
	if err != nil {
		exports.Logs.Error(err, "Wallet Failed to Get Entries")
		return err
	}

	for _, e := range transfers.Incoming {
		if err := IncomingTransferEntry(e, db); err != nil {
			return err // Exit on error during initial load
		}
	}

	return nil
}
