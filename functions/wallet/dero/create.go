package dero

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"

	"github.com/secretnamebasis/secret-app/functions/wallet/monero"
	"go.etcd.io/bbolt"
)

var response monero.IntegratedAddressResponse

func createRequest(e rpc.Entry, message string, db *bbolt.DB) {
	RequestInfo(e, message+" request")

	reply := createTransfer(e)
	result := SendTransfer(reply)

	if result != "" {
		updateDatabaseOnSuccess(db, message, e)
		exports.Logs.Info(Echo("ping replied successfully with pong "), "result", result, "payment_id", response.PaymentID)
		message = "contacts"
		updateContactsOnSuccess(db, message, e, response)
		exports.Logs.Info(Echo("reply back address paired with paymentID "), "address", getAddressFromEntry(e), "payment_id", response.PaymentID)
		// Log the successful completion
	}
}

func createTransfer(e rpc.Entry) rpc.Transfer_Params {
	where := getAddressFromEntry(e)
	request, _ := monero.MakeIntegratedAddress()
	response.IntegratedAddress = request["integrated_address"]
	response.PaymentID = request["payment_id"]
	comment := getCommentFromEntry(response.IntegratedAddress)

	return rpc.Transfer_Params{
		Transfers: []rpc.Transfer{
			{
				Destination: where,
				Amount:      uint64(1),
				Payload_RPC: prepareTransferPayload(comment),
			},
		},
	}
}

func getAddressFromEntry(e rpc.Entry) string {
	return e.Payload_RPC.Value(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress).(rpc.Address).String()
}

func getCommentFromEntry(integratedAddress string) rpc.Arguments {
	return rpc.Arguments{
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
			Value:    integratedAddress,
		},
	}
}

func prepareTransferPayload(comment rpc.Arguments) rpc.Arguments {
	return comment
}

func updateDatabaseOnSuccess(db *bbolt.DB, message string, e rpc.Entry) {
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(message))
		return b.Put([]byte(e.TXID), []byte("done"))
	})
	if err != nil {
		exports.Logs.Error(err, Echo("err updating db"))
		// Handle the error in updating the database

	}
}
func updateContactsOnSuccess(db *bbolt.DB, message string, e rpc.Entry, response monero.IntegratedAddressResponse) {
	address := e.Payload_RPC.Value(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress).(rpc.Address).String()
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(message))
		return b.Put([]byte(response.PaymentID), []byte(address))
	})

	if err != nil {
		exports.Logs.Error(err, Echo("err updating db"))
		// Handle the error in updating the database

	}
}
