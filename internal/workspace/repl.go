// internal/workspace/repl.go
package workspace

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/johnjallday/flow-workspace/internal/project"
)

// StartWorkspaceREPL starts an interactive REPL for the specified workspace directory.
func StartWorkspaceREPL(dbPath string, workspaceDir string) {
	reader := bufio.NewReader(os.Stdin)
	// Extract the base folder name from workspaceDir.
	currentWorkspace := filepath.Base(workspaceDir)
	fmt.Printf("Workspace REPL started for directory: %s\n", workspaceDir)
	fmt.Printf("Current Workspace: %s\n", currentWorkspace)
	printWorkspaceHelp()

	// Attempt to load the workspace's projects.toml
	projs, err := LoadProjectsToml(workspaceDir)
	if err != nil {
		fmt.Printf("Error loading '%s/projects.toml': %v\n", workspaceDir, err)
		// Optionally, you can create an empty Projects if needed:
		projs = &Projects{}
	}

	for {
		fmt.Printf("\n[workspace:%s] >> ", currentWorkspace)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		line = strings.TrimSpace(line)

		switch strings.ToLower(line) {
		case "exit":
			fmt.Println("Exiting Workspace REPL. Goodbye!")
			return

		case "help":
			printWorkspaceHelp()

		case "list projects":
			// Display the projects using your existing ListProjects() function
			ListProjects(projs)

		case "todo":
			// Lists aggregated TODOs for all projects in this workspace.
			// Note: we now call the function from the new aggregated workspace todo package.
			ListAllTodos(workspaceDir)

		case "select project":
			// Let the user choose a project by number, then start the Project REPL
			selectProject(dbPath, workspaceDir, projs, reader)

		case "update projects":
			// Update projects by scanning the workspace.
			updatedProjs, err := UpdateProjects(workspaceDir)
			if err != nil {
				fmt.Printf("Error updating projects: %v\n", err)
			} else {
				projs = updatedProjs
				fmt.Println("Projects updated successfully!")
				// Optionally, list the projects.
				ListProjects(projs)
			}

		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

// printWorkspaceHelp shows commands specific to the workspace-level REPL.
func printWorkspaceHelp() {
	fmt.Println(`Available commands (Workspace REPL):
  help             - Show this help message
  list projects    - List all projects in this workspace
  todo             - List aggregated TODOs from all projects in this workspace
  select project   - Choose a project to open the Project REPL
  exit             - Exit the Workspace REPL
  update projects  - Scan the workspace for new projects and update the projects.toml file`)
}

// selectProject lets the user pick from the loaded Projects, then starts the Project REPL.
func selectProject(dbPath string, workspaceDir string, projs *Projects, reader *bufio.Reader) {
	if projs == nil || len(projs.Projects) == 0 {
		fmt.Println("No projects found in this workspace.")
		return
	}

	fmt.Println("\nSelect a project to open its REPL:")
	for i, p := range projs.Projects {
		fmt.Printf("%d) %s (Path: %s)\n", i+1, p.Name, p.Path)
	}

	for {
		fmt.Print("Enter project number (or 'cancel' to abort): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "cancel") {
			return
		}

		idx, convErr := strconv.Atoi(input)
		if convErr != nil || idx < 1 || idx > len(projs.Projects) {
			fmt.Println("Invalid project number, try again.")
			continue
		}

		chosenProject := projs.Projects[idx-1]
		p := chosenProject.Path

		var projectDir string
		if filepath.IsAbs(p) {
			// It's already an absolute path
			projectDir = filepath.Clean(p)
		} else {
			// It's relative, so combine it
			projectDir = filepath.Join(workspaceDir, p)
		}

		fmt.Printf("Selected Project: %s\n", chosenProject.Name)
		// Start the project-level REPL.
		project.StartProjectREPL(dbPath, projectDir)
		return
	}
}
