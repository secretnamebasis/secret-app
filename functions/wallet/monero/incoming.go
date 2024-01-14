package monero

import (
	"time"

	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/local"
	"go.etcd.io/bbolt"
)

var err error

func IncomingTransferEntry(e Entry, db *bbolt.DB) error {
	// Simulating a condition that might lead to an error
	// var err = errors.New("error")

	// // if e.Amount <= 0 {
	// // 	exports.Logs.Error(err, "Monero Transfer amount is less than 0", "txid", e.TxID, "amount", e.Amount, "payID", e.PaymentID)
	// // }

	// // if e.PaymentID == "" {
	// // 	exports.Logs.Error(err, "Monero Transfer has no PaymentID", "txid", e.TxID, "amount", e.Amount, "payID", e.PaymentID)

	// // }

	// // var already_processed bool
	// _, err = isTransactionProcessed(db, "create", e.TxID)
	// if err != nil {
	// 	return err
	// }
	unixTime := int64(e.Timestamp) // Replace with your Unix timestamp

	t := time.Unix(unixTime, 0)
	formattedTime := t.Format("2006-01-02 15:04:05") // Replace the format as needed

	exports.Logs.Info(local.SayEcho("WALLET Monero: "), "time", formattedTime, "txid", e.TxID, "amount", e.Amount, "payment_id", e.PaymentID)

	return nil
}