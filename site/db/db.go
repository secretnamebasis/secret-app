package db

import "go.etcd.io/bbolt"

var (
	DB     *bbolt.DB
	bucket = []byte("items")
)

func InitDB() error {
	var err error
	DB, err = bbolt.Open("items.db", 0600, nil)
	if err != nil {
		return err
	}

	// Create a bucket if it doesn't exist
	err = DB.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})

	return err
}
