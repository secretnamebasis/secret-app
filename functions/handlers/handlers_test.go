package handlers_test

import (
	"testing"

	"github.com/secretnamebasis/secret-app/asserts"
	"github.com/secretnamebasis/secret-app/functions/handlers"

	"go.etcd.io/bbolt"
)

func TestHandleIncomingTransfers(t *testing.T) {
	t.Run(
		"Test Incoming Transfers for Pong",
		func(t *testing.T) {

			asserts.DBCreationWithBucket(t, func(db *bbolt.DB) error {

				got := handlers.HandleIncomingTransfers(db)
				if got != nil {
					t.Errorf("got %s", got)
				}
				return nil
			})
		})
}