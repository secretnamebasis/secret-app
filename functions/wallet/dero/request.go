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
