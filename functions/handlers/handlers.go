package handlers

import (
	"errors"
	"fmt"

	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet"

	"go.etcd.io/bbolt"
)

func HandleIncomingTransfers(db *bbolt.DB) error {
	LoopActivated := false
	exports.Logs.Info("Entering For Loop")
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
			// Simulating a condition that might lead to an error
			if e.Amount <= 0 {
				return errors.New("invalid transaction amount")
			}
			if e.Coinbase || !e.Incoming { // skip coinbase or outgoing, self generated transactions
				continue
			}

			// check whether the entry has been processed before, if yes skip it
			var already_processed bool
			db.View(func(tx *bbolt.Tx) error {
				if b := tx.Bucket([]byte("SALE")); b != nil {
					if ok := b.Get([]byte(e.TXID)); ok != nil { // if existing in bucket
						already_processed = true
					}
				}
				return nil
			})

			if already_processed { // if already processed skip it
				continue
			}

			// check whether this service should handle the transfer
			if !e.Payload_RPC.Has(rpc.RPC_DESTINATION_PORT, rpc.DataUint64) ||
				exports.DEST_PORT != e.Payload_RPC.Value(rpc.RPC_DESTINATION_PORT, rpc.DataUint64).(uint64) { // this service is expecting value to be specfic
				continue

			}

			exports.Logs.V(1).Info("to be processed", "txid", e.TXID)
			if exports.Expected_arguments.Has(rpc.RPC_VALUE_TRANSFER, rpc.DataUint64) { // this service is expecting value to be specfic
				value_expected := exports.Expected_arguments.Value(rpc.RPC_VALUE_TRANSFER, rpc.DataUint64).(uint64)
				if e.Amount != value_expected { // TODO we should mark it as faulty
					exports.Logs.Error(nil, fmt.Sprintf("user transferred %d, we were expecting %d. so we will not do anything", e.Amount, value_expected)) // this is an unexpected situation
					continue
				}

				if !e.Payload_RPC.Has(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress) {
					exports.Logs.Error(nil, fmt.Sprintf("user has not give his address so we cannot replyback")) // this is an unexpected situation
					continue
				}

				destination_expected := e.Payload_RPC.Value(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress).(rpc.Address).String()
				addr, err := rpc.NewAddress(destination_expected)
				if err != nil {
					exports.Logs.Error(err, "err while while parsing incoming addr")
					continue
				}
				addr.Mainnet = true // convert addresses to testnet form, by default it's expected to be mainnnet
				destination_expected = addr.String()

				exports.Logs.V(1).Info("tx should be replied", "txid", e.TXID, "replyback_address", destination_expected)

				//destination_expected := e.Sender

				// value received is what we are expecting, so time for response
				exports.Response[0].Value = e.SourcePort // source port now becomes destination port, similar to TCP
				exports.Response[2].Value = fmt.Sprintf("%s. You sent %s at height %d", exports.Pong, walletapi.FormatMoney(e.Amount), e.Height)

				//_, err :=  response.CheckPack(transaction.PAYLOAD0_LIMIT)) //  we only have 144 bytes for RPC

				// sender of ping now becomes destination
				var result rpc.Transfer_Result
				tparams := rpc.Transfer_Params{Transfers: []rpc.Transfer{{Destination: destination_expected, Amount: uint64(1), Payload_RPC: exports.Response}}}
				err = exports.RpcClient.CallFor(&result, "Transfer", tparams)
				if err != nil {
					exports.Logs.Error(err, "err while transfer")
					continue
				}

				err = db.Update(func(tx *bbolt.Tx) error {
					b := tx.Bucket([]byte("SALE"))
					return b.Put([]byte(e.TXID), []byte("done"))
				})
				if err != nil {
					exports.Logs.Error(err, "err updating db")
				} else {
					exports.Logs.Info("ping replied successfully with pong ", "result", result)
				}
				if exports.Testing == true {
					return nil
				}
			}
		}

	}
}
