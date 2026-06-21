package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHideAndShowFile(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	content := []byte("This is secret content")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	password := "testpassword123"

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: password,
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	if hideResult.OutputPath == "" {
		t.Error("OutputPath should not be empty")
	}

	if !strings.HasSuffix(hideResult.OutputPath, ".cfp") {
		t.Errorf("OutputPath should end with .cfp, got %s", hideResult.OutputPath)
	}

	if _, err := os.Stat(hideResult.OutputPath); os.IsNotExist(err) {
		t.Error("Encrypted file should exist")
	}

	showResult, err := Show(ShowOptions{
		Path:     hideResult.OutputPath,
		Password: password,
		Output:   filepath.Join(tmpDir, "restored.txt"),
		Force:    true,
	})
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	restoredContent, err := os.ReadFile(showResult.OutputPath)
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}

	if string(restoredContent) != string(content) {
		t.Errorf("restored content doesn't match: got %q, want %q", restoredContent, content)
	}
}

func TestHideAndShowWithGeneratedKey(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	content := []byte("Secret with generated key")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:        testFile,
		GenerateKey: true,
		Keep:        true,
		Profile:     "fast",
		Yes:         true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	if hideResult.GeneratedKey == "" {
		t.Error("GeneratedKey should not be empty")
	}

	if !strings.HasPrefix(hideResult.GeneratedKey, "cfp_") {
		t.Errorf("GeneratedKey should start with cfp_, got %s", hideResult.GeneratedKey)
	}

	showResult, err := Show(ShowOptions{
		Path:     hideResult.OutputPath,
		Password: hideResult.GeneratedKey,
		Output:   filepath.Join(tmpDir, "restored.txt"),
		Force:    true,
	})
	if err != nil {
		t.Fatalf("Show with generated key failed: %v", err)
	}

	restoredContent, err := os.ReadFile(showResult.OutputPath)
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}

	if string(restoredContent) != string(content) {
		t.Errorf("restored content doesn't match")
	}
}

func TestVerifyCorrectPassword(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	password := "testpassword123"

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: password,
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	err = Verify(VerifyOptions{
		Path:     hideResult.OutputPath,
		Password: password,
	})
	if err != nil {
		t.Errorf("Verify with correct password should succeed: %v", err)
	}
}

func TestVerifyWrongPassword(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	password := "testpassword123"

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: password,
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	err = Verify(VerifyOptions{
		Path:     hideResult.OutputPath,
		Password: "wrongpassword",
	})
	if err == nil {
		t.Error("Verify with wrong password should fail")
	}
}

func TestInfoWithoutPassword(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: "testpassword",
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	infoResult, err := Info(InfoOptions{
		Path: hideResult.OutputPath,
	})
	if err != nil {
		t.Fatalf("Info failed: %v", err)
	}

	if infoResult.Header == nil {
		t.Error("Header should not be nil")
	}

	if infoResult.Header.Format != "CFP" {
		t.Errorf("Header format should be CFP, got %s", infoResult.Header.Format)
	}

	if infoResult.Header.Cipher != "AES-256-GCM" {
		t.Errorf("Header cipher should be AES-256-GCM, got %s", infoResult.Header.Cipher)
	}

	if infoResult.Metadata != nil {
		t.Error("Metadata should be nil without --unlock")
	}
}

func TestInfoWithUnlock(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	password := "testpassword123"

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: password,
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	infoResult, err := Info(InfoOptions{
		Path:     hideResult.OutputPath,
		Password: password,
		Unlock:   true,
	})
	if err != nil {
		t.Fatalf("Info with unlock failed: %v", err)
	}

	if infoResult.Metadata == nil {
		t.Error("Metadata should not be nil with --unlock")
	}

	if infoResult.Metadata.OriginalName != "secret.txt" {
		t.Errorf("OriginalName should be secret.txt, got %s", infoResult.Metadata.OriginalName)
	}

	if infoResult.Metadata.OriginalType != "file" {
		t.Errorf("OriginalType should be file, got %s", infoResult.Metadata.OriginalType)
	}
}

func TestInfoWithUnlockWrongPassword(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: "correctpassword",
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	_, err = Info(InfoOptions{
		Path:     hideResult.OutputPath,
		Password: "wrongpassword",
		Unlock:   true,
	})
	if err == nil {
		t.Error("Info with wrong password should fail")
	}
}

