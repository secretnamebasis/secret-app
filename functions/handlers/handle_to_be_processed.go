package handlers

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"go.etcd.io/bbolt"
)

func handleToBeProcessed(e rpc.Entry, db *bbolt.DB) {
	logToBeProcessedInfo(e, "to be processed")

	switch dstPort := e.Payload_RPC.Value(rpc.RPC_DESTINATION_PORT, rpc.DataUint64).(uint64); dstPort {
	case exports.DEST_PORT:
		handleRequest(e, "create", db) //
	case uint64(2):
		handleRequest(e, "review", db)
	case uint64(3):
		handleRequest(e, "update", db)
	case uint64(4):
		handleRequest(e, "destroy", db)
	default:
		handleRequest(e, "", db)
	}
}
