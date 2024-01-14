package monero

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/local"

	"go.etcd.io/bbolt"
)

var err error

func IncomingTransferEntry(e Entry, db *bbolt.DB) error {

	if e.Amount <= 0 {
		return nil
	}

	unixTime := int64(e.Timestamp) // Replace with your Unix timestamp
	t := time.Unix(unixTime, 0)
	formattedTime := t.Format("2006-01-02 15:04:05") // Replace the format as needed

	if e.PaymentID == "0000000000000000" {
		return nil
	}

	var already_processed bool
	already_processed, err = isTransactionProcessed(db, "sale", e)
	if err != nil {
		return err
	}
	if already_processed != false {
		return err
	}
	var hasID bool
	hasID, err = hasPaymentID(db, "contacts", e)
	if err != nil {
		return err
	}
	switch hasID {
	case true:
		exports.Logs.Info(local.SayEcho("WALLET Monero: Found Transfer"), "time", formattedTime, "txid", e.TxID, "amount", e.Amount, "payment_id", e.PaymentID)
		reply, _ := createTransfer(e, db)
		if err != nil {
			return err
		}
		var transfers rpc.Transfer_Result
		err = exports.RpcClient.CallFor(
			&transfers,
			"Transfer",
			reply,
		)

		if transfers.TXID == "" {
			err := errors.New("Empty TXID")
			exports.Logs.Error(err, "TXID is \"\" string")
		}

		message := "sale"
		updateDatabaseOnSuccess(db, message, e)
		exports.Logs.Info(local.SayEcho("WALLET Write to database: "), "time", formattedTime, "dero txid", transfers.TXID, "monero txid", e.TxID)
	case false:
		// exports.Logs.Error(err, "Error Message", "time", formattedTime, "txid", e.TxID, "amount", e.Amount, "payment_id", e.PaymentID)
		return err
	}

	// exports.Logs.Info(local.SayEcho("WALLET Monero: "), "time", formattedTime, "txid", e.TxID, "amount", e.Amount, "payment_id", e.PaymentID)

	return nil
}

func getAddressFromEntry(e Entry, db *bbolt.DB) string {

	var byteAddress []byte

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("contacts"))
		if bucket != nil {
			byteAddress = bucket.Get([]byte(e.PaymentID))
		}
		return nil
	})
	if err != nil {
		return err.Error()
	}
	address, err := rpc.NewAddress(fmt.Sprintf("%s", byteAddress))
	if err != nil {
		return ""
	}
	exports.Logs.Info(local.SayEcho("searching"), "address", address.String())

	return address.String()
}

func getAmount(e Entry) uint64 {
	if e.Amount == 0 {
		return uint64(0)
	}

	xmrUsdtURL := "https://tradeogre.com/api/v1/ticker/xmr-usdt"
	deroUsdtURL := "https://tradeogre.com/api/v1/ticker/dero-usdt"

	getUsdtXmrQuote, _ := getTickerQuote(xmrUsdtURL)
	exports.Logs.Info(local.SayEcho("quoting"), "amount", getUsdtXmrQuote)

	amountUsdt := float64(e.Amount) * getUsdtXmrQuote * 1e-12

	getUsdtDeroAsk, _ := getTickerQuote(deroUsdtURL)
	amountDero := amountUsdt / getUsdtDeroAsk
	exports.Logs.Info(local.SayEcho("converting"), "amount", amountDero)

	// Convert amountDero to atomic units
	atomicUnits := amountDero * 1e5 // 1 DERO = 100000 atomic units
	tradeDero := atomicUnits * 0.99

	exports.Logs.Info(local.SayEcho("sending atomic units"), "amount", uint64(tradeDero))
	return uint64(tradeDero)
}

func createTransfer(e Entry, db *bbolt.DB) (rpc.Transfer_Params, error) {
	where := getAddressFromEntry(e, db)
	if where == "invalid checksum" {
		return rpc.Transfer_Params{}, errors.New("Invalid checksum")
	}

	if where == "" {
		exports.Logs.Error(err, "address is an empty string")
	}

	exports.Logs.Info(local.SayEcho("getting"), "where", where)

	amount := getAmount(e)
	message := "Thank you for using secret-swap"

	return rpc.Transfer_Params{
		Transfers: []rpc.Transfer{
			{
				Destination: where,
				Amount:      amount,
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
						Value:    message,
					},
				},
			},
		},
	}, nil
}

func updateDatabaseOnSuccess(db *bbolt.DB, message string, e Entry) {
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(message))
		return b.Put([]byte(e.TxID), []byte("done"))
	})
	if err != nil {
		exports.Logs.Error(err, "err updating db")
	}
}

func getTickerQuote(url string) (float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		exports.Logs.Error(err, "error in the http request")
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		exports.Logs.Error(err, "error reading response body")
		return 0, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		exports.Logs.Error(err, "error unmarshalling data")
		return 0, err
	}

	askString, ok := data["ask"].(string)
	if !ok {
		exports.Logs.Error(fmt.Errorf("ask value is not a string"), "error checking data type")
		return 0, fmt.Errorf("Error extracting ask value from JSON. JSON structure: %v", data)
	}

	quote, err := strconv.ParseFloat(askString, 64)
	if err != nil {
		exports.Logs.Error(err, "error converting ask value to float64")
		return 0, err
	}

	return quote, nil
}
