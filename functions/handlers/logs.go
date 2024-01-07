package handlers

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
)

func logTransferError(err error, e rpc.Entry, errorMessage string) {
	msg := dero.Echo(errorMessage)
	exports.Logs.Error(err, msg, "txid", e.TXID, "dst_port", e.DestinationPort, "amount", e.Amount)
}

func logToBeProcessedInfo(e rpc.Entry, message string) {
	msg := dero.Echo(message)
	exports.Logs.V(1).Info(msg, "txid", e.TXID, "dst_port", e.DestinationPort, "amount", e.Amount)
}

func logRequestInfo(e rpc.Entry, message string) {
	msg := dero.Echo(message)
	exports.Logs.Info(msg, "txid", e.TXID, "dst_port", e.DestinationPort, "amount", e.Amount)
}
