package handlers

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"go.etcd.io/bbolt"
)

func handleUpdateRequest(e rpc.Entry, message string, db *bbolt.DB) {
	exports.Logs.Info("Handling update request", "txid", e.TXID)
}
