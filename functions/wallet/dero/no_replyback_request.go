package dero

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/local"
	"github.com/secretnamebasis/secret-app/functions/wallet/monero"
	"go.etcd.io/bbolt"
)

/*
The user is going to be sending DERO to the exports.DEST_PORT
The user is simply sending DERO with an XMR address
We have already processed for DEST_PORT
We have to determine first, how much is it?
Then we have to determine if the entry has a comment
The reason we do this is because we need an XMR address
*/
func noReplyBackRequest(e rpc.Entry, message string, db *bbolt.DB) error {
	// Simulating a condition that might lead to an error
	var err = errors.New("error")

	if e.Amount <= 0 {
		exports.Logs.Error(err, "amount is less than 0", "txid", e.TXID, "dst_port", e.DestinationPort)
		return nil
	}

	// if !e.Payload_RPC.Has(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress) {
	// exports.Logs.Error(nil, fmt.Sprintf("user has not give his address so we cannot replyback")) // this is an unexpected situation
	// return nil
	// this isn't mandatory, but it would be helpful
	// }
	// address := e.Payload_RPC.Value(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress).(rpc.Address).String()
	// if address == "" {
	// 	// exports.Logs.Error(nil, fmt.Sprintf("user has given address as empty string so we cannot replyback")) // this is an unexpected situation
	// 	// again, would be helpful
	// 	return nil
	// }

	if !e.Payload_RPC.Has(rpc.RPC_COMMENT, rpc.DataString) {
		exports.Logs.Error(err, "entry does not contain comment", "txid", e.TXID, "dst_port", e.DestinationPort)
		return nil
	}

	comment := e.Payload_RPC.Value(rpc.RPC_COMMENT, rpc.DataString).(string)
	if comment == "" {
		// exports.Logs.Error(err, "entry cannot contain and empty string\"\" ", "txid", e.TXID, "dst_port", e.DestinationPort)
		return nil
	}

	// exports.Logs.Info("entry does contains comment", "txid", e.TXID, "comment", comment)
	isValidXMRAddress := monero.ValidateAddress(comment)

	switch isValidXMRAddress {
	case true:
		exports.Logs.Info(Echo("DERO: Found XMR address"), "time", e.Time, "txid", e.TXID, "amount", e.Amount, "XMR ADDR", comment)
		amount := getAmount(e)

		given := monero.Transfer_Params{
			Destinations: []monero.Transfer{
				{
					Amount:  amount, // Amount in atomic units
					Address: comment,
				},
			},
		}

		got := monero.SendTransfer(given)

		if got.Result.TxHash == "" {

			exports.Logs.Info(local.SayEcho("WALLET Write to database: "), "monero txid", got)
		}
		// 		}

		// 		message := "sale"
		// 		updateDatabaseOnSuccess(db, message, e)
		exports.Logs.Info(local.SayEcho("WALLET Write to database: "), "monero txid", got.Result.TxHash)
	case false:
		// exports.Logs.Info(Echo("DERO: Found Comment"), "time", e.Time, "txid", e.TXID, "amount", e.Amount, "Comment", comment)
		// return err
	}
	return nil
}

// exports.Logs.Info(local.SayEcho("WALLET Monero: "), "time", formattedTime, "txid", e.TxID, "amount", e.Amount, "payment_id", e.PaymentID)

// 	// var already_processed bool
// 	// already_processed, err = isTransactionProcessed(db, "create", e.TXID)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// if already_processed != false {
// 	// 	return err
// 	// }

// 	var hasID bool
// 	hasID, err = hasPaymentID(db, "contacts", e)
// 	if err != nil {
// 		return err
// 	}
//

func getAmount(e rpc.Entry) uint64 {
	if e.Amount == 0 {
		return uint64(0)
	}

	deroUsdtURL := "https://tradeogre.com/api/v1/ticker/dero-usdt"
	getUsdtDeroAsk, _ := getTickerQuote(deroUsdtURL)
	exports.Logs.Info(local.SayEcho("quoting"), "amount", getUsdtDeroAsk)
	amountUsdt := float64(e.Amount) * getUsdtDeroAsk * 1e-5

	xmrUsdtURL := "https://tradeogre.com/api/v1/ticker/xmr-usdt"
	getUsdtXmrQuote, _ := getTickerQuote(xmrUsdtURL)
	exports.Logs.Info(local.SayEcho("quoting"), "amount", getUsdtXmrQuote)

	amountXMR := amountUsdt / getUsdtXmrQuote
	exports.Logs.Info(local.SayEcho("converting"), "amount", amountXMR)

	// Convert amountDero to atomic units
	atomicUnits := amountXMR * 1e12 // 1 DERO = 100000 atomic units
	// exports.Logs.Info(local.SayEcho("converting"), "amount", atomicUnits)

	tradeMonero := atomicUnits * 0.99

	exports.Logs.Info(local.SayEcho("sending atomic units"), "amount", uint64(tradeMonero))
	return uint64(tradeMonero)
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

// func updateDatabaseOnSuccess(db *bbolt.DB, message string, e Entry) {
// 	err := db.Update(func(tx *bbolt.Tx) error {
// 		b := tx.Bucket([]byte(message))
// 		return b.Put([]byte(e.TxID), []byte("done"))
// 	})
// 	if err != nil {
// 		exports.Logs.Error(err, "err updating db")
// 	}
// }
