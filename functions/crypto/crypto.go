package crypto

import (
	"crypto/sha1"
	"fmt"

	"github.com/secretnamebasis/secret-app/functions/wallet"
)

func Sha1Sum(DEVELOPER_ADDRESS string) string {
	shasum := fmt.Sprintf("%x", sha1.Sum([]byte(wallet.Address())))
	return shasum
}
