package crypto

import (
	"crypto/sha1"
	"fmt"

	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
)

func Sha1Sum(DEVELOPER_ADDRESS string) string {
	shasum := fmt.Sprintf("%x", sha1.Sum([]byte(dero.Address())))
	return shasum
}
