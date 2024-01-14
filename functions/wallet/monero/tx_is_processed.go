package monero

import "go.etcd.io/bbolt"

func isTransactionProcessed(db *bbolt.DB, bucketName string, e Entry) (bool, error) {
	var already_processed bool

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket != nil {
			if existing := bucket.Get([]byte(e.TxID)); existing != nil {
				already_processed = true
			}
		}
		return nil
	})

	if err != nil {
		return false, err // Return false and the encountered error
	}

	return already_processed, nil
}

func hasPaymentID(db *bbolt.DB, bucketName string, e Entry) (bool, error) {
	var already_processed bool

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket != nil {
			if existing := bucket.Get([]byte(e.PaymentID)); existing != nil {
				already_processed = true
			}
		}
		return nil
	})

	if err != nil {
		return false, err // Return false and the encountered error
	}

	return already_processed, nil
}
