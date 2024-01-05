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

func processIncomingTransfers(db *bbolt.DB, LoopActivated *bool) error {
	var currentHeight int // Store the current wallet height
	var currentTransfers *rpc.Get_Transfers_Result
	for {
		height := wallet.Height() // Get the current wallet height

		if currentHeight != height {
			currentHeight = height // Update the current height

			transfers, err := wallet.GetIncomingTransfersByHeight(currentHeight)

			if transfers == nil {
				continue
			}

			if err != nil {
				return err
			}

			if currentTransfers != transfers {
				currentTransfers = transfers

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
		}

		if exports.Testing {
			time.Sleep(100 * time.Millisecond)
		} else {
			time.Sleep(900 * time.Millisecond)
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

func isTransactionProcessed(db *bbolt.DB, bucketName string, txID string) (bool, error) {
	var alreadyProcessed bool

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket != nil {
			if existing := bucket.Get([]byte(txID)); existing != nil {
				alreadyProcessed = true
			}
		}
		return nil
	})

	if err != nil {
		return false, err // Return false and the encountered error
	}

	return alreadyProcessed, nil
}
func handleCoinbaseTransfer(err error, e rpc.Entry) {
	logTransferError(err, e, "coinbase transfer")
}

func handleAlreadyProcessedTransfer(err error, e rpc.Entry) {
	logTransferError(err, e, "already processed transfer")
}

func handleNoDstPort(err error, e rpc.Entry) {
	logTransferError(err, e, "has no dst_port")
}

func handleRequest(e rpc.Entry, message string) {
	if message != "" {
		logRequestInfo(e, message)
		return
	}
	switch value := e.Payload_RPC.Value(rpc.RPC_VALUE_TRANSFER, rpc.DataUint64); value {
	case uint64(exports.PONG_AMOUNT):
		handleCreateRequest(e)
	case uint64(1):
		handleReviewRequest(e)
	case uint64(2):
		handleUpdateRequest(e)
	case uint64(3):
		handleDestroyRequest(e)
	}

}

func handleCreateRequest(e rpc.Entry) {
	exports.Logs.Info("Handling create request", "txid", e.TXID)
}

func handleReviewRequest(e rpc.Entry) {
	exports.Logs.Info("Handling review request", "txid", e.TXID)
}

func handleUpdateRequest(e rpc.Entry) {
	exports.Logs.Info("Handling update request", "txid", e.TXID)
}

func handleDestroyRequest(e rpc.Entry) {
	exports.Logs.Info("Handling destroy request", "txid", e.TXID)
}

func handleToBeProcessed(e rpc.Entry) {
	logToBeProcessedInfo(e, "to be processed")

	switch dstPort := e.Payload_RPC.Value(rpc.RPC_DESTINATION_PORT, rpc.DataUint64).(uint64); dstPort {
	case exports.DEST_PORT:
		handleRequest(e, "create")
	case uint64(2):
		handleRequest(e, "review")
	case uint64(3):
		handleRequest(e, "update")
	case uint64(4):
		handleRequest(e, "destroy")
	default:
		handleRequest(e, "")
	}
}
func IncomingTransferEntry(e rpc.Entry, db *bbolt.DB) error {
	// Simulating a condition that might lead to an error
	var err = errors.New("error")

	if e.Amount <= 0 {
		exports.Logs.Error(err, "amount is less than 0", "txid", e.TXID, "dst_port", e.DestinationPort)
	}

	var already_processed bool
	already_processed, err = isTransactionProcessed(db, "SALE", e.TXID)
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
		handleToBeProcessed(e)

	default:
		return nil
	}
	return nil
}
