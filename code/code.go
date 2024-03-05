package code

import (
	"fmt"

	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/site"
	"github.com/secretnamebasis/secret-app/site/config"
)

const (
	createBucket   = "create"
	saleBucket     = "sale"
	contactsBucket = "contacts"
)

func webapp() error {

	app := site.MakeWebsite()
	app.ListenTLS(
		":443",
		"/etc/letsencrypt/live/secretnamebasis.site/cert.pem",
		"/etc/letsencrypt/live/secretnamebasis.site/privkey.pem",
	)
	config := config.Server{Port: 443}
	if err := site.StartServer(app, config.Port); err != nil {
		exports.Logs.Error(err, "Error starting server")
	}
	return nil
}

func Run() error {

	setupLogger()

	if err := checkWalletConnection(); err != nil {
		return fmt.Errorf("failed to establish DERO wallet connection: %v", err)
	}

	if err := checkMoneroConnection(); err != nil {
		return fmt.Errorf("failed to establish Monero wallet connection: %v", err)
	}

	deroDB, deroDBName, err := createDeroDB()
	if err != nil {
		return fmt.Errorf("failed to create DERO database: %v", err)
	}

	logWalletInfo(deroDBName)

	err = webapp()
	if err != nil {
		return fmt.Errorf("failed to startwebsite: %v", err)
	}

	go func() error {
		if err := performWalletOperations(deroDB); err != nil {
			return fmt.Errorf("failed to perform wallet operations: %v", err)
		}
		return nil
	}()

	return nil
}
