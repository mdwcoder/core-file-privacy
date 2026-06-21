package archive

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractArchive(data []byte, destDir string, force bool) error {
	return ExtractArchiveWithStrip(data, destDir, force, 0)
}

func ExtractArchiveWithStrip(data []byte, destDir string, force bool, stripComponents int) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	tr := tar.NewReader(bytes.NewReader(data))

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		name := header.Name

		if stripComponents > 0 {
			parts := strings.SplitN(name, "/", stripComponents+1)
			if len(parts) <= stripComponents {
				continue
			}
			name = parts[stripComponents]
			if name == "" {
				continue
			}
		}

		if err := validatePath(name); err != nil {
			return err
		}

		targetPath := filepath.Join(destDir, name)
		targetPath = filepath.Clean(targetPath)

		rel, err := filepath.Rel(destDir, targetPath)
		if err != nil {
			return fmt.Errorf("failed to compute relative path: %w", err)
		}
		if len(rel) >= 2 && rel[:2] == ".." {
			return fmt.Errorf("path traversal detected: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

		case tar.TypeReg:
			dir := filepath.Dir(targetPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			if !force {
				if _, err := os.Stat(targetPath); err == nil {
					return fmt.Errorf("file already exists: %s (use --force to overwrite)", targetPath)
				}
			}

			file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(file, tr); err != nil {
				file.Close()
				return fmt.Errorf("failed to write file: %w", err)
			}

			file.Close()
		}
	}

	return nil
}
