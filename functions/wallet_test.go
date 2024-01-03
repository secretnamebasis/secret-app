package functions_test

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/secretnamebasis/secret-app/code"
	"github.com/secretnamebasis/secret-app/code/exports"
	"github.com/secretnamebasis/secret-app/functions"
)

func TestWallet(t *testing.T) {
	if t.Run(
		"TestWalletConnection",
		func(t *testing.T) {
			got := functions.Connection()
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
			got := functions.Echo(given)
			want := "WALLET " + functions.Echo(given) + "\n"
			assertCorrectMessage(t, got, want)
		},
	)
	t.Run(
		"TestWalletAddress",
		func(t *testing.T) {
			got := functions.Address()
			want := exports.DEVELOPER_ADDRESS
			assertCorrectMessage(t, got, want)

		},
	)
	t.Run(
		"TestWalletAddressSha1Sum",
		func(t *testing.T) {
			given := exports.DEVELOPER_ADDRESS
			got := code.Sha1Sum(given)
			want := fmt.Sprintf("%x", sha1.Sum([]byte(exports.DEVELOPER_ADDRESS)))
			if got != want {
				t.Errorf("got %q", got)
			}
		})
	t.Run(
		"TestWalletHeight",
		func(t *testing.T) {
			got := functions.Height()
			if got == 0 {
				t.Errorf("got %q", got)
			}
		},
	)
	t.Run(
		"TestWalletCreateServiceAddress",
		func(t *testing.T) {
			given := exports.DEVELOPER_ADDRESS
			got := functions.CreateServiceAddress(given)
			if got == "" {
				t.Errorf("got %s", got)
			}
		},
	)
	t.Run(
		"TestWalletCreateServiceAddressWithoutHardcodedValue",
		func(t *testing.T) {
			got := functions.CreateServiceAddressWithoutHardcodedValue(functions.CreateServiceAddress(exports.DEVELOPER_ADDRESS))
			if got == "" {
				t.Errorf("got %s", got)
			}
		},
	)
	t.Run(
		"TestWalletGetTransfers",
		func(t *testing.T) {
			_, got := functions.GetTransfers()
			if got != nil {
				t.Errorf(got.Error())
			}
		},
	)
}
