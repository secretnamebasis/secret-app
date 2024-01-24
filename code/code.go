package code

import (
	"fmt"

	"github.com/secretnamebasis/secret-app/functions/wallet/dero"

)

const (
	createBucket   = "create"
	saleBucket     = "sale"
	contactsBucket = "contacts"
)

func Run() error {
	setupLogger()

	if err := checkWalletConnection(); err != nil {
		return fmt.Errorf("Failed to establish wallet connection: %v", err)
	}

	if err := checkMoneroConnection(); err != nil {
		return fmt.Errorf("Failed to establish Monero wallet connection: %v", err)
	}

	deroDB, deroDBName, err := createDeroDB()
	if err != nil {
		return fmt.Errorf("Failed to create DERO database: %v", err)
	}

	logWalletInfo(deroDBName, dero.Address())

	if err := performWalletOperations(deroDB); err != nil {
		return fmt.Errorf("Failed to perform wallet operations: %v", err)
	}

	return nil
}
