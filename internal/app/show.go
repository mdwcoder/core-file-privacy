package app

import (
	"fmt"
	"os"

	"github.com/core-file-privacy/core-file-privacy/internal/archive"
	"github.com/core-file-privacy/core-file-privacy/internal/crypto"
	"github.com/core-file-privacy/core-file-privacy/internal/files"
	"github.com/core-file-privacy/core-file-privacy/internal/format"
)

type ShowOptions struct {
	Path            string
	Password        string
	Output          string
	Keep            bool
	DeleteEncrypted bool
	Force           bool
	Yes             bool
}

type ShowResult struct {
	OutputPath string
}

func Show(opts ShowOptions) (*ShowResult, error) {
	if opts.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	container, err := os.ReadFile(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read container: %w", err)
	}

	header, ciphertext, err := format.DecodeContainer(container)
	if err != nil {
		return nil, err
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

	metadata, data, err := format.ExtractMetadata(payload)
	if err != nil {
		return nil, err
	}

	outputPath := opts.Output
	if outputPath == "" {
		outputPath = metadata.OriginalName
	}

	if !opts.Force {
		if _, err := os.Stat(outputPath); err == nil {
			return nil, fmt.Errorf("output path already exists: %s (use --force to overwrite)", outputPath)
		}
	}

	if metadata.IsDirectory() {
		if err := archive.ExtractArchiveWithStrip(data, outputPath, opts.Force, 1); err != nil {
			return nil, err
		}
	} else {
		mode := os.FileMode(0644)
		if metadata.OriginalMode != "" {
			var m uint32
			fmt.Sscanf(metadata.OriginalMode, "%o", &m)
			mode = os.FileMode(m)
		}

		if err := files.AtomicWrite(outputPath, data, mode); err != nil {
			return nil, err
		}
	}

	if opts.DeleteEncrypted {
		if err := files.SecureRemove(opts.Path); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to delete encrypted file: %v\n", err)
		}
	}

	return &ShowResult{
		OutputPath: outputPath,
	}, nil
}
