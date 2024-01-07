package handlers

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/functions/logger"
)

func handleCoinbaseTransfer(err error, e rpc.Entry) {
	logger.TransferError(err, e, "coinbase transfer")
}
