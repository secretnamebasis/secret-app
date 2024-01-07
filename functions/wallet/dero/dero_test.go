package dero_test

import (
	"crypto/sha1"
	"fmt"
	"testing"
	"time"

	"github.com/deroproject/derohe/rpc"
	asserts_tests "github.com/secretnamebasis/secret-app/asserts"
	"github.com/secretnamebasis/secret-app/exports"
	crypto "github.com/secretnamebasis/secret-app/functions/crypto"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/functions/wallet/generators"
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

					got := dero.IncomingTransfers(db)
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
func TestWallet(t *testing.T) {
	if t.Run(
		"TestWalletConnection",
		func(t *testing.T) {
			got := dero.Connection()
			if got != true {
				t.Errorf("Your wallet is not connected")
			}
		},
	) != true {
		t.Skip("Skipping wallet-related tests. Wallet connection failed.")
	}
	t.Run(
		"TestWalletEcho",
		func(t *testing.T) {
			given := "secret"
			got := dero.Echo(given)
			want := "WALLET " + given + "\n"
			asserts_tests.CorrectMessage(t, got, want)
		},
	)
	t.Run(
		"TestWalletAddress",
		func(t *testing.T) {
			got := dero.Address()
			want := exports.DEVELOPER_ADDRESS
			asserts_tests.CorrectMessage(t, got, want)

		},
	)
	t.Run(
		"TestWalletAddressSha1Sum",
		func(t *testing.T) {
			given := exports.DEVELOPER_ADDRESS
			got := crypto.Sha1Sum(given)
			want := fmt.Sprintf("%x", sha1.Sum([]byte(exports.DEVELOPER_ADDRESS)))
			if got != want {
				t.Errorf("got %q", got)
			}
		})
	t.Run(
		"TestWalletHeight",
		func(t *testing.T) {
			got := dero.Height()
			if got == 0 {
				t.Errorf("got %q", got)
			}
		},
	)
	t.Run(
		"TestWalletCreateServiceAddress",
		func(t *testing.T) {
			given := exports.DEVELOPER_ADDRESS
			got := dero.CreateServiceAddress(given)
			if got == "" {
				t.Errorf("got %s", got)
			}
		},
	)
	t.Run(
		"TestWalletCreateServiceAddressWithoutHardcodedValue",
		func(t *testing.T) {
			got := dero.CreateServiceAddressWithoutHardcodedValue(dero.CreateServiceAddress(exports.DEVELOPER_ADDRESS))
			if got == "" {
				t.Errorf("got %s", got)
			}
		},
	)
	t.Run(
		"TestWalletGetTransfers",
		func(t *testing.T) {
			_, got := dero.GetIncomingTransfers()
			if got != nil {
				t.Errorf(got.Error())
			}
		},
	)
	t.Run(
		"TestWalletGetTransfersByHeight",
		func(t *testing.T) {
			given := dero.Height()
			_, got := dero.GetIncomingTransfersByHeight(given)
			if got != nil {
				t.Errorf(got.Error())
			}
		},
	)
	var SendTests = false
	if SendTests {

		t.Run(
			"TestSendTransfer",
			func(t *testing.T) {
				where := exports.SECRET_WALLET_ADDRESS

				given := rpc.Transfer_Params{
					Transfers: []rpc.Transfer{
						{
							Destination: where,
							Amount:      uint64(1),
							Payload_RPC: exports.Response,
						},
					},
				}

				got := dero.SendTransfer(given)
				if len(got) > 64 {
					t.Errorf("got: %s", got)
				}

			},
		)
		t.Run(
			"TestSendTransfers",
			func(t *testing.T) {
				where := []string{
					// This will always result in txid
					// and a successful wallet write
					// but as an error type write
					// to prevent double spends
					// this the error address dero1qy2jy9yjj50w2vgdefhssn738qdu8fdt58nnmm8lt7kx6msaey0n2qqtwluyy
					// not secret wallet address dero1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqg2ctjn4
					/*
						03 Jan 24 21:14 MST
						Height 3111337
						TopoHeight 3111337
						transaction 801c1c76114525d47fcca0f5cf7648ecd4d023b7b0762d16714240b24fc3c13d
						spent 0.00001
						DERO Destination: dero1qy2jy9yjj50w2vgdefhssn738qdu8fdt58nnmm8lt7kx6msaey0n2qqtwluyy
						Proof: deroproof1qysa92k2neq95fd9q20kcg0t4pw577fm07uxrnxyehkddmjh40mwqq9zvfyyskpqev8dh9usn8vl287gcuepyuugqqls2yhyny5mnvd2dn6xgfg4a53ky4j4qy7sv4hc
						RPC CALL arguments [
							Name:C
							Type:string
							Value:'You have purchased a really cool link. You sent 0.12345 at height 1462136'
							Name:D Type:uint64
							Value:'0'
							Name:S
							Type:uint64
							Value:'1311768465173141112']

						03 Jan 24 21:14 MST
						Height 3111336
						TopoHeight 3111336
						transaction 3ea1b66486aa6f429f69a8076808d05db071113b412e6a1b5b5c4cc4ec1d6130
						spent 0.00001
						DERO Destination: dero1qy2jy9yjj50w2vgdefhssn738qdu8fdt58nnmm8lt7kx6msaey0n2qqtwluyy
						Proof: deroproof1qyxd4zkumv9yynlgq5hp7f02djshf2snrt43usdyjat8vwkha77y7qdzvfyyskpqqk7qaze0zfn8l84ktsvykttxxp58flh7qwulwusxgm9zxz00w2qky4j4qyqynmsj
						RPC CALL arguments [
							Name:C
							Type:string
							Value:'You have purchased a really cool link. You sent 0.12345 at height 1460491'
							Name:D Type:uint64
							Value:'0'
							Name:S
							Type:uint64
							Value:'1311768465173141112']


					*/
					exports.SECRET_WALLET_ADDRESS,
					exports.SECRET_WALLET_ADDRESS,
					exports.SECRET_WALLET_ADDRESS,
					exports.SECRET_WALLET_ADDRESS,
					exports.SECRET_WALLET_ADDRESS,
					exports.SECRET_WALLET_ADDRESS,
					exports.SECRET_WALLET_ADDRESS,
					exports.SECRET_WALLET_ADDRESS,
				}
				count := []uint64{
					1,
					1,
					1,
					1,
					1,
					1,
					1,
					1,
				}
				response := []rpc.Arguments{
					exports.Response,
					exports.Response,
					exports.Response,
					exports.Response,
					exports.Response,
					exports.Response,
					exports.Response,
					exports.Response,
				}
				given := rpc.Transfer_Params{
					Transfers: generators.Transfers(where, count, response),
				}

				got := dero.SendTransfer(given)
				if len(got) > 64 {
					t.Errorf("got: %s", got)
				}

			},
		)
	}

}
