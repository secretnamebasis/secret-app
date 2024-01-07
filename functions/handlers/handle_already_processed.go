package handlers

import "github.com/deroproject/derohe/rpc"

func handleAlreadyProcessedTransfer(err error, e rpc.Entry) {
	logTransferError(err, e, "already processed transfer")
}
