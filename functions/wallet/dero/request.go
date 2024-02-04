package dero

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"

	"go.etcd.io/bbolt"
)

func request(e rpc.Entry, message string, db *bbolt.DB) {
	// if message != "" {
	// 	logRequestInfo(e, message)
	// }
	if !e.Payload_RPC.Has(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress) {
		// so if they don't have an replyback, that's okay
		// exports.Logs.Error(nil, fmt.Sprintf("user has not give his address so we cannot replyback")) // this is an unexpected situation
		// return nil
		// noReplyBackRequest(e, message, db)
	} else {
		switch value := e.Amount; value {
		case uint64(exports.PONG_AMOUNT):
			createRequest(e, message, db)
		case uint64(1):
			reviewRequest(e, message, db)
		case uint64(2):
			updateRequest(e, message, db)
		case uint64(3):
			destroyRequest(e, message, db)
		}
	}
}
