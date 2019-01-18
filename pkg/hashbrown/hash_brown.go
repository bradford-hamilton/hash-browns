package hashbrown

import (
	"crypto/sha512"
	"encoding/base64"
)

// Create takes a password and hashes it (SHA512) then base64s the digest
func Create(password string) string {
	h := sha512.New()
	h.Write([]byte(password))
	b := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return b
}
