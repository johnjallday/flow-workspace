package root

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	// Import the workspace REPL so we can call StartWorkspaceREPL
	"github.com/johnjallday/flow-workspace/internal/workspace"
)

// StartRootREPL starts an interactive REPL at the root level.

func StartRootREPL(dbPath string, rootDir string) {
	fmt.Println("Welcome to the ROOT-level REPL!")
	fmt.Printf("Root Directory: %s\n", rootDir)
	printRootHelp()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n[ROOT] >> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		line = strings.TrimSpace(line)

		switch strings.ToLower(line) {
		case "exit":
			fmt.Println("Exiting Root REPL. Goodbye!")
			return

		case "help":
			printRootHelp()

		case "list":
			// Lists all workspaces at the root level.
			ListWorkspaces(rootDir)

		case "projects":
			// List subdirectories that contain a projects.toml.
			ListProjects(rootDir)

		case "select":
			// Let the user select a workspace.
			ListWorkspaces(rootDir)
			selectedWorkspace := selectWorkspace(rootDir, reader)
			if selectedWorkspace != "" {
				workspace.StartWorkspaceREPL(dbPath, selectedWorkspace)
			}

		case "todo":
			// New command: aggregate and print all todos from all workspaces.
			ListAllTodos(rootDir)

		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

// selectWorkspace lets the user pick from the listed workspace directories
// and returns the selected directory's *absolute path* (or an empty string if canceled/invalid).
func selectWorkspace(rootDir string, reader *bufio.Reader) string {
	// Directories to skip
	skip := map[string]bool{
		".config":     true,
		".spacedrive": true,
		".DS_Store":   true,
		".TagStudio":  true,
	}

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		log.Printf("Failed to read root dir: %v\n", err)
		return ""
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() && !skip[e.Name()] {
			dirs = append(dirs, e.Name())
		}
	}
	if len(dirs) == 0 {
		fmt.Println("No workspaces found (or all were skipped).")
		return ""
	}

	for {
		fmt.Print("Enter the number of the workspace to switch to (or 'cancel'): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return ""
		}
		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "cancel") {
			return ""
		}

		idx, convErr := strconv.Atoi(input)
		if convErr != nil || idx < 1 || idx > len(dirs) {
			fmt.Println("Invalid selection. Try again.")
			continue
		}

		selected := dirs[idx-1]
		selectedPath := filepath.Join(rootDir, selected)
		fmt.Printf("Workspace selected: %s\n", selectedPath)

		return selectedPath
	}
}

// printRootHelp displays the available commands in the root-level REPL.

func printRootHelp() {
	fmt.Println(`Available commands (root-level):
  help      - Show this help message
  list      - List all workspaces in the root directory
  projects  - List subdirectories that contain 'projects.toml'
  select    - Select a workspace by number (and load workspace REPL)
  todo      - Aggregate and list all TODOs from every workspace
  exit      - Exit this Root REPL`)
}
