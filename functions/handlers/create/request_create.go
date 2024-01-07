package create

import (
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/logger"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/functions/wallet/monero"
	"go.etcd.io/bbolt"
)

func handleCreateRequest(e rpc.Entry, message string, db *bbolt.DB) {
	logger.RequestInfo(e, message+" request")

	reply := createTransfer(e)
	result := dero.SendTransfer(reply)

	if result != "" {
		updateDatabaseOnSuccess(db, message, e.TXID)
		exports.Logs.Info(dero.Echo("ping replied successfully with pong "), "result", result)
		// Log the successful completion
	}
}

func createTransfer(e rpc.Entry) rpc.Transfer_Params {
	where := getAddressFromEntry(e)
	integratedAddress := getIntegratedAddress()
	comment := getCommentFromEntry(integratedAddress)

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

func getIntegratedAddress() string {
	request, _ := monero.MakeIntegratedAddress()
	var response monero.IntegratedAddressResponse
	response.IntegratedAddress = request["integrated_address"]
	return response.IntegratedAddress
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

func updateDatabaseOnSuccess(db *bbolt.DB, message string, txID string) {
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(message))
		return b.Put([]byte(txID), []byte("done"))
	})
	if err != nil {
		exports.Logs.Error(err, dero.Echo("err updating db"))
		// Handle the error in updating the database

	}
}
