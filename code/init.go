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
	_, err := dero.Connection()
	if err != nil {
		return errors.New("wallet Connection Failure")
	}
	return nil
}

func checkMoneroConnection() error {
	if monero.Height() <= 0 {
		return errors.New("monero Wallet Connection Failure")
	}
	return nil
}

func makeDBName(s string) string {
	return fmt.Sprintf("%s_%s.bbolt.db", exports.APP_NAME, crypto.Sha1Sum(s))
}

func createDeroDB() (*bbolt.DB, string, error) {
	addr, _ := dero.Address()
	deroDBName := makeDBName(addr)
	db, err := createDB(deroDBName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create database: %v", err)
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
	if exports.Testing {
		return nil
	}
	return wallet.IncomingTransfers(deroDB)
}

func logWalletInfo(deroDBName string) {
	addr, _ := dero.Address()
	exports.Logs.Info(dero.Echo("DB"), "DB", deroDBName)
	exports.Logs.Info(dero.Echo("Address"), "Monero", monero.Address(0))
	exports.Logs.Info(dero.Echo("Address"), "DERO", addr)
}
