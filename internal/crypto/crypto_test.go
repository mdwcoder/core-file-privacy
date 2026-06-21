package crypto

import (
	"strings"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt(16)
	if err != nil {
		t.Fatalf("GenerateSalt failed: %v", err)
	}

	if len(salt) != 16 {
		t.Errorf("expected salt length 16, got %d", len(salt))
	}

	salt2, _ := GenerateSalt(16)
	if string(salt) == string(salt2) {
		t.Error("two salts should not be equal")
	}
}

func TestGenerateNonce(t *testing.T) {
	nonce, err := GenerateNonce(12)
	if err != nil {
		t.Fatalf("GenerateNonce failed: %v", err)
	}

	if len(nonce) != 12 {
		t.Errorf("expected nonce length 12, got %d", len(nonce))
	}
}

func TestGenerateKey(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if !strings.HasPrefix(key, "cfp_") {
		t.Errorf("key should start with cfp_, got %s", key)
	}

	parts := strings.Split(key[4:], "-")
	if len(parts) != 8 {
		t.Errorf("key should have 8 groups, got %d", len(parts))
	}

	for _, part := range parts {
		if len(part) != 4 {
			t.Errorf("each group should have 4 chars, got %d", len(part))
		}
	}

	if !ValidateKey(key) {
		t.Error("generated key should be valid")
	}
}

func TestValidateKey(t *testing.T) {
	tests := []struct {
		key   string
		valid bool
	}{
		{"cfp_ABCD-1234-EFGH-5678-IJKL-9012-MNOP-3456", true},
		{"cfp_ABCD-1234-EFGH-5678-IJKL-9012-MNOP-345", false},
		{"cfp_ABCD-1234-EFGH-5678-IJKL-9012-MNOP", false},
		{"ABC-ABCD-1234-EFGH-5678-IJKL-9012-MNOP-3456", false},
		{"cfp_abcd-1234-efgh-5678-ijkl-9012-mnop-3456", false},
		{"", false},
	}

	for _, tt := range tests {
		if ValidateKey(tt.key) != tt.valid {
			t.Errorf("ValidateKey(%s) = %v, want %v", tt.key, !tt.valid, tt.valid)
		}
	}
}

func TestDeriveKey(t *testing.T) {
	salt, _ := GenerateSalt(16)
	params := KDFParams{
		MemoryKib:   32768,
		Iterations:  2,
		Parallelism: 2,
		Salt:        salt,
	}

	key, err := DeriveKey("testpassword", params)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("expected key length 32, got %d", len(key))
	}

	key2, _ := DeriveKey("testpassword", params)
	if string(key) != string(key2) {
		t.Error("same password and params should produce same key")
	}

	key3, _ := DeriveKey("differentpassword", params)
	if string(key) == string(key3) {
		t.Error("different passwords should produce different keys")
	}
}

func TestDeriveKeyEmptyPassword(t *testing.T) {
	salt, _ := GenerateSalt(16)
	params := KDFParams{
		MemoryKib:   32768,
		Iterations:  2,
		Parallelism: 2,
		Salt:        salt,
	}

	_, err := DeriveKey("", params)
	if err == nil {
		t.Error("empty password should fail")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key, _ := GenerateRandomBytes(32)
	nonce, _ := GenerateNonce(12)
	plaintext := []byte("Hello, World!")

	ciphertext, err := Encrypt(plaintext, key, nonce)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key, nonce)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("decrypted text doesn't match: got %s, want %s", decrypted, plaintext)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1, _ := GenerateRandomBytes(32)
	key2, _ := GenerateRandomBytes(32)
	nonce, _ := GenerateNonce(12)
	plaintext := []byte("Hello, World!")

	ciphertext, _ := Encrypt(plaintext, key1, nonce)

	_, err := Decrypt(ciphertext, key2, nonce)
	if err == nil {
		t.Error("decrypt with wrong key should fail")
	}
}

func TestDecryptCorruptedData(t *testing.T) {
	key, _ := GenerateRandomBytes(32)
	nonce, _ := GenerateNonce(12)
	plaintext := []byte("Hello, World!")

	ciphertext, _ := Encrypt(plaintext, key, nonce)
	corrupted := make([]byte, len(ciphertext))
	copy(corrupted, ciphertext)
	corrupted[0] ^= 0xFF

	_, err := Decrypt(corrupted, key, nonce)
	if err == nil {
		t.Error("decrypt corrupted data should fail")
	}
}

func TestEncryptInvalidKeySize(t *testing.T) {
	key := make([]byte, 16)
	nonce := make([]byte, 12)

	_, err := Encrypt([]byte("test"), key, nonce)
	if err == nil {
		t.Error("encrypt with invalid key size should fail")
	}
}
