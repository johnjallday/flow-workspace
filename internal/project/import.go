package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// ImportProject imports a project from the given directory.
// If project_info.toml does not exist, it creates one with default values,
// and automatically sets the project_type based on file extensions.
func ImportProject(projectDir string) error {
	metaFile := filepath.Join(projectDir, "project_info.toml")

	// Check if the file already exists.
	if _, err := os.Stat(metaFile); err == nil {
		fmt.Printf("Project info already exists at %s\n", metaFile)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking for project_info.toml: %v", err)
	}

	// Define file extensions to check.
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

	musicExtensions := map[string]bool{
		".rpp": true, // REAPER project file
		".als": true, // Ableton Live Set
		".flp": true, // FL Studio project file
		// Add additional music file extensions as needed.
	}

	foundCoding := false
	foundMusic := false

	// Walk the directory (recursively) to determine project type.
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip file errors.
			return nil
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if musicExtensions[ext] {
			foundMusic = true
			// Stop early if a music file is found.
			return filepath.SkipDir
		}
		if codingExtensions[ext] {
			foundCoding = true
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking project directory: %v", err)
	}

	// Determine project type.
	projectType := "general"
	if foundMusic {
		projectType = "music"
	} else if foundCoding {
		projectType = "coding"
	}

	// Create a new Project with default values.
	proj := Project{
		Name:         filepath.Base(projectDir),
		Alias:        filepath.Base(projectDir),
		ProjectType:  projectType,
		Tags:         []string{},
		Notes:        []string{},
		Path:         projectDir,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}

	// Create the project_info.toml file.
	f, err := os.Create(metaFile)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", metaFile, err)
	}
	defer f.Close()

	// Encode the project data into TOML format.
	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(proj); err != nil {
		return fmt.Errorf("failed to encode project info: %v", err)
	}

	fmt.Printf("Created new project_info.toml at %s with project type '%s'\n", metaFile, proj.ProjectType)
	return nil
}
