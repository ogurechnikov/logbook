package usecases

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ogurechnikov/logbook/internal/domain"
)

// InitVault создаёт новое хранилище заметок:
// 1. Создаёт директорию
// 2. Инициализирует Git-репозиторий
// 3. Создаёт .gitignore
// 4. Создаёт README.md
func InitVault(path string, gitRepo domain.GitRepository) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exist", path)
	}

	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := gitRepo.Init(path); err != nil {
		return fmt.Errorf("failed to init git %w", err)
	}

	gitignore := []byte(".DS_Store\n*.swp\n*.swo\n*~\n")
	if err := os.WriteFile(filepath.Join(path, ".gitignore"), gitignore, 0o644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	readme := []byte("# Logbook\n\nWelcome to your journal.\n\nWrite. Commit. Remember.\n")
	if err := os.WriteFile(filepath.Join(path, "README.md"), readme, 0o644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	return nil
}
