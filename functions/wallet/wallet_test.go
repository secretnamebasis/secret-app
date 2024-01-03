package wallet_test

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/secretnamebasis/secret-app/asserts"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet"
)

func TestWallet(t *testing.T) {
	if t.Run(
		"TestWalletConnection",
		func(t *testing.T) {
			got := wallet.Connection()
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
			got := wallet.Echo(given)
			want := "WALLET " + given + "\n"
			asserts.CorrectMessage(t, got, want)
		},
	)
	t.Run(
		"TestWalletAddress",
		func(t *testing.T) {
			got := wallet.Address()
			want := exports.DEVELOPER_ADDRESS
			asserts.CorrectMessage(t, got, want)

		},
	)
	t.Run(
		"TestWalletAddressSha1Sum",
		func(t *testing.T) {
			given := exports.DEVELOPER_ADDRESS
			got := wallet.Sha1Sum(given)
			want := fmt.Sprintf("%x", sha1.Sum([]byte(exports.DEVELOPER_ADDRESS)))
			if got != want {
				t.Errorf("got %q", got)
			}
		})
	t.Run(
		"TestWalletHeight",
		func(t *testing.T) {
			got := wallet.Height()
			if got == 0 {
				t.Errorf("got %q", got)
			}
		},
	)
	t.Run(
		"TestWalletCreateServiceAddress",
		func(t *testing.T) {
			given := exports.DEVELOPER_ADDRESS
			got := wallet.CreateServiceAddress(given)
			if got == "" {
				t.Errorf("got %s", got)
			}
		},
	)
	t.Run(
		"TestWalletCreateServiceAddressWithoutHardcodedValue",
		func(t *testing.T) {
			got := wallet.CreateServiceAddressWithoutHardcodedValue(wallet.CreateServiceAddress(exports.DEVELOPER_ADDRESS))
			if got == "" {
				t.Errorf("got %s", got)
			}
		},
	)
	t.Run(
		"TestWalletGetTransfers",
		func(t *testing.T) {
			_, got := wallet.GetTransfers()
			if got != nil {
				t.Errorf(got.Error())
			}
		},
	)
}
