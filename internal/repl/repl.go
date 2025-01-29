package repl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/johnjallday/flow-workspace/internal/project"
)

// StartREPL initializes the REPL loop with the given root workspace directory.
func StartREPL(rootWorkspaceDir string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the TODO Manager!")
	fmt.Println("Type 'help' for available commands or 'exit' to quit.")

	var selectedWorkspace string
	var projectsInfo *project.Projects

	// Initial Workspace Selection
	selectedWorkspace, projectsInfo = selectWorkspace(rootWorkspaceDir, reader)
	if selectedWorkspace == "" {
		fmt.Println("No workspace selected. Exiting.")
		return
	}

	// Use projectsInfo to list projects
	project.ListProjects(projectsInfo)

	for {
		fmt.Printf("\n[%s] >> ", filepath.Base(selectedWorkspace))
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		line = strings.TrimSpace(line)

		// Detect Ctrl+L (ASCII Form Feed, 0x0C)
		if line == "\f" {
			clearScreen()
			continue
		}

		switch strings.ToLower(line) {
		case "exit":
			fmt.Println("Exiting TODO Manager. Goodbye!")
			return
		case "help":
			printHelp()
		case "workspace":
			// Allow workspace switching
			selectedWorkspace, projectsInfo = selectWorkspace(rootWorkspaceDir, reader)
			if selectedWorkspace == "" {
				fmt.Println("No workspace selected.")
			} else {
				// Use projectsInfo to list projects after switching
				project.ListProjects(projectsInfo)
			}
		case "todo":
			handleTodoSubCommand(reader, selectedWorkspace)
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

// selectWorkspace prompts the user to select a workspace from the root directory.
func selectWorkspace(rootWorkspaceDir string, reader *bufio.Reader) (string, *project.Projects) {
	workspaces, err := os.ReadDir(rootWorkspaceDir)
	if err != nil {
		log.Fatalf("Failed to read workspaces directory: %v", err)
	}

	// Filter only directories
	var workspaceDirs []os.DirEntry
	for _, ws := range workspaces {
		if ws.IsDir() {
			workspaceDirs = append(workspaceDirs, ws)
		}
	}

	if len(workspaceDirs) == 0 {
		fmt.Println("No workspaces found in the root directory.")
		return "", nil
	}

	fmt.Println("\nAvailable Workspaces:")
	for i, ws := range workspaceDirs {
		fmt.Printf("%d. %s\n", i+1, ws.Name())
	}

	fmt.Print("Select a workspace by number (or type 'cancel' to abort): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if strings.ToLower(input) == "cancel" {
		return "", nil
	}

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(workspaceDirs) {
		fmt.Println("Invalid selection.")
		return "", nil
	}

	selectedWS := workspaceDirs[index-1].Name()
	selectedWSPath := filepath.Join(rootWorkspaceDir, selectedWS)
	projectsTomlPath := filepath.Join(selectedWSPath, "projects.toml")

	// Load projects.toml
	projectsInfo, err := project.LoadProjectsInfo(projectsTomlPath)
	if err != nil {
		fmt.Printf("Failed to load 'projects.toml' in workspace '%s': %v\n", selectedWS, err)
		// Optionally, initialize an empty Projects struct
		projectsInfo = &project.Projects{
			Projects: []project.Project{},
		}
	}

	fmt.Printf("Selected Workspace: %s\n", selectedWS)

	return selectedWSPath, projectsInfo
}

// printHelp displays the available commands.
func printHelp() {
	fmt.Println(`Available commands:
  help        - Print this help message
  workspace   - Switch to a different workspace
  todo        - Manage your TODO list (View, Input, Complete, Cleanup, Delete, Filter, Edit)
  exit        - Exit the application
`)
}

// handleTodoSubCommand handles the 'todo' command by offering options to View, Input, Complete, Cleanup, Delete, Filter, Edit.
func handleTodoSubCommand(reader *bufio.Reader, workspacePath string) {
	project.HandleTodoView(workspacePath)
}

// clearScreen clears the terminal screen.
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
