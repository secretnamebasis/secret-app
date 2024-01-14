package dero

import (
	"github.com/deroproject/derohe/rpc"
	"go.etcd.io/bbolt"
)

func reviewRequest(e rpc.Entry, message string, db *bbolt.DB) {
	// exports.Logs.Info(Echo("DERO Handling review request"), "txid", e.TXID)
}
