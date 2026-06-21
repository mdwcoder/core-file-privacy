package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	data := []byte("test content")
	if err := AtomicWrite(testFile, data, 0644); err != nil {
		t.Fatalf("AtomicWrite failed: %v", err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("content mismatch: got %s, want test content", content)
	}

	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	if info.Mode().Perm() != 0644 {
		t.Errorf("permission mismatch: got %o, want 0644", info.Mode().Perm())
	}
}

func TestAtomicWriteOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := AtomicWrite(testFile, []byte("first"), 0644); err != nil {
		t.Fatalf("first AtomicWrite failed: %v", err)
	}

	if err := AtomicWrite(testFile, []byte("second"), 0644); err != nil {
		t.Fatalf("second AtomicWrite failed: %v", err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(content) != "second" {
		t.Errorf("content should be overwritten: got %s, want second", content)
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if FileExists(testFile) {
		t.Error("file should not exist yet")
	}

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if !FileExists(testFile) {
		t.Error("file should exist now")
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "new", "nested", "dir")

	if err := EnsureDir(newDir); err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}

	info, err := os.Stat(newDir)
	if err != nil {
		t.Fatalf("directory should exist: %v", err)
	}

	if !info.IsDir() {
		t.Error("should be a directory")
	}
}

func TestSecureRemove(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := SecureRemove(testFile); err != nil {
		t.Fatalf("SecureRemove failed: %v", err)
	}

	if FileExists(testFile) {
		t.Error("file should be removed")
	}
}

func TestSecureRemoveAll(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "testdir")

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}

	testFile := filepath.Join(testDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := SecureRemoveAll(testDir); err != nil {
		t.Fatalf("SecureRemoveAll failed: %v", err)
	}

	if FileExists(testDir) {
		t.Error("directory should be removed")
	}
}
