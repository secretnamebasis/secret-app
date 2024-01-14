package dero

import (
	"github.com/deroproject/derohe/rpc"
	"go.etcd.io/bbolt"
)

func updateRequest(e rpc.Entry, message string, db *bbolt.DB) {
	// exports.Logs.Info(Echo("DERO Handling update request"), "txid", e.TXID)
}
