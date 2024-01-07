package handlers

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"go.etcd.io/bbolt"
)

func handleDestroyRequest(e rpc.Entry, message string, db *bbolt.DB) {
	exports.Logs.Info("Handling destroy request", "txid", e.TXID)
}
