package handlers

import "github.com/deroproject/derohe/rpc"

func handleCoinbaseTransfer(err error, e rpc.Entry) {
	logTransferError(err, e, "coinbase transfer")
}
