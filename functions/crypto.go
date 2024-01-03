package functions

import (
	"crypto/sha1"
	"fmt"
)

func Sha1Sum(DEVELOPER_ADDRESS string) string {
	shasum := fmt.Sprintf("%x", sha1.Sum([]byte(Address())))
	return shasum
}
