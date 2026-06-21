package app

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/core-file-privacy/core-file-privacy/internal/crypto"
	"github.com/core-file-privacy/core-file-privacy/internal/format"
)

type InfoOptions struct {
	Path     string
	Password string
	Unlock   bool
}

type InfoResult struct {
	Header   *format.CFPHeader
	Metadata *format.Metadata
}

func Info(opts InfoOptions) (*InfoResult, error) {
	container, err := os.ReadFile(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read container: %w", err)
	}

	header, ciphertext, err := format.DecodeContainer(container)
	if err != nil {
		return nil, err
	}

	result := &InfoResult{
		Header: header,
	}

	if opts.Unlock {
		if opts.Password == "" {
			return nil, fmt.Errorf("password is required for --unlock")
		}

		kdfParams, err := header.GetKDFParams()
		if err != nil {
			return nil, err
		}

		key, err := crypto.DeriveKey(opts.Password, kdfParams)
		if err != nil {
			return nil, err
		}

		payload, err := format.DecryptPayload(ciphertext, key, header)
		if err != nil {
			return nil, fmt.Errorf("decryption failed: wrong password/key or corrupted file")
		}

		metadata, _, err := format.ExtractMetadata(payload)
		if err != nil {
			return nil, err
		}

		result.Metadata = metadata
	}

	return result, nil
}

func InfoJSON(opts InfoOptions) (string, error) {
	container, err := os.ReadFile(opts.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read container: %w", err)
	}

	header, _, err := format.DecodeContainer(container)
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(header, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
