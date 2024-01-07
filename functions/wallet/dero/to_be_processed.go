package dero

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/handlers"

	"go.etcd.io/bbolt"
)

func ToBeProcessed(e rpc.Entry, db *bbolt.DB) {
	ToBeProcessedInfo(e, "to be processed")

	switch dstPort := e.Payload_RPC.Value(rpc.RPC_DESTINATION_PORT, rpc.DataUint64).(uint64); dstPort {
	case exports.DEST_PORT:
		handlers.Request(e, "create", db) //
	case uint64(2):
		handlers.Request(e, "review", db)
	case uint64(3):
		handlers.Request(e, "update", db)
	case uint64(4):
		handlers.Request(e, "destroy", db)
	default:
		handlers.Request(e, "", db)
	}
}
