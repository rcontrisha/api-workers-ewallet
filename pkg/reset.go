package pkg

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateRandomToken nge-generate cryptographically secure random string.
// Parameter byteCount idealnya 32 (akan menghasilkan 64 karakter hex).
func GenerateRandomToken(byteCount int) (string, error) {
	bytes := make([]byte, byteCount)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken nge-hash raw token pakai SHA-256.
// Hasil hash ini yang akan dipakai sebagai Key di Redis (misal: "reset_password:<HASH>").
func HashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}
