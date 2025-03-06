package root

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	// Import the workspace package to load and list projects
	"github.com/johnjallday/flow-workspace/internal/todo"
	"github.com/johnjallday/flow-workspace/internal/workspace"
)

// ListWorkspaces reads the root directory and prints out
// all sub-directories except for certain hidden/system ones.
func ListWorkspaces(rootDir string) {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		log.Printf("Failed to read root dir '%s': %v\n", rootDir, err)
		return
	}

	// Directories to skip
	skip := map[string]bool{
		".config":     true,
		".spacedrive": true,
		".DS_Store":   true,
		".TagStudio":  true,
	}

	fmt.Printf("\nAvailable Workspaces in %s:\n", rootDir)
	count := 0

	for _, e := range entries {
		name := e.Name()
		if skip[name] {
			continue
		}
		if e.IsDir() {
			count++
			fmt.Printf("%d) %s\n", count, name)
		}
	}

	if count == 0 {
		fmt.Println("No workspaces found (or all were skipped).")
	}
}

// ListProjects looks for a `projects.toml` file in each subdirectory of rootDir
// (one level only). If found, it loads and prints the contained projects.
func ListProjects(rootDir string) {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		log.Printf("Failed to read root dir '%s': %v\n", rootDir, err)
		return
	}

	skip := map[string]bool{
		".config":     true,
		".spacedrive": true,
		".DS_Store":   true,
		".TagStudio":  true,
	}

	fmt.Printf("\nScanning for 'projects.toml' in subfolders of: %s\n", rootDir)
	foundAny := false

	for _, e := range entries {
		name := e.Name()
		if skip[name] {
			continue
		}
		if e.IsDir() {
			// Check if this directory has 'projects.toml'
			candidate := filepath.Join(rootDir, name, "projects.toml")
			if _, err := os.Stat(candidate); err == nil {
				foundAny = true
				fmt.Printf("\nFound 'projects.toml' in: %s\n", name)

				// Load and parse the projects.toml
				projs, loadErr := workspace.LoadProjectsToml(filepath.Join(rootDir, name))
				if loadErr != nil {
					fmt.Printf("  -> Error loading 'projects.toml': %v\n", loadErr)
					continue
				}

				// Print the loaded projects from this subdirectory
				workspace.ListProjects(projs)
			}
		}
	}

	if !foundAny {
		fmt.Println("\nNo 'projects.toml' found in any subdirectory.")
	}
}

// ListAllTodos aggregates and prints all TODOs from every workspace found under rootDir.
func ListAllTodos(rootDir string) {
	// Directories to skip at the root level.
	skip := map[string]bool{
		".config":     true,
		".spacedrive": true,
		".DS_Store":   true,
		".TagStudio":  true,
	}

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		log.Printf("Failed to read root directory '%s': %v\n", rootDir, err)
		return
	}

	var aggregatedTodos []todo.Todo

	// Loop through each subdirectory (workspace) in the root directory.
	for _, e := range entries {
		if !e.IsDir() || skip[e.Name()] {
			continue
		}

		workspacePath := filepath.Join(rootDir, e.Name())
		projectsTomlPath := filepath.Join(workspacePath, "projects.toml")
		if _, err := os.Stat(projectsTomlPath); os.IsNotExist(err) {
			// Skip this workspace if no projects.toml exists.
			fmt.Printf("Skipping workspace '%s': no projects.toml found.\n", workspacePath)
			continue
		}

		// Load the projects from the workspace’s projects.toml.
		projs, err := workspace.LoadProjectsToml(workspacePath)
		if err != nil {
			fmt.Printf("Skipping workspace '%s': error loading projects.toml: %v\n", workspacePath, err)
			continue
		}

		// For each project, decide which folder to check for a todo.md.
		for _, proj := range projs.Projects {
			var projectDir string

			// Use the project Path if it’s provided (and not just the default "./")
			if proj.Path != "" && proj.Path != "./" {
				if filepath.IsAbs(proj.Path) {
					projectDir = filepath.Clean(proj.Path)
				} else {
					projectDir = filepath.Join(workspacePath, proj.Path)
				}
			} else {
				// Default: assume the project is in a folder named after the project.
				projectDir = filepath.Join(workspacePath, proj.Name)
			}

			todoFile := filepath.Join(projectDir, "todo.md")
			if _, err := os.Stat(todoFile); os.IsNotExist(err) {
				// Skip if todo.md does not exist.
				continue
			}

			// Load the tasks from the todo.md file.
			tasks, err := todo.LoadAllTodos(todoFile)
			if err != nil {
				fmt.Printf("Error loading todos from '%s': %v\n", todoFile, err)
				continue
			}

			// Optionally, annotate each task with the project and workspace names.
			for i := range tasks {
				if tasks[i].ProjectName == "" {
					tasks[i].ProjectName = proj.Name
				}
				if tasks[i].WorkspaceName == "" {
					tasks[i].WorkspaceName = e.Name() // using the workspace folder name
				}
			}

			aggregatedTodos = append(aggregatedTodos, tasks...)
		}
	}

	// Print the aggregated list.
	if len(aggregatedTodos) == 0 {
		fmt.Println("No TODOs found in any workspace.")
		return
	}

	fmt.Println("\nAggregated TODOs across all Workspaces:")

	todo.PrintTodos(aggregatedTodos)
}
