package asserts_tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/secretnamebasis/secret-app/exports"

	"github.com/secretnamebasis/secret-app/functions/crypto"
	"github.com/secretnamebasis/secret-app/functions/database"
	"go.etcd.io/bbolt"
)

func DBCreationWithBucket(t *testing.T, fn func(db *bbolt.DB) error) {
	DBCreation(
		t, func(db *bbolt.DB) error {
			err := database.CreateBucket(db, []byte("SALE"))
			if err != nil {
				return fmt.Errorf("Error creating 'SALE' bucket: %s", err)
			}

			err = BucketExists(t, db, []byte("SALE"))
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func DBCreation(t *testing.T, fn func(db *bbolt.DB) error) {

	given := fmt.Sprintf("test_%s_%s.bbolt.db", exports.APP_NAME, crypto.Sha1Sum(exports.DEVELOPER_ADDRESS))

	defer func() {
		err := os.Remove(given)
		if err != nil {
			t.Errorf("Error cleaning up: %s", err)
		}
	}()

	db, err := database.CreateDB(given)
	if err != nil {
		t.Fatalf("Error creating database: %s", err)
	}

	err = fn(db)
	if err != nil {
		t.Fatalf("Error during test: %s", err)
	}
}

func BucketExists(t *testing.T, db *bbolt.DB, bucketName []byte) error {
	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Expected '%s' bucket to exist, but it doesn't", bucketName)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%s does not exist: %s", bucketName, err)
	}
	return nil
}

func CorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func ErrMessage(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("err %q", err)
	}
}
