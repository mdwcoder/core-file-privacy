package archive

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAndExtractArchive(t *testing.T) {
	tmpDir := t.TempDir()

	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}

	testFile := filepath.Join(sourceDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	subDir := filepath.Join(sourceDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	subFile := filepath.Join(subDir, "nested.txt")
	if err := os.WriteFile(subFile, []byte("nested content"), 0644); err != nil {
		t.Fatalf("failed to write nested file: %v", err)
	}

	archiveData, err := CreateArchive(sourceDir)
	if err != nil {
		t.Fatalf("CreateArchive failed: %v", err)
	}

	if len(archiveData) == 0 {
		t.Error("archive should not be empty")
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := ExtractArchive(archiveData, extractDir, false); err != nil {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	extractedFile := filepath.Join(extractDir, "source", "test.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("content mismatch: got %s, want test content", content)
	}

	extractedNested := filepath.Join(extractDir, "source", "subdir", "nested.txt")
	nestedContent, err := os.ReadFile(extractedNested)
	if err != nil {
		t.Fatalf("failed to read extracted nested file: %v", err)
	}

	if string(nestedContent) != "nested content" {
		t.Errorf("nested content mismatch: got %s, want nested content", nestedContent)
	}
}

func TestExtractArchivePathTraversal(t *testing.T) {
	tmpDir := t.TempDir()

	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}

	testFile := filepath.Join(sourceDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	archiveData, err := CreateArchive(sourceDir)
	if err != nil {
		t.Fatalf("CreateArchive failed: %v", err)
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	err = ExtractArchive(archiveData, extractDir, false)
	if err != nil {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	parentDir := filepath.Dir(tmpDir)
	escapedFile := filepath.Join(parentDir, "escaped.txt")
	if _, err := os.Stat(escapedFile); err == nil {
		t.Error("path traversal should not write outside extract directory")
	}
}

func TestExtractArchiveNoOverwrite(t *testing.T) {
	tmpDir := t.TempDir()

	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}

	testFile := filepath.Join(sourceDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("original"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	archiveData, err := CreateArchive(sourceDir)
	if err != nil {
		t.Fatalf("CreateArchive failed: %v", err)
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := ExtractArchive(archiveData, extractDir, false); err != nil {
		t.Fatalf("first ExtractArchive failed: %v", err)
	}

	err = ExtractArchive(archiveData, extractDir, false)
	if err == nil {
		t.Error("second extract without force should fail")
	}

	err = ExtractArchive(archiveData, extractDir, true)
	if err != nil {
		t.Errorf("extract with force should succeed: %v", err)
	}
}

func TestCreateArchiveNotDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := CreateArchive(testFile)
	if err == nil {
		t.Error("CreateArchive on file should fail")
	}
}
