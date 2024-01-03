package code_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/secretnamebasis/secret-app/code"
	"github.com/secretnamebasis/secret-app/functions"

	"github.com/secretnamebasis/secret-app/exports"
)

func TestRunApp(t *testing.T) {
	Testing := true
	if !Testing {

		given := fmt.Sprintf("%s_%s.bbolt.db", exports.APP_NAME, functions.Sha1Sum(exports.DEVELOPER_ADDRESS))
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
