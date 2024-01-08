package wallet_test

import (
	"testing"
	"time"

	"github.com/deroproject/derohe/rpc"
	asserts_tests "github.com/secretnamebasis/secret-app/asserts"
	"github.com/secretnamebasis/secret-app/functions/wallet"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"go.etcd.io/bbolt"
)

var entry = []rpc.Entry{
	{
		Height:         3114209,
		TopoHeight:     3114209,
		BlockHash:      "49925d620bae1f4f007f5cff4d57ccf468edf2361649a48bafda22127e4ef168",
		MinerReward:    0,
		TransactionPos: 0,
		Pos:            0,
		Coinbase:       false,
		Incoming:       true,
		TXID:           "ad1d19d5f74a147037112c80f58aec99c69a8f19cb0e950876ca8f7fbfa49c40",
		Destination:    "",
		Amount:         28000,
		Fees:           181,
		Proof:          "deroproof1qyw8ed6u0r500de7zcjhw6y4u7equcgngzms7dlhj9fk5w8mjrsn2qdzvfyyskpqrry09zjphyqgzhyqtx8vjm5n4k9e36le2yusfuvamqn0a7ttlutxy4j4r9kkq85kfre",
		Status:         0,
		Time:           time.Date(2024, 1, 4, 12, 18, 21, 242000000, time.FixedZone("", -7*60*60)),
		EWData:         "102d35d9ed14f467cd7aa08f0aecc8bcf48a4156311dd3d6680c7073d6567fbb01195b6e060776c8d1e2615fda8801377a094ff111544a36afa1b5db2b6188ce3301",
		Data:           []byte("DqJiQ1NgYkRVAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="),
		PayloadType:    0,
		Payload:        []byte("omJDU2BiRFUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		Payload_RPC: rpc.Arguments{
			{
				Name:     "C",
				DataType: "S",
				Value:    "",
			},
			{
				Name:     "D",
				DataType: "U",
				Value:    0,
			},
		},
		Sender:          "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g",
		DestinationPort: 0,
		SourcePort:      0,
	},
	{
		Height:         0,
		TopoHeight:     0,
		BlockHash:      "",
		MinerReward:    0,
		TransactionPos: 0,
		Pos:            0,
		Coinbase:       false,
		Incoming:       true,
		TXID:           "",
		Destination:    "",
		Amount:         0,
		Fees:           0,
		Proof:          "",
		Status:         0,
		Time:           time.Date(2024, 1, 4, 12, 18, 21, 242000000, time.FixedZone("", -7*60*60)),
		EWData:         "",
		Data:           []byte("=="),
		PayloadType:    0,
		Payload:        []byte(""),
		Payload_RPC: rpc.Arguments{
			{
				Name:     "C",
				DataType: "S",
				Value:    "",
			},
			{
				Name:     "D",
				DataType: "U",
				Value:    0,
			},
		},
		Sender:          "",
		DestinationPort: 0,
		SourcePort:      0,
	},
}

func TestHandleIncomingTransfers(t *testing.T) {
	t.Run(
		"Test Handle Incoming Transfers",
		func(t *testing.T) {

			asserts_tests.DBCreationWithBucket(
				t,
				func(db *bbolt.DB) error {

					got := wallet.IncomingTransfers(db)
					if got == nil {
						t.Errorf("got %s", got)
					}
					return nil
				},
			)
		},
	)
	t.Run(
		"Test Handle Incoming Transfer Entry",
		func(t *testing.T) {

			asserts_tests.DBCreationWithBucket(
				t,
				func(db *bbolt.DB) error {

					given := entry[1]
					given.Amount = 0

					if got := dero.IncomingTransferEntry(given, db); got != nil {
						t.Errorf("got %s", got)
					}

					return nil
				},
			)
		},
	)
	// t.Run("Test Incoming Transfer Entry Switch",
	// 	func(t *testing.T) {
	// 		asserts_tests.DBCreationWithBucket(
	// 			t, func(db *bbolt.DB) error {
	// 				given := entry[1]
	// 				if got := handlers.IncomingTransferEntrySwitch(given, db); got != nil {
	// 					t.Errorf("got %s", got)
	// 				}
	// 				return nil
	// 			},
	// 		)
	// 	},
	// )
}
