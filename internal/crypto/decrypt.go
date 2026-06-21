package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func Decrypt(ciphertext, key, nonce []byte) ([]byte, error) {
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

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: wrong password/key or corrupted file")
	}

	return plaintext, nil
}
