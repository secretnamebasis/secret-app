package src_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/secretnamebasis/secret-app/asserts"

	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions"

	"go.etcd.io/bbolt"
)

func TestRunApp(t *testing.T) {
	Testing := true
	if !Testing {

		given := fmt.Sprintf("%s_%s.bbolt.db", exports.APP_NAME, code.Sha1Sum(exports.DEVELOPER_ADDRESS))
		defer func() {
			err := os.Remove(given)
			if err != nil {
				t.Errorf("Error cleaning up: %s", err)
			}
		}()

		if code.RunApp() != nil {
			t.Errorf("App is not running when trying to run app")
		}
	}

}
func TestSayVar(t *testing.T) {
	username := "secret"
	if functions.SayEcho(username) != username {
		t.Errorf("App is not returning strings")
	}
}

func TestSayHelloVar(t *testing.T) {
	given := "secret"
	got := functions.SayHello(given)
	want := "Hello, secret"
	asserts.CorrectMessage(t, got, want)
}

func TestPing(t *testing.T) {
	got := functions.Ping()
	if got != true {
		t.Errorf("App is not returning pinging")
	}
}

func TestHandleIncomingTransfers(t *testing.T) {

	asserts.DBCreationWithBucket(t, func(db *bbolt.DB) error {

		got := code.HandleIncomingTransfers(db)
		if got != nil {
			t.Errorf("got %s", got)
		}
		return nil
	})
}

func TestLogger(t *testing.T) {
	got := code.Logger()
	if got != nil {
		t.Errorf("got %q", got)
	}
}

func TestRoundTrip(t *testing.T) {
	// Given
	givenUsername := "testuser"
	givenPassword := "testpassword"

	// Mock HTTP server for testing
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Basic authentication not set")
			return
		}

		// Got
		gotUsername := username
		gotPassword := password

		// Want
		wantUsername := givenUsername
		wantPassword := givenPassword

		// Check if Got matches Want
		if gotUsername != wantUsername || gotPassword != wantPassword {
			t.Errorf("Unexpected basic auth headers. Got username: %s, Want username: %s, Got password: %s,  Want password: %s",
				gotUsername, wantUsername, gotPassword, wantPassword)
		}
	}))
	defer ts.Close()

	// Create a TransportWithBasicAuth instance
	authTransport := &exports.TransportWithBasicAuth{
		Username: givenUsername,
		Password: givenPassword,
		Base:     http.DefaultTransport,
	}

	// Create a request
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// When
	// Execute RoundTrip method
	_, err = authTransport.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}

}

func TestDB(t *testing.T) {
	if functions.Connection() != true {
		t.Skip("Skipping wallet-related tests. Wallet connection failed.")
	}

	given := fmt.Sprintf("test_%s_%s.bbolt.db", exports.APP_NAME, code.Sha1Sum(exports.DEVELOPER_ADDRESS))

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
				err := code.CreateBucket(db, []byte("SALE"))
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
