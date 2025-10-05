package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken menghasilkan token random (hex string) dengan panjang n byte.
// Misalnya n = 32 -> hasil token panjang 64 karakter.
func GenerateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
