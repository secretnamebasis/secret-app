package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	asserts_tests "github.com/secretnamebasis/secret-app/asserts"
	"github.com/secretnamebasis/secret-app/exports"

	"github.com/secretnamebasis/secret-app/functions/crypto"
	"github.com/secretnamebasis/secret-app/functions/database"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"go.etcd.io/bbolt"
)

func TestDB(t *testing.T) {
	if dero.Connection() != true {
		t.Skip("Skipping wallet-related tests. Wallet connection failed.")
	}

	given := fmt.Sprintf("test_%s_%s.bbolt.db", exports.APP_NAME, crypto.Sha1Sum(exports.DEVELOPER_ADDRESS))
	t.Run("TestOpenDB",
		func(t *testing.T) {
			asserts_tests.DBCreation(t, func(db *bbolt.DB) error {
				db_name := crypto.Sha1Sum(dero.Address())
				db, err := bbolt.Open(db_name, 0600, nil)
				if err != nil {
					fmt.Printf(err.Error())
					return err
				}
				log.Printf("Database '%s' created successfully", db_name)
				exports.Logs.Info(
					dero.Echo(
						"Database Created",
					),
				)

				return nil

			})
		},
	)
	t.Run(
		"TestCreateDB",
		func(t *testing.T) {
			asserts_tests.DBCreation(t, func(db *bbolt.DB) error {
				_, err := os.Stat(given)
				if err != nil {
					return fmt.Errorf("Error checking file existence: %s", err)
				}
				return nil
			})
		},
	)
	t.Run(
		"TestUpdateDB",
		func(t *testing.T) {
			asserts_tests.DBCreation(t, func(db *bbolt.DB) error {
				return db.Update(func(tx *bbolt.Tx) error {
					_, err := tx.CreateBucketIfNotExists([]byte("SALE"))
					return err
				})
			})
		},
	)

	t.Run(
		"TestCreateSalesBucket",
		func(t *testing.T) {

			asserts_tests.DBCreation(t, func(db *bbolt.DB) error {
				err := database.CreateBucket(db, []byte("SALE"))
				if err != nil {
					return fmt.Errorf("Error creating 'SALE' bucket: %s", err)
				}

				err = asserts_tests.BucketExists(t, db, []byte("SALE"))
				if err != nil {
					return err
				}

				return nil
			})
		},
	)
}
