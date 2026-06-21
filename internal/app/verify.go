package app

import (
	"fmt"
	"os"

	"github.com/core-file-privacy/core-file-privacy/internal/crypto"
	"github.com/core-file-privacy/core-file-privacy/internal/format"
)

type VerifyOptions struct {
	Path     string
	Password string
}

func Verify(opts VerifyOptions) error {
	if opts.Password == "" {
		return fmt.Errorf("password is required")
	}

	container, err := os.ReadFile(opts.Path)
	if err != nil {
		return fmt.Errorf("failed to read container: %w", err)
	}

	header, ciphertext, err := format.DecodeContainer(container)
	if err != nil {
		return err
	}

	kdfParams, err := header.GetKDFParams()
	if err != nil {
		return err
	}

	key, err := crypto.DeriveKey(opts.Password, kdfParams)
	if err != nil {
		return err
	}

	_, err = format.DecryptPayload(ciphertext, key, header)
	if err != nil {
		return fmt.Errorf("decryption failed: wrong password/key or corrupted file")
	}

	return nil
}
