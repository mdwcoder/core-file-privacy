package format

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/core-file-privacy/core-file-privacy/internal/crypto"
)

func EncodeContainer(header *CFPHeader, metadata *Metadata, data []byte, key []byte) ([]byte, error) {
	metadataJSON, err := metadata.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize metadata: %w", err)
	}

	nonce, err := header.GetNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	payload := createPayload(metadataJSON, data)

	ciphertext, err := crypto.Encrypt(payload, key, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	headerJSON, err := header.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize header: %w", err)
	}

	container := createContainer(headerJSON, ciphertext)

	return container, nil
}

func createPayload(metadataJSON, data []byte) []byte {
	var buf bytes.Buffer

	metadataLen := uint32(len(metadataJSON))
	binary.Write(&buf, binary.BigEndian, metadataLen)
	buf.Write(metadataJSON)
	buf.Write(data)

	return buf.Bytes()
}

func createContainer(headerJSON, ciphertext []byte) []byte {
	var buf bytes.Buffer

	buf.WriteString(MagicBytes)

	headerLen := uint32(len(headerJSON))
	binary.Write(&buf, binary.BigEndian, headerLen)
	buf.Write(headerJSON)
	buf.Write(ciphertext)

	return buf.Bytes()
}
