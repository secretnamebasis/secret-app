package handlers

import (
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"go.etcd.io/bbolt"
)

func initialLoad(db *bbolt.DB) error {
	if loaded {
		return nil
	}
	exports.Logs.Info(dero.Echo("Initial Loading of Wallet Entries"))

	transfers, err := dero.GetIncomingTransfers()
	if err != nil {
		exports.Logs.Error(err, "Wallet Failed to Get Entries")
		return err
	}

	for _, e := range transfers.Entries {
		if err := IncomingTransferEntry(e, db); err != nil {
			return err // Exit on error during initial load
		}
	}

	return nil
}
