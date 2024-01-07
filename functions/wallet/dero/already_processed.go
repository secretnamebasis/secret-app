package dero

import (
	"github.com/deroproject/derohe/rpc"
)

func AlreadyProcessedTransfer(err error, e rpc.Entry) {
	TransferError(err, e, "already processed transfer")
}
