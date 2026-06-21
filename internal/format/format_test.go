package format

import (
	"testing"

	"github.com/core-file-privacy/core-file-privacy/internal/crypto"
)

func TestNewHeader(t *testing.T) {
	salt, _ := crypto.GenerateSalt(16)
	nonce, _ := crypto.GenerateNonce(12)
	params := crypto.KDFParams{
		MemoryKib:   65536,
		Iterations:  3,
		Parallelism: 4,
		Salt:        salt,
	}

	header := NewHeader("file", crypto.ProfileDefault, params, nonce)

	if header.Format != "CFP" {
		t.Errorf("expected format CFP, got %s", header.Format)
	}

	if header.Version != 1 {
		t.Errorf("expected version 1, got %d", header.Version)
	}

	if header.Cipher != "AES-256-GCM" {
		t.Errorf("expected cipher AES-256-GCM, got %s", header.Cipher)
	}

	if header.KDF != "Argon2id" {
		t.Errorf("expected KDF Argon2id, got %s", header.KDF)
	}

	if err := header.Validate(); err != nil {
		t.Errorf("header validation failed: %v", err)
	}
}

func TestHeaderSerialization(t *testing.T) {
	salt, _ := crypto.GenerateSalt(16)
	nonce, _ := crypto.GenerateNonce(12)
	params := crypto.KDFParams{
		MemoryKib:   65536,
		Iterations:  3,
		Parallelism: 4,
		Salt:        salt,
	}

	header := NewHeader("file", crypto.ProfileDefault, params, nonce)

	json, err := header.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	decoded, err := HeaderFromJSON(json)
	if err != nil {
		t.Fatalf("HeaderFromJSON failed: %v", err)
	}

	if decoded.Format != header.Format {
		t.Errorf("format mismatch: got %s, want %s", decoded.Format, header.Format)
	}

	if decoded.Version != header.Version {
		t.Errorf("version mismatch: got %d, want %d", decoded.Version, header.Version)
	}
}

func TestMetadataSerialization(t *testing.T) {
	metadata := NewFileMetadata("test.txt", "0644", 1234)

	json, err := metadata.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	decoded, err := MetadataFromJSON(json)
	if err != nil {
		t.Fatalf("MetadataFromJSON failed: %v", err)
	}

	if decoded.OriginalName != metadata.OriginalName {
		t.Errorf("name mismatch: got %s, want %s", decoded.OriginalName, metadata.OriginalName)
	}

	if decoded.OriginalType != metadata.OriginalType {
		t.Errorf("type mismatch: got %s, want %s", decoded.OriginalType, metadata.OriginalType)
	}

	if decoded.OriginalSize != metadata.OriginalSize {
		t.Errorf("size mismatch: got %d, want %d", decoded.OriginalSize, metadata.OriginalSize)
	}
}

func TestDirectoryMetadata(t *testing.T) {
	metadata := NewDirectoryMetadata("test-folder")

	if !metadata.IsDirectory() {
		t.Error("directory metadata should report IsDirectory as true")
	}

	if !metadata.Archived {
		t.Error("directory metadata should have Archived as true")
	}
}

func TestEncodeDecodeContainer(t *testing.T) {
	salt, _ := crypto.GenerateSalt(16)
	nonce, _ := crypto.GenerateNonce(12)
	params := crypto.KDFParams{
		MemoryKib:   32768,
		Iterations:  2,
		Parallelism: 2,
		Salt:        salt,
	}

	header := NewHeader("file", crypto.ProfileFast, params, nonce)
	metadata := NewFileMetadata("test.txt", "0644", 13)
	data := []byte("Hello, World!")

	key, _ := crypto.DeriveKey("testpassword", params)

	container, err := EncodeContainer(header, metadata, data, key)
	if err != nil {
		t.Fatalf("EncodeContainer failed: %v", err)
	}

	if string(container[:4]) != "CFP1" {
		t.Error("container should start with CFP1 magic bytes")
	}

	decodedHeader, ciphertext, err := DecodeContainer(container)
	if err != nil {
		t.Fatalf("DecodeContainer failed: %v", err)
	}

	if decodedHeader.Format != header.Format {
		t.Errorf("header format mismatch")
	}

	payload, err := DecryptPayload(ciphertext, key, decodedHeader)
	if err != nil {
		t.Fatalf("DecryptPayload failed: %v", err)
	}

	decodedMetadata, decodedData, err := ExtractMetadata(payload)
	if err != nil {
		t.Fatalf("ExtractMetadata failed: %v", err)
	}

	if decodedMetadata.OriginalName != metadata.OriginalName {
		t.Errorf("metadata name mismatch")
	}

	if string(decodedData) != string(data) {
		t.Errorf("data mismatch: got %s, want %s", decodedData, data)
	}
}

func TestDecodeContainerInvalidMagic(t *testing.T) {
	_, _, err := DecodeContainer([]byte("INVALID"))
	if err == nil {
		t.Error("decode with invalid magic should fail")
	}
}

func TestDecodeContainerTooShort(t *testing.T) {
	_, _, err := DecodeContainer([]byte("CF"))
	if err == nil {
		t.Error("decode with too short data should fail")
	}
}
