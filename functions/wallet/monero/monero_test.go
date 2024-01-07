package monero_test

import (
	"testing"

	"github.com/secretnamebasis/secret-app/functions/wallet/monero"
)

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
