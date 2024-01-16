package code

import (
	"errors"
	"fmt"

	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/crypto"
	"github.com/secretnamebasis/secret-app/functions/database"
	"github.com/secretnamebasis/secret-app/functions/wallet"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/functions/wallet/monero"
	"go.etcd.io/bbolt"
)

func setupLogger() {
	dero.Logger()
	exports.Logs.Info(dero.Echo("Logger has started"))
}

func checkWalletConnection() error {
	if !dero.Connection() {
		return errors.New("Wallet Connection Failure")
	}
	return nil
}

func checkMoneroConnection() error {
	if monero.Height() <= 0 {
		return errors.New("Monero Wallet Connection Failure")
	}
	return nil
}

func makeDBName(s string) string {
	return fmt.Sprintf("%s_%s.bbolt.db", exports.APP_NAME, crypto.Sha1Sum(s))
}

func createDeroDB() (*bbolt.DB, string, error) {
	deroDBName := makeDBName(dero.Address())
	db, err := createDB(deroDBName)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to create database: %v", err)
	}
	return db, deroDBName, nil
}

func createDB(dbName string) (*bbolt.DB, error) {
	db, err := database.CreateDB(dbName)
	if err != nil {
		exports.Logs.Error(err, err.Error())
		return nil, err
	}
	createBuckets(db)
	return db, nil
}

func createBuckets(db *bbolt.DB) {
	buckets := []string{createBucket, saleBucket, contactsBucket}
	for _, bucket := range buckets {
		database.CreateBucket(db, []byte(bucket))
	}
}

func performWalletOperations(deroDB *bbolt.DB) error {
	if exports.Testing == true {
		return nil
	}
	return wallet.IncomingTransfers(deroDB)
}

func logWalletInfo(deroDBName string, deroAddress string) {
	exports.Logs.Info(dero.Echo("DERO ID Created: " + deroDBName))
	exports.Logs.Info(dero.Echo("DERO Address: " + deroAddress))
	exports.Logs.Info(dero.Echo("Monero Address: " + monero.Address(0)))
	exports.Logs.Info(dero.Echo("DERO Integrated Address with Expected Arguments: " + dero.CreateServiceAddress(deroAddress)))
	exports.Logs.Info(dero.Echo("DERO Integrated Address without Hardcoded Value: " + dero.CreateServiceAddressWithoutHardcodedValue(deroAddress)))
}
