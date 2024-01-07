package dero

import (
	"errors"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"go.etcd.io/bbolt"
)

var loaded bool

func IncomingTransfers(db *bbolt.DB) error {
	LoopActivated := false
	exports.Logs.Info(Echo("Entering For Loop"))

	if err := Load(db); err != nil {
		return err // Exit on error during initial load
	}

	return IncomingTransfers(db, &LoopActivated)
}

func IncomingTransferEntry(e rpc.Entry, db *bbolt.DB) error {
	// Simulating a condition that might lead to an error
	var err = errors.New("error")

	if e.Amount <= 0 {
		exports.Logs.Error(err, "amount is less than 0", "txid", e.TXID, "dst_port", e.DestinationPort)
	}

	if !e.Payload_RPC.Has(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress) {
		// exports.Logs.Error(nil, fmt.Sprintf("user has not give his address so we cannot replyback")) // this is an unexpected situation
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

	case e.Payload_RPC.Has(rpc.RPC_DESTINATION_PORT, rpc.DataUint64) && exports.Expected_arguments.Has(rpc.RPC_VALUE_TRANSFER, rpc.DataUint64):
		ToBeProcessed(e, db)

	default:
		return nil
	}
	return nil
}