func TestRekeyWithNewPassword(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	content := []byte("Secret content for rekey")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	oldPassword := "oldpassword123"
	newPassword := "newpassword456"

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: oldPassword,
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	rekeyResult, err := Rekey(RekeyOptions{
		Path:        hideResult.OutputPath,
		OldPassword: oldPassword,
		NewPassword: newPassword,
		Yes:         true,
	})
	if err != nil {
		t.Fatalf("Rekey failed: %v", err)
	}

	if rekeyResult.GeneratedKey != "" {
		t.Error("GeneratedKey should be empty when not using --generate-key")
	}

	err = Verify(VerifyOptions{
		Path:     hideResult.OutputPath,
		Password: newPassword,
	})
	if err != nil {
		t.Errorf("Verify with new password should succeed: %v", err)
	}

	err = Verify(VerifyOptions{
		Path:     hideResult.OutputPath,
		Password: oldPassword,
	})
	if err == nil {
		t.Error("Verify with old password should fail after rekey")
	}

	showResult, err := Show(ShowOptions{
		Path:     hideResult.OutputPath,
		Password: newPassword,
		Output:   filepath.Join(tmpDir, "restored.txt"),
		Force:    true,
	})
	if err != nil {
		t.Fatalf("Show with new password failed: %v", err)
	}

	restoredContent, err := os.ReadFile(showResult.OutputPath)
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}

	if string(restoredContent) != string(content) {
		t.Errorf("restored content doesn't match after rekey")
	}
}

func TestRekeyWithGeneratedKey(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	oldPassword := "oldpassword123"

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: oldPassword,
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	rekeyResult, err := Rekey(RekeyOptions{
		Path:        hideResult.OutputPath,
		OldPassword: oldPassword,
		GenerateKey: true,
		Yes:         true,
	})
	if err != nil {
		t.Fatalf("Rekey with generated key failed: %v", err)
	}

	if rekeyResult.GeneratedKey == "" {
		t.Error("GeneratedKey should not be empty with --generate-key")
	}

	err = Verify(VerifyOptions{
		Path:     hideResult.OutputPath,
		Password: rekeyResult.GeneratedKey,
	})
	if err != nil {
		t.Errorf("Verify with generated key should succeed: %v", err)
	}
}

func TestHideDirectoryWithoutArchive(t *testing.T) {
	tmpDir := t.TempDir()

	testDir := filepath.Join(tmpDir, "testdir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}

	_, err := Hide(HideOptions{
		Path:     testDir,
		Password: "testpassword",
		Archive:  false,
		Yes:      true,
	})
	if err == nil {
		t.Error("Hide directory without --archive should fail")
	}
}

func TestHideAndShowDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	testDir := filepath.Join(tmpDir, "private-folder")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}

	file1 := filepath.Join(testDir, "file1.txt")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("failed to write file1: %v", err)
	}

	subDir := filepath.Join(testDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	file2 := filepath.Join(subDir, "file2.txt")
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatalf("failed to write file2: %v", err)
	}

	password := "testpassword123"

	hideResult, err := Hide(HideOptions{
		Path:     testDir,
		Password: password,
		Archive:  true,
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide directory failed: %v", err)
	}

	if _, err := os.Stat(hideResult.OutputPath); os.IsNotExist(err) {
		t.Error("Encrypted directory file should exist")
	}

	infoResult, err := Info(InfoOptions{
		Path:     hideResult.OutputPath,
		Password: password,
		Unlock:   true,
	})
	if err != nil {
		t.Fatalf("Info failed: %v", err)
	}

	if infoResult.Metadata.OriginalType != "directory" {
		t.Errorf("OriginalType should be directory, got %s", infoResult.Metadata.OriginalType)
	}

	if !infoResult.Metadata.Archived {
		t.Error("Archived should be true for directory")
	}

	showResult, err := Show(ShowOptions{
		Path:     hideResult.OutputPath,
		Password: password,
		Output:   filepath.Join(tmpDir, "restored-folder"),
		Force:    true,
	})
	if err != nil {
		t.Fatalf("Show directory failed: %v", err)
	}

	restoredFile1 := filepath.Join(showResult.OutputPath, "file1.txt")
	content1, err := os.ReadFile(restoredFile1)
	if err != nil {
		t.Fatalf("failed to read restored file1: %v", err)
	}
	if string(content1) != "content1" {
		t.Errorf("file1 content doesn't match")
	}

	restoredFile2 := filepath.Join(showResult.OutputPath, "subdir", "file2.txt")
	content2, err := os.ReadFile(restoredFile2)
	if err != nil {
		t.Fatalf("failed to read restored file2: %v", err)
	}
	if string(content2) != "content2" {
		t.Errorf("file2 content doesn't match")
	}
}

func TestHideWithHidden(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "visible.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: "testpassword",
		Hidden:   true,
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide with hidden failed: %v", err)
	}

	if !strings.HasPrefix(filepath.Base(hideResult.OutputPath), ".") {
		t.Errorf("Hidden file should start with dot, got %s", hideResult.OutputPath)
	}

	if _, err := os.Stat(hideResult.OutputPath); os.IsNotExist(err) {
		t.Error("Hidden file should exist")
	}
}

