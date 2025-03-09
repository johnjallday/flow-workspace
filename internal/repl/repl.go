package repl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/johnjallday/flow-workspace/internal/project"
	"github.com/johnjallday/flow-workspace/internal/root"
	"github.com/johnjallday/flow-workspace/internal/workspace"
)

// StartREPL detects the appropriate REPL to launch and allows clearing the screen with Ctrl+L.
func StartREPL(dbPath string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current directory: %v", err)
		}

		// Detect REPL scope and launch the appropriate one.
		if info, err := os.Stat(filepath.Join(cwd, ".config")); err == nil && info.IsDir() {
			fmt.Println("Detected .config folder. Launching Root REPL.")
			root.StartRootREPL(dbPath, cwd)
			return
		}

		if _, err := os.Stat(filepath.Join(cwd, "ws_info.toml")); err == nil {
			fmt.Println("Detected ws_info.toml. Launching Workspace REPL.")
			workspace.StartWorkspaceREPL(dbPath, cwd)
			return
		}

		if _, err := os.Stat(filepath.Join(cwd, "project_info.toml")); err == nil {
			fmt.Println("Detected project_info.toml. Launching Project REPL.")
			project.StartProjectREPL(dbPath, cwd)
			return
		}

		// No known scope file found.
		// Scan for potential coding projects in the current directory, its parent, and immediate subdirectories.
		candidates := []string{}
		// Add current directory.
		candidates = append(candidates, cwd)

		// Add parent directory if it is different.
		parent := filepath.Dir(cwd)
		if parent != cwd {
			candidates = append(candidates, parent)
		}

		// Add immediate subdirectories.
		entries, err := os.ReadDir(cwd)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					candidates = append(candidates, filepath.Join(cwd, entry.Name()))
				}
			}
		}

		// Filter candidates that appear to be coding projects.
		var validCandidates []string
		for _, cand := range candidates {
			if isCodingProject(cand) {
				validCandidates = append(validCandidates, cand)
			}
		}

		if len(validCandidates) > 0 {
			fmt.Println("The following directories appear to be coding projects:")
			for i, cand := range validCandidates {
				fmt.Printf("  %d) %s\n", i+1, cand)
			}
			fmt.Print("This looks like a coding project. Would you like to import one of these directories? (Enter number, or press Enter to retry): ")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line != "" {
				index, err := strconv.Atoi(line)
				if err == nil && index >= 1 && index <= len(validCandidates) {
					selected := validCandidates[index-1]
					if err := project.ImportProject(selected); err != nil {
						fmt.Printf("Error importing project: %v\n", err)
					} else {
						fmt.Printf("Project imported successfully. Launching Project REPL for %s\n", selected)
						project.StartProjectREPL(dbPath, selected)
						return
					}
				}
			}
		} else {
			fmt.Println("Unrecognized scope for TODO REPL. Press 'Enter' to retry or 'Ctrl + L' to clear.")
			fmt.Print("\n[repl] >> ")
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
			line = strings.TrimSpace(line)
			// Handle Ctrl+L (ASCII 12).
			if line == "\x0c" {
				clearScreen()
				continue
			}
			if line != "" {
				fmt.Println("Exiting REPL.")
				return
			}
		}
	}
}

// isCodingProject checks if the given directory contains a go.mod file or at least one file
// with a common coding file extension.
func isCodingProject(dir string) bool {
	// Check if go.mod exists in the directory.
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		return true
	}

	codingExtensions := map[string]bool{
		".go":    true,
		".py":    true,
		".js":    true,
		".java":  true,
		".c":     true,
		".cpp":   true,
		".cs":    true,
		".rb":    true,
		".php":   true,
		".ts":    true,
		".swift": true,
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if codingExtensions[ext] {
				return true
			}
		}
	}
	return false
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("clear") // Linux
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls") // Windows
	case "darwin":
		cmd = exec.Command("clear") // macOS
	default:
		fmt.Println("CLS for", runtime.GOOS, "not implemented")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
