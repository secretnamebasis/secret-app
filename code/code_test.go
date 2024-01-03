package code_test

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/secretnamebasis/secret-app/code"
	"github.com/secretnamebasis/secret-app/functions"
	"go.etcd.io/bbolt"
)

func TestRunApp(t *testing.T) {
	Testing := true
	if !Testing {

		given := fmt.Sprintf("%s_%s.bbolt.db", code.APP_NAME, code.Sha1Sum(functions.Address()))
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
	if code.Echo(username) != username {
		t.Errorf("App is not returning strings")
	}
}

func TestSayHelloVar(t *testing.T) {
	given := "secret"
	got := code.SayHello(given)
	want := "Hello, secret"
	assertCorrectMessage(t, got, want)
}

func TestPing(t *testing.T) {
	got := code.Ping()
	if got != true {
		t.Errorf("App is not returning pinging")
	}
}

func TestWallet(t *testing.T) {
	if t.Run(
		"TestWalletConnection",
		func(t *testing.T) {
			got := functions.Connection()
			if got != true {
				t.Errorf("Your wallet is not connected")
			}
		},
	) != true {
		t.Skip("Skipping wallet-related tests. Wallet connection failed.")
	}
	t.Run(
		"TestWalletEcho",
		func(t *testing.T) {
			given := "secret"
			got := functions.Echo(given)
			want := "WALLET " + code.Echo(given) + "\n"
			assertCorrectMessage(t, got, want)
		},
	)
	t.Run(
		"TestWalletAddress",
		func(t *testing.T) {
			got := functions.Address()
			want := code.DEVELOPER_ADDRESS
			assertCorrectMessage(t, got, want)

		},
	)
	t.Run(
		"TestWalletAddressSha1Sum",
		func(t *testing.T) {
			given := code.DEVELOPER_ADDRESS
			got := code.Sha1Sum(given)
			want := fmt.Sprintf("%x", sha1.Sum([]byte(code.DEVELOPER_ADDRESS)))
			if got != want {
				t.Errorf("got %q", got)
			}
		})
	t.Run(
		"TestWalletHeight",
		func(t *testing.T) {
			got := functions.Height()
			if got == 0 {
				t.Errorf("got %q", got)
			}
		},
	)
	t.Run(
		"TestWalletCreateServiceAddress",
		func(t *testing.T) {
			given := code.DEVELOPER_ADDRESS
			got := code.CreateServiceAddress(given)
			if got == "" {
				t.Errorf("got %s", got)
			}
		},
	)
	t.Run(
		"TestWalletCreateServiceAddressWithoutHardcodedValue",
		func(t *testing.T) {
			got := code.CreateServiceAddressWithoutHardcodedValue(code.CreateServiceAddress(code.DEVELOPER_ADDRESS))
			if got == "" {
				t.Errorf("got %s", got)
			}
		},
	)
	t.Run(
		"TestWalletGetTransfers",
		func(t *testing.T) {
			_, got := functions.GetTransfers()
			if got != nil {
				t.Errorf(got.Error())
			}
		},
	)
}

func TestHandleIncomingTransfers(t *testing.T) {

	assertDBCreationWithBucket(t, func(db *bbolt.DB) error {

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
	authTransport := &code.TransportWithBasicAuth{
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

	given := fmt.Sprintf("test_%s_%s.bbolt.db", code.APP_NAME, code.Sha1Sum(code.DEVELOPER_ADDRESS))

	t.Run(
		"TestCreateDB",
		func(t *testing.T) {
			assertDBCreation(t, func(db *bbolt.DB) error {
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
			assertDBCreation(t, func(db *bbolt.DB) error {
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

			assertDBCreation(t, func(db *bbolt.DB) error {
				err := code.CreateBucket(db, []byte("SALE"))
				if err != nil {
					return fmt.Errorf("Error creating 'SALE' bucket: %s", err)
				}

				err = assertBucketExists(t, db, []byte("SALE"))
				if err != nil {
					return err
				}

				return nil
			})
		},
	)
}

func assertDBCreationWithBucket(t *testing.T, fn func(db *bbolt.DB) error) {
	assertDBCreation(
		t, func(db *bbolt.DB) error {
			err := code.CreateBucket(db, []byte("SALE"))
			if err != nil {
				return fmt.Errorf("Error creating 'SALE' bucket: %s", err)
			}

			err = assertBucketExists(t, db, []byte("SALE"))
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func assertDBCreation(t *testing.T, fn func(db *bbolt.DB) error) {

	given := fmt.Sprintf("test_%s_%s.bbolt.db", code.APP_NAME, code.Sha1Sum(code.DEVELOPER_ADDRESS))

	defer func() {
		err := os.Remove(given)
		if err != nil {
			t.Errorf("Error cleaning up: %s", err)
		}
	}()

	db, err := code.CreateDB(given)
	if err != nil {
		t.Fatalf("Error creating database: %s", err)
	}

	err = fn(db)
	if err != nil {
		t.Fatalf("Error during test: %s", err)
	}
}

func assertBucketExists(t *testing.T, db *bbolt.DB, bucketName []byte) error {
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

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertErrMessage(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("err %q", err)
	}
}
