package app_test

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/secretnamebasis/secret-app/app"
	"go.etcd.io/bbolt"
)

func TestRunApp(t *testing.T) {

	if app.RunApp() != nil {
		t.Errorf("App is not running when trying to run app")
	}
}

func TestSayVar(t *testing.T) {
	name := "Alixander"
	if app.Echo(name) != name {
		t.Errorf("App is not returning strings")
	}
}

func TestSayHelloVar(t *testing.T) {
	given := "secret"
	got := app.SayHello(given)
	want := "Hello, secret"
	assertCorrectMessage(t, got, want)
}

func TestPing(t *testing.T) {
	got := app.Ping()
	if got != true {
		t.Errorf("App is not returning pinging")
	}
}

func TestWalletEcho(t *testing.T) {
	given := "secret"
	got := app.WalletEcho(given)
	want := "WALLET " + app.Echo(given) + "\n"
	assertCorrectMessage(t, got, want)

}

func TestWalletAddress(t *testing.T) {
	got := app.WalletAddress()
	want := app.DEVELOPER_ADDRESS
	assertCorrectMessage(t, got, want)
}

func TestWalletHeight(t *testing.T) {
	got := app.WalletHeight()
	if got == 0 {
		t.Errorf("got %q", got)
	}
}

func TestLogger(t *testing.T) {
	got := app.Logger()
	if got != nil {
		t.Errorf("got %q", got)
	}
}

func TestSha1SumCrypto(t *testing.T) {
	given := app.DEVELOPER_ADDRESS
	got := app.Sha1Sum(given)
	want := fmt.Sprintf("%x", sha1.Sum([]byte(app.DEVELOPER_ADDRESS)))
	if got != want {
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
	authTransport := &app.TransportWithBasicAuth{
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

func TestCreateDB(t *testing.T) {
	given := fmt.Sprintf("test_%s_%s.bbolt.db", app.APP_NAME, app.Sha1Sum(app.DEVELOPER_ADDRESS))

	t.Run(
		"create",
		func(t *testing.T) {
			assertDBCreation(t, given, func(db *bbolt.DB) error {
				_, err := os.Stat(given)
				if err != nil {
					return fmt.Errorf("Error checking file existence: %s", err)
				}
				return nil
			})
		},
	)
}

func TestCreateSalesBucket(t *testing.T) {
	given := fmt.Sprintf("test_%s_%s.bbolt.db", app.APP_NAME, app.Sha1Sum(app.DEVELOPER_ADDRESS))

	assertDBCreation(t, given, func(db *bbolt.DB) error {
		err := app.CreateBucket(db, []byte("SALE"))
		if err != nil {
			return fmt.Errorf("Error creating 'SALE' bucket: %s", err)
		}

		err = assertBucketExists(t, db, []byte("SALE"))
		if err != nil {
			return err
		}

		return nil
	})
}

func TestUpdateDB(t *testing.T) {
	given := fmt.Sprintf("test_%s_%s.bbolt.db", app.APP_NAME, app.Sha1Sum(app.DEVELOPER_ADDRESS))

	t.Run(
		"update",
		func(t *testing.T) {
			assertDBCreation(t, given, func(db *bbolt.DB) error {
				return db.Update(func(tx *bbolt.Tx) error {
					_, err := tx.CreateBucketIfNotExists([]byte("SALE"))
					return err
				})
			})
		},
	)
}

func TestCreateServiceAddress(t *testing.T) {
	got := app.CreateServiceAddress()
	if got == "" {
		t.Fatalf("got %s", got)
	}
}

func assertDBCreation(t *testing.T, given string, fn func(db *bbolt.DB) error) {

	given = fmt.Sprintf("test_%s_%s.bbolt.db", app.APP_NAME, app.Sha1Sum(app.DEVELOPER_ADDRESS))

	defer func() {
		err := os.Remove(given)
		if err != nil {
			t.Errorf("Error cleaning up: %s", err)
		}
	}()

	db, err := app.CreateDB(given)
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
