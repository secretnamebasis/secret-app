package dero

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"

	"go.etcd.io/bbolt"
)

func ToBeProcessed(e rpc.Entry, db *bbolt.DB) {
	ToBeProcessedInfo(e, "DERO to be processed")

	switch dstPort := e.Payload_RPC.Value(rpc.RPC_DESTINATION_PORT, rpc.DataUint64).(uint64); dstPort {
	case exports.DEST_PORT:
		request(e, "create", db)
	case uint64(2):
		request(e, "review", db)
	case uint64(3):
		request(e, "update", db)
	case uint64(4):
		request(e, "destroy", db)
	default:
		request(e, "", db)
	}
}
