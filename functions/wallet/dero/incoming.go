package dero

import (
	"errors"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"go.etcd.io/bbolt"
)

var loaded bool

func IncomingTransferEntry(e rpc.Entry, db *bbolt.DB) error {
	// Simulating a condition that might lead to an error
	var err = errors.New("error")

	if e.Amount <= 0 {
		exports.Logs.Error(err, "amount is less than 0", "txid", e.TXID, "dst_port", e.DestinationPort)
		return nil
	}

	if !e.Payload_RPC.Has(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress) {
		// so if they don't have an replyback, that's okay
		// exports.Logs.Error(nil, fmt.Sprintf("user has not give his address so we cannot replyback")) // this is an unexpected situation
		// return nil
	}

	var already_processed bool
	already_processed, err = isTransactionProcessed(db, "create", e.TXID)
	if err != nil {
		return err
	}

	switch {
	case e.Coinbase:
		CoinbaseTransfer(err, e)

	case already_processed:
		AlreadyProcessedTransfer(err, e)

	case !e.Payload_RPC.Has(rpc.RPC_DESTINATION_PORT, rpc.DataUint64):
		NoDstPort(err, e)

	case e.Payload_RPC.Has(rpc.RPC_DESTINATION_PORT, rpc.DataUint64):
		ToBeProcessed(e, db)

	default:
		return nil
	}
	return nil
}
