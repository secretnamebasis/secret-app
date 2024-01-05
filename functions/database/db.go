package database

import (
	"fmt"

	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet"
	"go.etcd.io/bbolt"
)

func CreateBucket(db *bbolt.DB, bucketName []byte) error {
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)

		exports.Logs.Info(
			wallet.
				Echo(
					"Sales Initiated",
				),
		)
		return err
	})
}

func CreateDB(db_name string) (*bbolt.DB, error) {

	db, err := bbolt.Open(db_name, 0600, nil)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	exports.Logs.Info(
		wallet.
			Echo(
				"Database Created",
			),
	)

	return db, nil
}
