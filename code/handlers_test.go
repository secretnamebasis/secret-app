package code_test

import (
	"testing"

	"github.com/secretnamebasis/secret-app/asserts"
	"github.com/secretnamebasis/secret-app/code"

	"go.etcd.io/bbolt"
)

func TestHandleIncomingTransfers(t *testing.T) {

	asserts.DBCreationWithBucket(t, func(db *bbolt.DB) error {

		got := code.HandleIncomingTransfers(db)
		if got != nil {
			t.Errorf("got %s", got)
		}
		return nil
	})
}
