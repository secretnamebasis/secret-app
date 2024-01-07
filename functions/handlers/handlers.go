package handlers

import (
	"errors"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"

	"go.etcd.io/bbolt"
)

var loaded bool
var currentHeight int

func IncomingTransfers(db *bbolt.DB) error {
	LoopActivated := false
	exports.Logs.Info(dero.Echo("Entering For Loop"))

	if err := initialLoad(db); err != nil {
		return err // Exit on error during initial load
	}

	return processIncomingTransfers(db, &LoopActivated)
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
		handleCoinbaseTransfer(err, e)

	case already_processed:
		handleAlreadyProcessedTransfer(err, e)

	case !e.Payload_RPC.Has(rpc.RPC_DESTINATION_PORT, rpc.DataUint64):
		handleNoDstPort(err, e)

	case e.Payload_RPC.Has(rpc.RPC_DESTINATION_PORT, rpc.DataUint64) && exports.Expected_arguments.Has(rpc.RPC_VALUE_TRANSFER, rpc.DataUint64):
		handleToBeProcessed(e, db)

	default:
		return nil
	}
	return nil
}
