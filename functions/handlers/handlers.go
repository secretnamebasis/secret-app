package handlers

import (
	"errors"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet"

	"go.etcd.io/bbolt"
)

var loaded bool

func initialLoad(db *bbolt.DB) error {
	if loaded {
		return nil
	}
	exports.Logs.Info(wallet.Echo("Initial Loading of Wallet Entries"))

	transfers, err := wallet.GetIncomingTransfers()
	if err != nil {
		exports.Logs.Error(err, "Wallet Failed to Get Entries")
		return err
	}

	for _, e := range transfers.Entries {
		if err := IncomingTransferEntry(e, db); err != nil {
			return err // Exit on error during initial load
		}
	}
	loaded = true
	return nil
}

func processIncomingTransfers(db *bbolt.DB, LoopActivated *bool) error {
	var currentHeight int // Store the current wallet height

	for {
		height := wallet.Height() // Get the current wallet height

		if currentHeight != height {
			currentHeight = height // Update the current height

			transfers, err := wallet.GetIncomingTransfersByHeight(currentHeight)
			if err != nil {
				exports.Logs.Error(err, "Wallet Failed to Get Entries")
				continue // Continue to the next iteration on error
			}

			if !*LoopActivated {
				exports.Logs.Info("Wallet Entries are Instantiated")
				*LoopActivated = true
			}

			for _, e := range transfers.Entries {
				if err := IncomingTransferEntry(e, db); err != nil {
					return err // Exit inner loop on error
				}
			}
		}

		if exports.Testing {
			time.Sleep(1 * time.Second) // Adjust the delay for testing mode
		} else {
			time.Sleep(18 * time.Second) // Normal delay for production
		}
	}
}

func IncomingTransfers(db *bbolt.DB) error {
	LoopActivated := false
	exports.Logs.Info(wallet.Echo("Entering For Loop"))

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

	var already_processed bool
	db.View(
		func(tx *bbolt.Tx) error {
			if b := tx.Bucket([]byte("SALE")); b != nil {
				if ok := b.Get([]byte(e.TXID)); ok != nil { // if existing in bucket
					already_processed = true
				}
			}
			return nil
		},
	)

	s := ""

	switch {
	case e.Coinbase:
		logTransferError(err, e, "coinbase transfer")

	case already_processed:
		logTransferError(err, e, "already processed transfer")

	case !e.Payload_RPC.Has(rpc.RPC_DESTINATION_PORT, rpc.DataUint64):
		logTransferError(err, e, "has no dst_port")

	case e.Payload_RPC.Has(rpc.RPC_DESTINATION_PORT, rpc.DataUint64) && exports.Expected_arguments.Has(rpc.RPC_VALUE_TRANSFER, rpc.DataUint64):
		logToBeProcessedInfo(e, "to be processed")

		create := exports.DEST_PORT
		review := uint64(2)
		update := uint64(3)
		destroy := uint64(4)

		// Nested switch within the above case block
		switch dstPort := e.Payload_RPC.Value(rpc.RPC_DESTINATION_PORT, rpc.DataUint64).(uint64); {
		case dstPort == create:
			s = "has create request"
			logRequestInfo(e, s)

		case dstPort == review:
			s = "has review request"
			logRequestInfo(e, s)

		case dstPort == update:
			s = "has update request"
			logRequestInfo(e, s)

		case dstPort == destroy:
			s = "has destroy request"
			logRequestInfo(e, s)

		default:
			logRequestInfo(e, s)

		}
	default:
		return nil
	}
	return nil
}

func logTransferError(err error, e rpc.Entry, errorMessage string) {
	msg := wallet.Echo(errorMessage)
	exports.Logs.Error(err, msg, "txid", e.TXID, "dst_port", e.DestinationPort)
}

func logToBeProcessedInfo(e rpc.Entry, message string) {
	msg := wallet.Echo(message)
	exports.Logs.V(1).Info(msg, "txid", e.TXID, "dst_port", e.DestinationPort)
}

func logRequestInfo(e rpc.Entry, message string) {
	msg := wallet.Echo(message)
	exports.Logs.Info(msg, "txid", e.TXID, "dst_port", e.DestinationPort)
}
