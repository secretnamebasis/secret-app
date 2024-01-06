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
var currentHeight int

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

	var currentTransfers *rpc.Get_Transfers_Result

	checkAndProcess := func(transfers *rpc.Get_Transfers_Result) error {
		if currentTransfers != transfers {
			currentTransfers = transfers

			if !*LoopActivated {

				*LoopActivated = true
			}

			for _, e := range transfers.Entries {
				if err := IncomingTransferEntry(e, db); err != nil {
					return err
				}
			}
		}
		return nil
	}

	for {
		height := wallet.Height()

		if currentHeight != height {
			currentHeight = height

			transfers, err := wallet.GetIncomingTransfersByHeight(currentHeight)
			if transfers == nil {
				continue
			}
			if err != nil {
				return err
			}

			if err := checkAndProcess(transfers); err != nil {
				return err
			}
		}

		sleepDuration := 1 * time.Second
		if exports.Testing {
			sleepDuration = 100 * time.Millisecond
		}
		time.Sleep(sleepDuration)
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

func isTransactionProcessed(db *bbolt.DB, bucketName string, TXID string) (bool, error) {
	var alreadyProcessed bool

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket != nil {
			if existing := bucket.Get([]byte(TXID)); existing != nil {
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
	// logTransferError(err, e, "has no dst_port")
}

func handleRequest(e rpc.Entry, message string, db *bbolt.DB) {
	// if message != "" {
	// 	logRequestInfo(e, message)
	// }
	switch value := e.Amount; value {
	case uint64(exports.PONG_AMOUNT):
		handleCreateRequest(e, message, db)
	case uint64(1):
		handleReviewRequest(e, message, db)
	case uint64(2):
		handleUpdateRequest(e, message, db)
	case uint64(3):
		handleDestroyRequest(e, message, db)
	}

}

func handleCreateRequest(e rpc.Entry, message string, db *bbolt.DB) {
	exports.Logs.Info(
		wallet.Echo(message+" request"),
		"txid", e.TXID,
		"amount", e.Amount,
		"dst_port", e.DestinationPort,
		"comment", e.Payload_RPC.Value(rpc.RPC_COMMENT, rpc.DataString),
		"reply_back", e.Payload_RPC.Value(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress),
	)

	where := e.Payload_RPC.Value(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress).(rpc.Address).String()

	reply := rpc.Transfer_Params{
		Transfers: []rpc.Transfer{
			{
				Destination: where,
				Amount:      uint64(1),
				Payload_RPC: rpc.Arguments{
					{
						Name:     rpc.RPC_DESTINATION_PORT,
						DataType: rpc.DataUint64,
						Value:    uint64(0),
					},
					{
						Name:     rpc.RPC_SOURCE_PORT,
						DataType: rpc.DataUint64,
						Value:    exports.DEST_PORT,
					},
					{
						Name:     rpc.RPC_COMMENT,
						DataType: rpc.DataString,
						Value:    exports.Pong,
					},
				},
			},
		},
	}
	result := wallet.SendTransfer(reply)

	// update database
	if result != "" {
		// Perform further actions based on the result
		// ...

		// If processing is successful, update the database
		err := db.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(message))
			return b.Put([]byte(e.TXID), []byte("done"))
		})
		if err != nil {
			exports.Logs.Error(err, "err updating db")
			// Handle the error in updating the database
		} else {
			exports.Logs.Info("ping replied successfully with pong ", "result", result)
			// Log the successful completion
		}
	}
}

func handleReviewRequest(e rpc.Entry, message string, db *bbolt.DB) {
	exports.Logs.Info("Handling review request", "txid", e.TXID)
}

func handleUpdateRequest(e rpc.Entry, message string, db *bbolt.DB) {
	exports.Logs.Info("Handling update request", "txid", e.TXID)
}

func handleDestroyRequest(e rpc.Entry, message string, db *bbolt.DB) {
	exports.Logs.Info("Handling destroy request", "txid", e.TXID)
}

func handleToBeProcessed(e rpc.Entry, db *bbolt.DB) {
	logToBeProcessedInfo(e, "to be processed")

	switch dstPort := e.Payload_RPC.Value(rpc.RPC_DESTINATION_PORT, rpc.DataUint64).(uint64); dstPort {
	case exports.DEST_PORT:
		handleRequest(e, "create", db) //
	case uint64(2):
		handleRequest(e, "review", db)
	case uint64(3):
		handleRequest(e, "update", db)
	case uint64(4):
		handleRequest(e, "destroy", db)
	default:
		handleRequest(e, "", db)
	}
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
