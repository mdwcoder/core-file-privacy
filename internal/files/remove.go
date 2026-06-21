package files

import (
	"fmt"
	"os"
)

func SecureRemove(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}
	return nil
}

func SecureRemoveAll(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove path: %w", err)
	}
	return nil
}
