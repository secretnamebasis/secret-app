package crypto

import (
	"crypto/sha1"
	"fmt"
)

func Sha1Sum(s string) string {
	shasum := fmt.Sprintf("%x", sha1.Sum([]byte(s)))
	return shasum
}
