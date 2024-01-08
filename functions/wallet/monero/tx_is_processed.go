package monero

import "go.etcd.io/bbolt"

func isTransactionProcessed(db *bbolt.DB, bucketName string, TxID string) (bool, error) {
	var alreadyProcessed bool

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket != nil {
			if existing := bucket.Get([]byte(TxID)); existing != nil {
				alreadyProcessed = true
			}
		}
		return nil
	})

	if err != nil {
		return false, err // Return false and the encountered error
	}

	return alreadyProcessed, nil
}
