package app

import (
	"fmt"
	"os"

	"github.com/core-file-privacy/core-file-privacy/internal/crypto"
	"github.com/core-file-privacy/core-file-privacy/internal/files"
	"github.com/core-file-privacy/core-file-privacy/internal/format"
)

type RekeyOptions struct {
	Path          string
	OldPassword   string
	NewPassword   string
	GenerateKey   bool
	Profile       string
	Yes           bool
}

type RekeyResult struct {
	GeneratedKey string
}

func Rekey(opts RekeyOptions) (*RekeyResult, error) {
	if opts.OldPassword == "" {
		return nil, fmt.Errorf("old password is required")
	}

	container, err := os.ReadFile(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read container: %w", err)
	}

	header, ciphertext, err := format.DecodeContainer(container)
	if err != nil {
		return nil, err
	}

	currentKDFParams, err := header.GetKDFParams()
	if err != nil {
		return nil, err
	}

	currentKey, err := crypto.DeriveKey(opts.OldPassword, currentKDFParams)
	if err != nil {
		return nil, err
	}

	payload, err := format.DecryptPayload(ciphertext, currentKey, header)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: wrong password/key or corrupted file")
	}

	metadata, data, err := format.ExtractMetadata(payload)
	if err != nil {
		return nil, err
	}

	profile := crypto.Profile(opts.Profile)
	if profile == "" {
		profile = crypto.Profile(header.Profile)
	}

	newKDFParams, err := crypto.GetProfileParams(profile)
	if err != nil {
		return nil, err
	}

	newSalt, err := crypto.GenerateSalt(16)
	if err != nil {
		return nil, err
	}
	newKDFParams = newKDFParams.WithSalt(newSalt)

	newNonce, err := crypto.GenerateNonce(crypto.NonceSize)
	if err != nil {
		return nil, err
	}

	newPassword := opts.NewPassword
	var generatedKey string

	if opts.GenerateKey {
		generatedKey, err = crypto.GenerateKey()
		if err != nil {
			return nil, err
		}
		newPassword = generatedKey
	}

	if newPassword == "" {
		return nil, fmt.Errorf("new password is required")
	}

	newKey, err := crypto.DeriveKey(newPassword, newKDFParams)
	if err != nil {
		return nil, err
	}

	newHeader := format.NewHeader(header.Type, profile, newKDFParams, newNonce)

	newContainer, err := format.EncodeContainer(newHeader, metadata, data, newKey)
	if err != nil {
		return nil, err
	}

	if err := files.AtomicWrite(opts.Path, newContainer, 0644); err != nil {
		return nil, err
	}

	return &RekeyResult{
		GeneratedKey: generatedKey,
	}, nil
}
