package code

import (
	"fmt"

	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/site"
	"github.com/secretnamebasis/secret-app/site/config"
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

	go func() {
		app := site.MakeWebsite()
		config := config.Server{Port: 3000}
		if err := site.StartServer(app, config.Port); err != nil {
			exports.Logs.Error(err, "Error starting server")
		}
	}()

	if err := performWalletOperations(deroDB); err != nil {
		return fmt.Errorf("Failed to perform wallet operations: %v", err)
	}

	return nil
}
