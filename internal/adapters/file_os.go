package adapters

import (
	"os"
	"path/filepath"
	"strings"
)

// FileOS — реализация FileRepository через стандартную файловую систему.
type FileOS struct{}

func NewFileOS() *FileOS {
	return &FileOS{}
}

func (f *FileOS) List(vaultPath string) ([]string, error) {
	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".md") {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func (f *FileOS) Read(vaultPath, filename string) (string, error) {
	content, err := os.ReadFile(filepath.Join(vaultPath, filename))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
