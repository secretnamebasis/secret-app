package dero

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
)

func TransferError(err error, e rpc.Entry, errorMessage string) {
	msg := Echo(errorMessage)
	exports.Logs.Error(err, msg, "txid", e.TXID, "dst_port", e.DestinationPort, "amount", e.Amount)
}

func ToBeProcessedInfo(e rpc.Entry, message string) {
	msg := Echo(message)
	exports.Logs.V(1).Info(msg, "txid", e.TXID, "dst_port", e.DestinationPort, "amount", e.Amount)
}

func RequestInfo(e rpc.Entry, message string) {
	msg := Echo(message)
	exports.Logs.Info(msg, "txid", e.TXID, "dst_port", e.DestinationPort, "amount", e.Amount)
}
