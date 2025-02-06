package repl

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/johnjallday/flow-workspace/internal/project"
	"github.com/johnjallday/flow-workspace/internal/root"
	"github.com/johnjallday/flow-workspace/internal/workspace"
)

func StartREPL() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Scope detection: check for known markers in cwd.
	if info, err := os.Stat(filepath.Join(cwd, ".config")); err == nil && info.IsDir() {
		fmt.Println("Detected .config folder. Launching Root REPL.")
		root.StartRootREPL(cwd)
		return
	}

	if _, err := os.Stat(filepath.Join(cwd, "ws_info.toml")); err == nil {
		fmt.Println("Detected ws_info.toml. Launching Workspace REPL.")
		workspace.StartWorkspaceREPL(cwd)
		return
	}

	if _, err := os.Stat(filepath.Join(cwd, "project_info.toml")); err == nil {
		fmt.Println("Detected project_info.toml. Launching Project REPL.")
		project.StartProjectREPL(cwd)
		return
	}

	// Otherwise, fall back to a default message.
	fmt.Println("Unrecognized scope for TODO REPL. Exiting.")
}
