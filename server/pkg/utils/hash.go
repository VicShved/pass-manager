package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

// HashSha256 return hash of input string
func HashSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	result := h.Sum(nil)
	return base64.RawStdEncoding.EncodeToString(result)
}
