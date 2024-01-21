package tests

import (
	"testing"

	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/monero"
)

var entry = []monero.Entry{
	{
		Address:                         "42oK8BJRrY5DYXbWxMS5j3Zamjkmsk6vDRS8tRR5TUFJTggTKovWzkien1Vp8bXvKAP1hDFJwZjxUgRqjfmY9sNPFvSea4w",
		Amount:                          1000000000,
		Amounts:                         []uint64{1000000000},
		Confirmations:                   103963,
		DoubleSpendSeen:                 false,
		Fee:                             44440000,
		Height:                          2953189,
		Locked:                          false,
		Note:                            "",
		PaymentID:                       "1e70dcce10f38392",
		SubaddrIndex:                    monero.SubaddrIndex{Major: 0, Minor: 0},
		SubaddrIndices:                  []monero.SubaddrIndex{{Major: 0, Minor: 0}},
		SuggestedConfirmationsThreshold: 1,
		Timestamp:                       1692187838,
		TxID:                            "7c1238251d6cd215726f9ad716b6a0dc2b24b17d3b621ad875d252a62542cdcf",
		Type:                            "in",
		UnlockTime:                      0,
	},
}

func TestMoneroWallet(t *testing.T) {

	t.Run("Test Monero Height",
		func(t *testing.T) {
			got := monero.Height()
			if got == 0 {
				t.Errorf("Error obtaining height, got %v", got)
			}
		},
	)
	t.Run("Test Monero Transfers",
		func(t *testing.T) {
			_, got := monero.GetIncomingTransfers()
			if got != nil {
				t.Errorf(got.Error())
			}
		},
	)
	t.Run("Test Monero Address", func(t *testing.T) {
		given := uint64(0)
		got := monero.Address(given)
		if got == "" {
			t.Errorf("err.Error()")
		}
	},
	)
	t.Run("Test Valid Monero Address", func(t *testing.T) {
		given := "42oK8BJRrY5DYXbWxMS5j3Zamjkmsk6vDRS8tRR5TUFJTggTKovWzkien1Vp8bXvKAP1hDFJwZjxUgRqjfmY9sNPFvSea4w"
		got := monero.ValidateAddress(given)
		if got != true {
			t.Errorf("err.Error()")
		}
	},
	)

	t.Run("Test Monero Transfers by Height", func(t *testing.T) { // this shiz is crazy
		given := 2953189
		got, _ := monero.GetIncomingTransfersByHeight(given)

		if got.In == nil {
			t.Errorf("%v", got.In)
		}
	})

	t.Run("Test Wallet Integrated Address",
		func(t *testing.T) {
			result, err := monero.MakeIntegratedAddress()
			if err != nil {
				t.Errorf("Error obtained %s", err)
				return
			}
			// Create IntegratedAddressResponse instance to get integrated address and payment ID
			var response monero.IntegratedAddressResponse
			response.IntegratedAddress = result["integrated_address"]
			response.PaymentID = result["payment_id"]

			if response.IntegratedAddress == "" {
				t.Errorf("Integrated address is empty")
			}

			if response.PaymentID == "" {
				t.Errorf("Payment ID is empty")
			}
		},
	)
}
func TestSendXMRTransfer(t *testing.T) {
	where := exports.TEST_XMR_ADDRESS

	given := monero.Transfer_Params{
		Destinations: []monero.Transfer{
			{
				Amount:  100000000, // Amount in atomic units, adjust as needed
				Address: where,
				// Add other fields as needed
			},
			// Add more Transfer objects as needed
		},
	}

	got := monero.SendTransfer(given)
	if len(got.Result.TxHash) != 64 {
		t.Errorf("Failed to send Monero transfer: %s", got.Result.TxHash)
	}
}
