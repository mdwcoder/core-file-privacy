package format

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/core-file-privacy/core-file-privacy/internal/crypto"
)

const (
	MagicBytes = "CFP1"
	Version    = 1
)

type CFPHeader struct {
	Format           string          `json:"format"`
	Version          int             `json:"version"`
	Type             string          `json:"type"`
	Cipher           string          `json:"cipher"`
	KDF              string          `json:"kdf"`
	Profile          string          `json:"profile"`
	KDFParams        KDFParamsJSON   `json:"kdf_params"`
	Nonce            string          `json:"nonce"`
	CreatedAt        string          `json:"created_at"`
	MetadataEncrypted bool           `json:"metadata_encrypted"`
}

type KDFParamsJSON struct {
	MemoryKib   uint32 `json:"memory_kib"`
	Iterations  uint32 `json:"iterations"`
	Parallelism uint8  `json:"parallelism"`
	Salt        string `json:"salt"`
}

func NewHeader(fileType string, profile crypto.Profile, kdfParams crypto.KDFParams, nonce []byte) *CFPHeader {
	return &CFPHeader{
		Format:   "CFP",
		Version:  Version,
		Type:     fileType,
		Cipher:   "AES-256-GCM",
		KDF:      "Argon2id",
		Profile:  string(profile),
		KDFParams: KDFParamsJSON{
			MemoryKib:   kdfParams.MemoryKib,
			Iterations:  kdfParams.Iterations,
			Parallelism: kdfParams.Parallelism,
			Salt:        base64.RawURLEncoding.EncodeToString(kdfParams.Salt),
		},
		Nonce:             base64.RawURLEncoding.EncodeToString(nonce),
		CreatedAt:         time.Now().UTC().Format(time.RFC3339),
		MetadataEncrypted: true,
	}
}

func (h *CFPHeader) Validate() error {
	if h.Format != "CFP" {
		return fmt.Errorf("invalid format: %s", h.Format)
	}
	if h.Version != Version {
		return fmt.Errorf("unsupported version: %d", h.Version)
	}
	if h.Cipher != "AES-256-GCM" {
		return fmt.Errorf("unsupported cipher: %s", h.Cipher)
	}
	if h.KDF != "Argon2id" {
		return fmt.Errorf("unsupported KDF: %s", h.KDF)
	}
	return nil
}

func (h *CFPHeader) GetKDFParams() (crypto.KDFParams, error) {
	salt, err := base64.RawURLEncoding.DecodeString(h.KDFParams.Salt)
	if err != nil {
		return crypto.KDFParams{}, fmt.Errorf("failed to decode salt: %w", err)
	}

	return crypto.KDFParams{
		MemoryKib:   h.KDFParams.MemoryKib,
		Iterations:  h.KDFParams.Iterations,
		Parallelism: h.KDFParams.Parallelism,
		Salt:        salt,
	}, nil
}

func (h *CFPHeader) GetNonce() ([]byte, error) {
	nonce, err := base64.RawURLEncoding.DecodeString(h.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %w", err)
	}
	return nonce, nil
}

func (h *CFPHeader) ToJSON() ([]byte, error) {
	return json.Marshal(h)
}

func HeaderFromJSON(data []byte) (*CFPHeader, error) {
	var header CFPHeader
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, fmt.Errorf("failed to parse header: %w", err)
	}
	return &header, nil
}
