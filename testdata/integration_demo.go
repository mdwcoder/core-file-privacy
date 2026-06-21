package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/core-file-privacy/core-file-privacy/internal/archive"
	"github.com/core-file-privacy/core-file-privacy/internal/crypto"
	"github.com/core-file-privacy/core-file-privacy/internal/files"
	"github.com/core-file-privacy/core-file-privacy/internal/format"
)

func main() {
	fmt.Println("=== Core File Privacy Integration Test ===")
	fmt.Println()

	tmpDir := "/tmp/cfp-integration-test"
	os.RemoveAll(tmpDir)
	os.MkdirAll(tmpDir, 0755)

	// Test 1: Generate secure key
	fmt.Println("Test 1: Generate secure key")
	key, err := crypto.GenerateKey()
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Generated key: %s\n", key)
	if !crypto.ValidateKey(key) {
		fmt.Println("FAIL: Generated key is not valid")
		os.Exit(1)
	}
	fmt.Println("PASS: Key generated and validated")
	fmt.Println()

	// Test 2: Encrypt a file
	fmt.Println("Test 2: Encrypt a file")
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("This is secret content for testing")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	salt, _ := crypto.GenerateSalt(16)
	nonce, _ := crypto.GenerateNonce(12)
	params := crypto.KDFParams{
		MemoryKib:   32768,
		Iterations:  2,
		Parallelism: 2,
		Salt:        salt,
	}

	derivedKey, err := crypto.DeriveKey(key, params)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	header := format.NewHeader("file", crypto.ProfileFast, params, nonce)
	metadata := format.NewFileMetadata("test.txt", "0644", int64(len(testContent)))

	container, err := format.EncodeContainer(header, metadata, testContent, derivedKey)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	encryptedFile := filepath.Join(tmpDir, "test.txt.cfp")
	if err := files.AtomicWrite(encryptedFile, container, 0644); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Encrypted file created: %s\n", encryptedFile)
	fmt.Println("PASS: File encrypted successfully")
	fmt.Println()

	// Test 3: Decrypt the file
	fmt.Println("Test 3: Decrypt the file")
	containerData, err := os.ReadFile(encryptedFile)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	decodedHeader, ciphertext, err := format.DecodeContainer(containerData)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	decodedParams, err := decodedHeader.GetKDFParams()
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	decryptedKey, err := crypto.DeriveKey(key, decodedParams)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	payload, err := format.DecryptPayload(ciphertext, decryptedKey, decodedHeader)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	decodedMetadata, decryptedData, err := format.ExtractMetadata(payload)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	if string(decryptedData) != string(testContent) {
		fmt.Printf("FAIL: Decrypted content doesn't match\n")
		os.Exit(1)
	}

	if decodedMetadata.OriginalName != "test.txt" {
		fmt.Printf("FAIL: Metadata name doesn't match\n")
		os.Exit(1)
	}
	fmt.Println("PASS: File decrypted successfully")
	fmt.Println()

	// Test 4: Wrong password
	fmt.Println("Test 4: Wrong password should fail")
	wrongKey, _ := crypto.DeriveKey("wrongpassword", decodedParams)
	_, err = format.DecryptPayload(ciphertext, wrongKey, decodedHeader)
	if err == nil {
		fmt.Println("FAIL: Decryption with wrong key should fail")
		os.Exit(1)
	}
	fmt.Println("PASS: Wrong password correctly rejected")
	fmt.Println()

	// Test 5: Archive directory
	fmt.Println("Test 5: Archive and encrypt directory")
	testDir := filepath.Join(tmpDir, "testdir")
	os.MkdirAll(testDir, 0755)
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("content2"), 0644)

	archiveData, err := archive.CreateArchive(testDir)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	dirHeader := format.NewHeader("archive", crypto.ProfileFast, params, nonce)
	dirMetadata := format.NewDirectoryMetadata("testdir")

	dirContainer, err := format.EncodeContainer(dirHeader, dirMetadata, archiveData, derivedKey)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	dirEncrypted := filepath.Join(tmpDir, "testdir.cfp")
	if err := files.AtomicWrite(dirEncrypted, dirContainer, 0644); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Encrypted directory: %s\n", dirEncrypted)
	fmt.Println("PASS: Directory archived and encrypted")
	fmt.Println()

	// Test 6: Extract archive
	fmt.Println("Test 6: Decrypt and extract directory")
	dirContainerData, _ := os.ReadFile(dirEncrypted)
	_, dirCiphertext, _ := format.DecodeContainer(dirContainerData)
	dirPayload, _ := format.DecryptPayload(dirCiphertext, derivedKey, decodedHeader)
	_, dirArchiveData, _ := format.ExtractMetadata(dirPayload)

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := archive.ExtractArchive(dirArchiveData, extractDir, false); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	extractedFile1 := filepath.Join(extractDir, "testdir", "file1.txt")
	content1, _ := os.ReadFile(extractedFile1)
	if string(content1) != "content1" {
		fmt.Println("FAIL: Extracted content doesn't match")
		os.Exit(1)
	}
	fmt.Println("PASS: Directory extracted successfully")
	fmt.Println()

	// Test 7: Hidden file
	fmt.Println("Test 7: Hide file (Unix)")
	hiddenFile := filepath.Join(tmpDir, "visible.txt")
	os.WriteFile(hiddenFile, []byte("test"), 0644)
	if err := files.HideFile(hiddenFile); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	hiddenPath := filepath.Join(tmpDir, ".visible.txt")
	if !files.FileExists(hiddenPath) {
		fmt.Println("FAIL: Hidden file not found")
		os.Exit(1)
	}
	if !files.IsHidden(hiddenPath) {
		fmt.Println("FAIL: File should be hidden")
		os.Exit(1)
	}
	fmt.Println("PASS: File hidden successfully")
	fmt.Println()

	// Test 8: Info without password
	fmt.Println("Test 8: Read container info without password")
	infoContainer, _ := os.ReadFile(encryptedFile)
	infoHeader, _, err := format.DecodeContainer(infoContainer)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	if infoHeader.Format != "CFP" {
		fmt.Println("FAIL: Invalid format")
		os.Exit(1)
	}
	if infoHeader.Cipher != "AES-256-GCM" {
		fmt.Println("FAIL: Invalid cipher")
		os.Exit(1)
	}
	fmt.Printf("Container info: Format=%s, Cipher=%s, KDF=%s, Profile=%s\n",
		infoHeader.Format, infoHeader.Cipher, infoHeader.KDF, infoHeader.Profile)
	fmt.Println("PASS: Info read successfully")
	fmt.Println()

	fmt.Println("=== All Tests Passed ===")
	fmt.Println()
	fmt.Printf("Test files are in: %s\n", tmpDir)
	fmt.Println("You can safely delete this directory.")
}
