package dero_test

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/deroproject/derohe/rpc"
	asserts_tests "github.com/secretnamebasis/secret-app/asserts"
	"github.com/secretnamebasis/secret-app/exports"
	crypto "github.com/secretnamebasis/secret-app/functions/crypto"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/functions/wallet/generators"
)

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
