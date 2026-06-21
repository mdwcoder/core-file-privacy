package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	keyPrefix    = "cfp_"
	keyGroups    = 8
	keyGroupSize = 4
	keyChars     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateKey() (string, error) {
	var groups []string

	for i := 0; i < keyGroups; i++ {
		var group strings.Builder
		for j := 0; j < keyGroupSize; j++ {
			idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(keyChars))))
			if err != nil {
				return "", fmt.Errorf("failed to generate random index: %w", err)
			}
			group.WriteByte(keyChars[idx.Int64()])
		}
		groups = append(groups, group.String())
	}

	return keyPrefix + strings.Join(groups, "-"), nil
}

func ValidateKey(key string) bool {
	if !strings.HasPrefix(key, keyPrefix) {
		return false
	}

	parts := strings.Split(key[len(keyPrefix):], "-")
	if len(parts) != keyGroups {
		return false
	}

	for _, part := range parts {
		if len(part) != keyGroupSize {
			return false
		}
		for _, c := range part {
			if !strings.ContainsRune(keyChars, c) {
				return false
			}
		}
	}

	return true
}
