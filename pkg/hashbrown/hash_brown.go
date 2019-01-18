package hashbrown

import (
	"crypto/sha512"
	"encoding/base64"
)

// Create takes a password (string) and hashes it (through SHA512)
// then returns a base64 encoded version of the digest
func Create(password string) string {
	h := sha512.New()
	h.Write([]byte(password))
	b := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return b
}
