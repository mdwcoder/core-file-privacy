package format

import (
	"encoding/json"
	"fmt"
)

type Metadata struct {
	OriginalName string `json:"original_name"`
	OriginalType string `json:"original_type"`
	OriginalMode string `json:"original_mode,omitempty"`
	OriginalSize int64  `json:"original_size,omitempty"`
	Archived     bool   `json:"archived"`
}

func NewFileMetadata(name string, mode string, size int64) *Metadata {
	return &Metadata{
		OriginalName: name,
		OriginalType: "file",
		OriginalMode: mode,
		OriginalSize: size,
		Archived:     false,
	}
}

func NewDirectoryMetadata(name string) *Metadata {
	return &Metadata{
		OriginalName: name,
		OriginalType: "directory",
		Archived:     true,
	}
}

func (m *Metadata) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func MetadataFromJSON(data []byte) (*Metadata, error) {
	var metadata Metadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	return &metadata, nil
}

func (m *Metadata) IsDirectory() bool {
	return m.OriginalType == "directory" || m.Archived
}
