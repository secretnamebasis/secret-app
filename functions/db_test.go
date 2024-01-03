package functions_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/secretnamebasis/secret-app/asserts"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions"
	"github.com/secretnamebasis/secret-app/src"
	"go.etcd.io/bbolt"
)

func TestDB(t *testing.T) {
	if functions.Connection() != true {
		t.Skip("Skipping wallet-related tests. Wallet connection failed.")
	}

	given := fmt.Sprintf("test_%s_%s.bbolt.db", exports.APP_NAME, src.Sha1Sum(exports.DEVELOPER_ADDRESS))

	t.Run(
		"TestCreateDB",
		func(t *testing.T) {
			asserts.DBCreation(t, func(db *bbolt.DB) error {
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
			asserts.DBCreation(t, func(db *bbolt.DB) error {
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

			asserts.DBCreation(t, func(db *bbolt.DB) error {
				err := functions.CreateBucket(db, []byte("SALE"))
				if err != nil {
					return fmt.Errorf("Error creating 'SALE' bucket: %s", err)
				}

				err = asserts.BucketExists(t, db, []byte("SALE"))
				if err != nil {
					return err
				}

				return nil
			})
		},
	)
}
