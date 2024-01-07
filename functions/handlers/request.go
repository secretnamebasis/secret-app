package handlers

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/handlers/create"
	"github.com/secretnamebasis/secret-app/functions/handlers/destroy"
	"github.com/secretnamebasis/secret-app/functions/handlers/review"
	"github.com/secretnamebasis/secret-app/functions/handlers/update"
	"go.etcd.io/bbolt"
)

func Request(e rpc.Entry, message string, db *bbolt.DB) {
	// if message != "" {
	// 	logRequestInfo(e, message)
	// }
	switch value := e.Amount; value {
	case uint64(exports.PONG_AMOUNT):
		create.Request(e, message, db)
	case uint64(1):
		review.Request(e, message, db)
	case uint64(2):
		update.Request(e, message, db)
	case uint64(3):
		destroy.Request(e, message, db)
	}

}
