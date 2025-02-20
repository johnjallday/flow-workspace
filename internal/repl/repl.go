package repl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/johnjallday/flow-workspace/internal/project"
	"github.com/johnjallday/flow-workspace/internal/root"
	"github.com/johnjallday/flow-workspace/internal/workspace"
)

// StartREPL detects the appropriate REPL to launch and allows clearing the screen with Ctrl+L.
func StartREPL() {
	reader := bufio.NewReader(os.Stdin)

	for {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current directory: %v", err)
		}

		// Detect REPL scope and launch the appropriate one.
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

		// If no valid REPL is found, prompt the user.
		fmt.Println("Unrecognized scope for TODO REPL. Press 'Enter' to retry or 'Ctrl + L' to clear.")

		// Read user input
		fmt.Print("\n[repl] >> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		line = strings.TrimSpace(line)

		// Handle Ctrl+L
		if line == "\x0c" { // ASCII 12 = Ctrl + L
			clearScreen()
			continue // Restart the loop to re-detect the REPL
		}

		// Exit if the user provides any other input
		if line != "" {
			fmt.Println("Exiting REPL.")
			return
		}
	}
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
