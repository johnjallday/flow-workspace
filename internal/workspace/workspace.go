package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/johnjallday/flow-workspace/internal/project"
)

// Workspace represents a directory containing multiple projects.toml entries.
type Workspace struct {
	Path string
	// ... other metadata if needed
}

// Projects represents a collection of Projects aggregated into projects.toml.
type Projects struct {
	Projects []project.Project `toml:"projects"`
}

// LoadProjectsToml loads an existing projects.toml from a workspace.
func LoadProjectsToml(workspacePath string) (*Projects, error) {
	projectsTomlPath := filepath.Join(workspacePath, "projects.toml")
	if _, err := os.Stat(projectsTomlPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file '%s' does not exist", projectsTomlPath)
	}

	var projs Projects
	if _, err := toml.DecodeFile(projectsTomlPath, &projs); err != nil {
		return nil, fmt.Errorf("error decoding TOML file: %w", err)
	}

	return &projs, nil
}

// SaveProjectsToml writes a Projects struct to projects.toml in the workspace.
func SaveProjectsToml(projs *Projects, workspacePath string) error {
	projectsTomlPath := filepath.Join(workspacePath, "projects.toml")
	output, err := toml.Marshal(projs)
	if err != nil {
		return fmt.Errorf("failed to marshal projects to TOML: %w", err)
	}

	if err := os.WriteFile(projectsTomlPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write '%s': %w", projectsTomlPath, err)
	}

	return nil
}

// Skip directories and files matching these names
var skipDirs = map[string]bool{
	".TagStudio": true,
	".config":    true,
}

var skipFiles = map[string]bool{
	".DS_Store": true,
}

// UpdateProjects scans the workspace directory for project_info.toml files,
// aggregates them into a Projects struct, and writes the data to projects.toml.
// It returns the aggregated Projects pointer or an error.
func UpdateProjects(workspaceDir string) (*Projects, error) {
	var allProjects []project.Project

	// Directories and files to skip.
	skipDirs := map[string]bool{
		".TagStudio": true,
		".config":    true,
	}
	skipFiles := map[string]bool{
		".DS_Store": true,
	}

	// Walk the workspace directory to find all project_info.toml files.
	err := filepath.Walk(workspaceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Failed to access path %s: %v\n", path, err)
			return nil // Continue walking.
		}

		// Skip directories we don't want to descend into.
		if info.IsDir() && skipDirs[info.Name()] {
			return filepath.SkipDir
		}

		// Skip unwanted files.
		if !info.IsDir() && skipFiles[info.Name()] {
			return nil
		}

		// Look for project_info.toml files.
		if !info.IsDir() && info.Name() == "project_info.toml" {
			fmt.Printf("Found project_info.toml: %s\n", path)
			p, loadErr := project.LoadProjectInfo(path)
			if loadErr != nil {
				fmt.Printf("Failed to load project info from '%s': %v\n", path, loadErr)
				return nil // Continue walking even if one fails.
			}
			allProjects = append(allProjects, *p)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking path '%s': %w", workspaceDir, err)
	}

	if len(allProjects) == 0 {
		return nil, fmt.Errorf("no project_info.toml files found in workspace %s", workspaceDir)
	}

	// Aggregate the discovered projects.
	aggregated := Projects{
		Projects: allProjects,
	}

	// Define the path for projects.toml in the workspace.
	projectsTomlPath := filepath.Join(workspaceDir, "projects.toml")

	// Marshal the aggregated projects into TOML format.
	output, err := toml.Marshal(aggregated)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal projects to TOML: %w", err)
	}

	// Write the TOML data to the projects.toml file.
	err = os.WriteFile(projectsTomlPath, output, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write to '%s': %w", projectsTomlPath, err)
	}

	return &aggregated, nil
}

// ScanAndAggregateProjects scans the workspace directory for project_info.toml files
// and aggregates them into a Projects struct. It ignores the directories .TagStudio, .config,
// and any file named .DS_Store.
func ScanAndAggregateProjects(rootDir string) (*Projects, error) {
	var allProjects []project.Project

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Failed to access path %s: %v\n", path, err)
			return nil // Continue walking
		}

		// Skip directories we don't want to descend into
		if info.IsDir() && skipDirs[info.Name()] {
			return filepath.SkipDir
		}

		// Skip files we don't want to process
		if !info.IsDir() && skipFiles[info.Name()] {
			return nil
		}

		// Check if the file is named 'project_info.toml'
		if !info.IsDir() && info.Name() == "project_info.toml" {
			fmt.Printf("Found project_info.toml: %s\n", path)
			p, loadErr := project.LoadProjectInfo(path)
			if loadErr != nil {
				fmt.Printf("Failed to load project info: %v\n", loadErr)
				return nil // Continue walking
			}
			allProjects = append(allProjects, *p)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking path '%s': %w", rootDir, err)
	}

	if len(allProjects) == 0 {
		return nil, fmt.Errorf("no project_info.toml files found")
	}

	aggregated := Projects{
		Projects: allProjects,
	}

	return &aggregated, nil
}

// ListProjects prints all projects in a Projects struct.
func ListProjects(projs *Projects) {
	if projs == nil || len(projs.Projects) == 0 {
		fmt.Println("No projects found.")
		return
	}
	fmt.Println("\nProjects in this Workspace:")
	for i, proj := range projs.Projects {
		fmt.Printf("--------------------------------------------------\n")
		fmt.Printf("Project #%d\n", i+1)
		fmt.Printf("Name         : %s\n", proj.Name)
		fmt.Printf("Alias        : %s\n", proj.Alias)
		fmt.Printf("Type         : %s\n", proj.ProjectType)
		fmt.Printf("Tags         : %v\n", proj.Tags)
		fmt.Printf("Date Created : %s\n", proj.DateCreated.Format("2006-01-02"))
		fmt.Printf("Date Modified: %s\n", proj.DateModified.Format("2006-01-02"))
		if proj.GitURL != "" {
			fmt.Printf("Git URL      : %s\n", proj.GitURL)
		}
		fmt.Printf("Path         : %s\n", proj.Path)
		fmt.Printf("--------------------------------------------------\n\n")
	}
}