func TestHideDeleteOriginal(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "delete-me.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:           testFile,
		Password:       "testpassword",
		DeleteOriginal: true,
		Profile:        "fast",
		Yes:            true,
	})
	if err != nil {
		t.Fatalf("Hide with delete-original failed: %v", err)
	}

	if _, err := os.Stat(hideResult.OutputPath); os.IsNotExist(err) {
		t.Error("Encrypted file should exist")
	}

	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Original file should be deleted")
	}
}

func TestShowDeleteEncrypted(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: "testpassword",
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	_, err = Show(ShowOptions{
		Path:            hideResult.OutputPath,
		Password:        "testpassword",
		Output:          filepath.Join(tmpDir, "restored.txt"),
		DeleteEncrypted: true,
		Force:           true,
	})
	if err != nil {
		t.Fatalf("Show with delete-encrypted failed: %v", err)
	}

	if _, err := os.Stat(hideResult.OutputPath); !os.IsNotExist(err) {
		t.Error("Encrypted file should be deleted")
	}
}

func TestShowWithoutForce(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: "testpassword",
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "restored.txt")
	if err := os.WriteFile(outputPath, []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to write existing file: %v", err)
	}

	_, err = Show(ShowOptions{
		Path:     hideResult.OutputPath,
		Password: "testpassword",
		Output:   outputPath,
		Force:    false,
	})
	if err == nil {
		t.Error("Show without --force on existing file should fail")
	}
}

func TestHideNameModeRandom(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: "testpassword",
		NameMode: "random",
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide with random name failed: %v", err)
	}

	if filepath.Base(hideResult.OutputPath) == "secret.txt.cfp" {
		t.Error("Output should have random name, not secret.txt.cfp")
	}

	if !strings.HasSuffix(hideResult.OutputPath, ".cfp") {
		t.Error("Output should still end with .cfp")
	}
}

func TestHideProfiles(t *testing.T) {
	tmpDir := t.TempDir()

	profiles := []string{"fast", "default", "paranoid"}

	for _, profile := range profiles {
		t.Run(profile, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "secret-"+profile+".txt")
			if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			hideResult, err := Hide(HideOptions{
				Path:     testFile,
				Password: "testpassword",
				Keep:     true,
				Profile:  profile,
				Yes:      true,
			})
			if err != nil {
				t.Fatalf("Hide with profile %s failed: %v", profile, err)
			}

			infoResult, err := Info(InfoOptions{
				Path: hideResult.OutputPath,
			})
			if err != nil {
				t.Fatalf("Info failed: %v", err)
			}

			if infoResult.Header.Profile != profile {
				t.Errorf("Profile should be %s, got %s", profile, infoResult.Header.Profile)
			}
		})
	}
}

