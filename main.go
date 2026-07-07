package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogurechnikov/logbook/internal/adapters"
	"github.com/ogurechnikov/logbook/internal/adapters/tui"
	"github.com/ogurechnikov/logbook/internal/usecases"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		handleInit(os.Args)
	}

	vaultPath := "."
	if len(os.Args) > 1 {
		vaultPath = os.Args[1]
	}
	runTUI(vaultPath)
}

func handleInit(args []string) {
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: logbook init <path>")
		os.Exit(1)
	}

	path := args[2]
	gitRepo := adapters.NewGoGitRepository()

	if err := usecases.InitVault(path, gitRepo); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("⚓ Journal created at %s\n", path)
	fmt.Println("Write. Commit. Remember.")
}

func runTUI(vaultPath string) {
	fileRepo := adapters.NewFileOS()
	model := tui.NewModel(vaultPath, fileRepo)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
