package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/core-file-privacy/core-file-privacy/internal/archive"
	"github.com/core-file-privacy/core-file-privacy/internal/crypto"
	"github.com/core-file-privacy/core-file-privacy/internal/files"
	"github.com/core-file-privacy/core-file-privacy/internal/format"
)

type HideOptions struct {
	Path           string
	Password       string
	GenerateKey    bool
	Keep           bool
	DeleteOriginal bool
	Hidden         bool
	Output         string
	Archive        bool
	NameMode       string
	Profile        string
	Yes            bool
}

type HideResult struct {
	OutputPath   string
	GeneratedKey string
	OriginalKept bool
}

func Hide(opts HideOptions) (*HideResult, error) {
	info, err := os.Stat(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	if info.IsDir() && !opts.Archive {
		return nil, fmt.Errorf("the provided path is a directory.\n\nUse:\n  cfp hide %s --archive\n\nThis will create a single encrypted .cfp container.", opts.Path)
	}

	profile := crypto.Profile(opts.Profile)
	if profile == "" {
		profile = crypto.ProfileDefault
	}

	kdfParams, err := crypto.GetProfileParams(profile)
	if err != nil {
		return nil, err
	}

	salt, err := crypto.GenerateSalt(16)
	if err != nil {
		return nil, err
	}
	kdfParams = kdfParams.WithSalt(salt)

	nonce, err := crypto.GenerateNonce(crypto.NonceSize)
	if err != nil {
		return nil, err
	}

	password := opts.Password
	var generatedKey string

	if opts.GenerateKey {
		generatedKey, err = crypto.GenerateKey()
		if err != nil {
			return nil, err
		}
		password = generatedKey
	}

	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	var data []byte
	var metadata *format.Metadata

	if info.IsDir() {
		data, err = archive.CreateArchive(opts.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to create archive: %w", err)
		}
		metadata = format.NewDirectoryMetadata(filepath.Base(opts.Path))
	} else {
		data, err = os.ReadFile(opts.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		metadata = format.NewFileMetadata(
			filepath.Base(opts.Path),
			fmt.Sprintf("%04o", info.Mode().Perm()),
			info.Size(),
		)
	}

	key, err := crypto.DeriveKey(password, kdfParams)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	fileType := "file"
	if info.IsDir() {
		fileType = "archive"
	}

	header := format.NewHeader(fileType, profile, kdfParams, nonce)

	container, err := format.EncodeContainer(header, metadata, data, key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %w", err)
	}

	outputPath := opts.Output
	if outputPath == "" {
		outputPath = opts.Path + ".cfp"
	}

	if opts.NameMode == "random" {
		dir := filepath.Dir(outputPath)
		outputPath = filepath.Join(dir, generateRandomName()+".cfp")
	}

	if err := files.AtomicWrite(outputPath, container, 0644); err != nil {
		return nil, fmt.Errorf("failed to write container: %w", err)
	}

	if opts.Hidden {
		if err := files.HideFile(outputPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to hide file: %v\n", err)
		} else {
			outputPath = files.GetHiddenName(outputPath)
		}
	}

	originalKept := true
	if opts.DeleteOriginal || (!opts.Keep && !opts.Yes) {
		if opts.DeleteOriginal || opts.Yes {
			if info.IsDir() {
				if err := files.SecureRemoveAll(opts.Path); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to delete original: %v\n", err)
					originalKept = true
				} else {
					originalKept = false
				}
			} else {
				if err := files.SecureRemove(opts.Path); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to delete original: %v\n", err)
					originalKept = true
				} else {
					originalKept = false
				}
			}
		}
	}

	return &HideResult{
		OutputPath:   outputPath,
		GeneratedKey: generatedKey,
		OriginalKept: originalKept,
	}, nil
}

func generateRandomName() string {
	bytes, _ := crypto.GenerateRandomBytes(8)
	return fmt.Sprintf("%x", bytes)
}
