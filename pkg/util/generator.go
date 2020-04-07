package util

import (
	"crypto/rand"
	"fmt"
)

// GenerateToken generates a 16-byte access token
func GenerateToken() string {
	return tokenGenerator(16)
}

// GenerateShortToken generates a 8-byte generic token
func GenerateShortToken() string {
	return tokenGenerator(8)
}

func tokenGenerator(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
