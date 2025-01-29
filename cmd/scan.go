package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/johnjallday/flow-workspace/internal/project"
)

// Usage prints how to use the scan tool.
func Usage() {
	fmt.Println("Usage: scan <root_directory>")
	fmt.Println("Example: scan /Users/jj/Workspace/johnj-programming")
}

func main() {
	// Check command-line arguments
	if len(os.Args) != 2 {
		Usage()
		os.Exit(1)
	}

	rootDir := os.Args[1]

	// Verify that the root directory exists
	info, err := os.Stat(rootDir)
	if os.IsNotExist(err) {
		log.Fatalf("Directory '%s' does not exist.", rootDir)
	}
	if !info.IsDir() {
		log.Fatalf("Path '%s' is not a directory.", rootDir)
	}

	var projects []project.Project

	// Walk through the directory recursively
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failed to access path %s: %v", path, err)
			return nil // Continue walking
		}

		// Check if the file is named 'project_info.toml'
		if !info.IsDir() && info.Name() == "project_info.toml" {
			fmt.Printf("Found project_info.toml: %s\n", path)
			p, err := project.LoadProjectInfo(path)
			if err != nil {
				log.Printf("Failed to load project info from '%s': %v", path, err)
				return nil // Continue walking
			}

			// If git_url exists in the project_info.toml, ensure it's captured
			// You might need to update the Project struct if git_url is not already present
			// Assuming it's present as per your example

			projects = append(projects, *p)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Error walking the path '%s': %v", rootDir, err)
	}

	if len(projects) == 0 {
		log.Println("No project_info.toml files found.")
		return
	}

	aggregatedProjects := project.Projects{
		Projects: projects,
	}

	// Define the path for 'projects.toml' at the root directory
	projectsTomlPath := filepath.Join(rootDir, "projects.toml")

	// Marshal the aggregated projects into TOML format
	output, err := toml.Marshal(aggregatedProjects)
	if err != nil {
		log.Fatalf("Failed to marshal projects to TOML: %v", err)
	}

	// Write the output to 'projects.toml'
	err = os.WriteFile(projectsTomlPath, output, 0644)
	if err != nil {
		log.Fatalf("Failed to write to '%s': %v", projectsTomlPath, err)
	}

	fmt.Printf("Successfully created '%s' with %d projects.\n", projectsTomlPath, len(projects))
}