func TestHideEmptyPassword(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := Hide(HideOptions{
		Path:     testFile,
		Password: "",
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err == nil {
		t.Error("Hide with empty password should fail")
	}
}

func TestShowEmptyPassword(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: "testpassword",
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	_, err = Show(ShowOptions{
		Path:     hideResult.OutputPath,
		Password: "",
		Output:   filepath.Join(tmpDir, "restored.txt"),
		Force:    true,
	})
	if err == nil {
		t.Error("Show with empty password should fail")
	}
}

func TestRekeyEmptyOldPassword(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: "testpassword",
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	_, err = Rekey(RekeyOptions{
		Path:        hideResult.OutputPath,
		OldPassword: "",
		NewPassword: "newpassword",
		Yes:         true,
	})
	if err == nil {
		t.Error("Rekey with empty old password should fail")
	}
}

func TestRekeyEmptyNewPassword(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: "testpassword",
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	_, err = Rekey(RekeyOptions{
		Path:        hideResult.OutputPath,
		OldPassword: "testpassword",
		NewPassword: "",
		GenerateKey: false,
		Yes:         true,
	})
	if err == nil {
		t.Error("Rekey with empty new password should fail")
	}
}

func TestHideNonExistentFile(t *testing.T) {
	_, err := Hide(HideOptions{
		Path:     "/nonexistent/file.txt",
		Password: "testpassword",
		Profile:  "fast",
		Yes:      true,
	})
	if err == nil {
		t.Error("Hide non-existent file should fail")
	}
}

func TestShowNonExistentFile(t *testing.T) {
	_, err := Show(ShowOptions{
		Path:     "/nonexistent/file.cfp",
		Password: "testpassword",
		Output:   "/tmp/restored.txt",
		Force:    true,
	})
	if err == nil {
		t.Error("Show non-existent file should fail")
	}
}

func TestVerifyNonExistentFile(t *testing.T) {
	err := Verify(VerifyOptions{
		Path:     "/nonexistent/file.cfp",
		Password: "testpassword",
	})
	if err == nil {
		t.Error("Verify non-existent file should fail")
	}
}

func TestInfoNonExistentFile(t *testing.T) {
	_, err := Info(InfoOptions{
		Path: "/nonexistent/file.cfp",
	})
	if err == nil {
		t.Error("Info non-existent file should fail")
	}
}

func TestRekeyNonExistentFile(t *testing.T) {
	_, err := Rekey(RekeyOptions{
		Path:        "/nonexistent/file.cfp",
		OldPassword: "testpassword",
		NewPassword: "newpassword",
		Yes:         true,
	})
	if err == nil {
		t.Error("Rekey non-existent file should fail")
	}
}

func TestRekeyWrongOldPassword(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: "testpassword",
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	_, err = Rekey(RekeyOptions{
		Path:        hideResult.OutputPath,
		OldPassword: "wrongpassword",
		NewPassword: "newpassword",
		Yes:         true,
	})
	if err == nil {
		t.Error("Rekey with wrong old password should fail")
	}
}

func TestHideShowLargeFile(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "large.txt")
	largeContent := make([]byte, 1024*1024)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	if err := os.WriteFile(testFile, largeContent, 0644); err != nil {
		t.Fatalf("failed to write large file: %v", err)
	}

	password := "testpassword123"

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: password,
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide large file failed: %v", err)
	}

	showResult, err := Show(ShowOptions{
		Path:     hideResult.OutputPath,
		Password: password,
		Output:   filepath.Join(tmpDir, "restored.txt"),
		Force:    true,
	})
	if err != nil {
		t.Fatalf("Show large file failed: %v", err)
	}

	restoredContent, err := os.ReadFile(showResult.OutputPath)
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}

	if len(restoredContent) != len(largeContent) {
		t.Errorf("restored file size doesn't match: got %d, want %d", len(restoredContent), len(largeContent))
	}

	for i := range largeContent {
		if restoredContent[i] != largeContent[i] {
			t.Errorf("content mismatch at byte %d", i)
			break
		}
	}
}

func TestHideShowEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "empty.txt")
	if err := os.WriteFile(testFile, []byte{}, 0644); err != nil {
		t.Fatalf("failed to write empty file: %v", err)
	}

	password := "testpassword123"

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: password,
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide empty file failed: %v", err)
	}

	showResult, err := Show(ShowOptions{
		Path:     hideResult.OutputPath,
		Password: password,
		Output:   filepath.Join(tmpDir, "restored.txt"),
		Force:    true,
	})
	if err != nil {
		t.Fatalf("Show empty file failed: %v", err)
	}

	restoredContent, err := os.ReadFile(showResult.OutputPath)
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}

	if len(restoredContent) != 0 {
		t.Errorf("restored file should be empty, got %d bytes", len(restoredContent))
	}
}

func TestHideShowSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "file with spaces & special-chars_123.txt")
	content := []byte("Content with special characters: ñ, á, é, 中文, 日本語")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	password := "testpassword123"

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: password,
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide file with special chars failed: %v", err)
	}

	showResult, err := Show(ShowOptions{
		Path:     hideResult.OutputPath,
		Password: password,
		Output:   filepath.Join(tmpDir, "restored.txt"),
		Force:    true,
	})
	if err != nil {
		t.Fatalf("Show file with special chars failed: %v", err)
	}

	restoredContent, err := os.ReadFile(showResult.OutputPath)
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}

	if string(restoredContent) != string(content) {
		t.Errorf("restored content doesn't match")
	}
}

func TestMultipleRekeys(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "secret.txt")
	content := []byte("Secret content")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	passwords := []string{"password1", "password2", "password3", "password4"}

	hideResult, err := Hide(HideOptions{
		Path:     testFile,
		Password: passwords[0],
		Keep:     true,
		Profile:  "fast",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Hide failed: %v", err)
	}

	for i := 1; i < len(passwords); i++ {
		_, err := Rekey(RekeyOptions{
			Path:        hideResult.OutputPath,
			OldPassword: passwords[i-1],
			NewPassword: passwords[i],
			Yes:         true,
		})
		if err != nil {
			t.Fatalf("Rekey %d failed: %v", i, err)
		}
	}

	showResult, err := Show(ShowOptions{
		Path:     hideResult.OutputPath,
		Password: passwords[len(passwords)-1],
		Output:   filepath.Join(tmpDir, "restored.txt"),
		Force:    true,
	})
	if err != nil {
		t.Fatalf("Show after multiple rekeys failed: %v", err)
	}

	restoredContent, err := os.ReadFile(showResult.OutputPath)
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}

	if string(restoredContent) != string(content) {
		t.Errorf("restored content doesn't match after multiple rekeys")
	}
}
