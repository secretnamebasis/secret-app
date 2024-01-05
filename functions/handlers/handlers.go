package handlers

import (
	"errors"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet"

	"go.etcd.io/bbolt"
)

func IncomingTransfers(db *bbolt.DB) error {
	LoopActivated := false
	exports.Logs.Info(wallet.Echo("Entering For Loop"))

	for {
		transfers, err := wallet.GetTransfers()
		if err != nil {
			exports.Logs.Error(err, "Wallet Failed to Get Entries")
		}
		if !LoopActivated {
			exports.Logs.Info("Wallet Entries are Instantiated")
			LoopActivated = true
		}

		for _, e := range transfers.Entries {
			if err := IncomingTransferEntry(e, db); err != nil {
				return err // Exit inner loop on error
			}
		}
		if exports.Testing {
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(18 * time.Second)
			/*
				it would be better to do a check to see if there was a height change in the wallet
				that way, every height is checked instead of in the likelihood that there is a transfer
				so we need to go back and we need to look at how do we manage the first data entry, that is height change
				once there, we can manage how we handle the heights first, and the transfer data pertinate to the height,
				if there are no changes, then we shouldn't have to process all transfers again.
				any change in height state will trigger the review of transfers at that height;
				if none, no need to handle all the transfers again,
				that's error prone,
				 and then we have to use a database to store "seen" data ,
				 and now we are pounding away on the database when we already have the data
				 we don't need to review data over and over again
				 we just need to see if there is a change, and move on.


			*/
		}

	}
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
