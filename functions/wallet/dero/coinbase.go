package dero

import (
	"github.com/deroproject/derohe/rpc"
)

func CoinbaseTransfer(err error, e rpc.Entry) {
	TransferError(err, e, "coinbase transfer")
}
