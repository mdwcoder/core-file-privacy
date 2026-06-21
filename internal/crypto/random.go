package crypto

import (
	"crypto/rand"
	"fmt"
	"io"
)

func GenerateSalt(length int) ([]byte, error) {
	return GenerateRandomBytes(length)
}

func GenerateNonce(length int) ([]byte, error) {
	return GenerateRandomBytes(length)
}

func GenerateRandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}
