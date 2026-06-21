package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

const (
	KeySize   = 32
	NonceSize = 12
)

func Encrypt(plaintext, key, nonce []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, fmt.Errorf("key must be %d bytes, got %d", KeySize, len(key))
	}

	if len(nonce) != NonceSize {
		return nil, fmt.Errorf("nonce must be %d bytes, got %d", NonceSize, len(nonce))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}
