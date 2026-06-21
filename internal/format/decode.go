package format

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/core-file-privacy/core-file-privacy/internal/crypto"
)

func DecodeContainer(data []byte) (*CFPHeader, []byte, error) {
	if len(data) < len(MagicBytes)+4 {
		return nil, nil, fmt.Errorf("invalid container: too short")
	}

	magic := string(data[:len(MagicBytes)])
	if magic != MagicBytes {
		return nil, nil, fmt.Errorf("invalid container: bad magic bytes")
	}

	offset := len(MagicBytes)

	var headerLen uint32
	if err := binary.Read(bytes.NewReader(data[offset:offset+4]), binary.BigEndian, &headerLen); err != nil {
		return nil, nil, fmt.Errorf("failed to read header length: %w", err)
	}
	offset += 4

	if uint32(len(data)-offset) < headerLen {
		return nil, nil, fmt.Errorf("invalid container: header length exceeds data")
	}

	headerJSON := data[offset : offset+int(headerLen)]
	offset += int(headerLen)

	header, err := HeaderFromJSON(headerJSON)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse header: %w", err)
	}

	if err := header.Validate(); err != nil {
		return nil, nil, fmt.Errorf("invalid header: %w", err)
	}

	ciphertext := data[offset:]

	return header, ciphertext, nil
}

func DecryptPayload(ciphertext []byte, key []byte, header *CFPHeader) ([]byte, error) {
	nonce, err := header.GetNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	plaintext, err := crypto.Decrypt(ciphertext, key, nonce)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func ExtractMetadata(payload []byte) (*Metadata, []byte, error) {
	if len(payload) < 4 {
		return nil, nil, fmt.Errorf("invalid payload: too short")
	}

	var metadataLen uint32
	if err := binary.Read(bytes.NewReader(payload[:4]), binary.BigEndian, &metadataLen); err != nil {
		return nil, nil, fmt.Errorf("failed to read metadata length: %w", err)
	}

	if uint32(len(payload)-4) < metadataLen {
		return nil, nil, fmt.Errorf("invalid payload: metadata length exceeds data")
	}

	metadataJSON := payload[4 : 4+metadataLen]
	data := payload[4+metadataLen:]

	metadata, err := MetadataFromJSON(metadataJSON)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return metadata, data, nil
}
