//go:build windows

package files

import (
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	setFileAttributesW = kernel32.NewProc("SetFileAttributesW")
	getFileAttributesW = kernel32.NewProc("GetFileAttributesW")
)

const (
	FILE_ATTRIBUTE_HIDDEN = 0x2
	INVALID_FILE_ATTRIBUTES = ^uintptr(0)
)

func HideFile(path string) error {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	ret, _, err := setFileAttributesW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		FILE_ATTRIBUTE_HIDDEN,
	)

	if ret == 0 {
		return err
	}

	return nil
}

func IsHidden(path string) bool {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}

	attrs, _, _ := getFileAttributesW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
	)

	if attrs == INVALID_FILE_ATTRIBUTES {
		return false
	}

	return attrs&FILE_ATTRIBUTE_HIDDEN != 0
}

func UnhideFile(path string) error {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	currentAttrs, _, _ := getFileAttributesW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
	)

	if currentAttrs == INVALID_FILE_ATTRIBUTES {
		return nil
	}

	newAttrs := currentAttrs &^ FILE_ATTRIBUTE_HIDDEN

	ret, _, err := setFileAttributesW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(newAttrs),
	)

	if ret == 0 {
		return err
	}

	return nil
}

func GetHiddenName(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	if !strings.HasPrefix(base, ".") {
		return filepath.Join(dir, "."+base)
	}

	return path
}
