package handlers

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"go.etcd.io/bbolt"
)

func handleRequest(e rpc.Entry, message string, db *bbolt.DB) {
	// if message != "" {
	// 	logRequestInfo(e, message)
	// }
	switch value := e.Amount; value {
	case uint64(exports.PONG_AMOUNT):
		handleCreateRequest(e, message, db)
	case uint64(1):
		handleReviewRequest(e, message, db)
	case uint64(2):
		handleUpdateRequest(e, message, db)
	case uint64(3):
		handleDestroyRequest(e, message, db)
	}

}
