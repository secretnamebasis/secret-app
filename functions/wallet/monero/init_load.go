package monero

import (
	"github.com/secretnamebasis/secret-app/exports"
	"go.etcd.io/bbolt"
)

var loaded bool

func Load(db *bbolt.DB) error {
	if loaded {
		return nil
	}

	transfers, err := GetIncomingTransfers()
	if err != nil {
		exports.Logs.Error(err, "Wallet Failed to Get Entries")
		return err
	}

	for _, e := range transfers.In {
		if err := IncomingTransferEntry(e, db); err != nil {
			return err // Exit on error during initial load
		}
	}

	return nil
}
