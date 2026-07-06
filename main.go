package main

import (
	"fmt"
	"os"

	"github.com/ogurechnikov/logbook/internal/adapters"
	"github.com/ogurechnikov/logbook/internal/usecases"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		handleInit(os.Args)
	}
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
