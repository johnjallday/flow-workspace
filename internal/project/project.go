package project

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

// Project represents a single project's metadata from project_info.toml.
type Project struct {
	Name         string        `toml:"name"`
	Alias        string        `toml:"alias"`
	ProjectType  string        `toml:"project_type"`
	Tags         []string      `toml:"tags"`
	DateCreated  time.Time     `toml:"date_created"`
	DateModified time.Time     `toml:"date_modified"`
	Notes        []string      `toml:"notes"`
	Path         string        `toml:"path"`
	GitURL       string        `toml:"git_url,omitempty"`
	MusicDetails *MusicDetails `toml:"music_details,omitempty"`
}

// LoadProjectInfo reads and parses a project_info.toml file into a Project.
func LoadProjectInfo(filename string) (*Project, error) {
	var proj Project

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("file '%s' does not exist", filename)
	}

	if _, err := toml.DecodeFile(filename, &proj); err != nil {
		return nil, fmt.Errorf("error decoding TOML file: %w", err)
	}

	if proj.Name == "" {
		return nil, fmt.Errorf("'name' field cannot be empty in '%s'", filename)
	}

	// Set defaults if not provided
	if proj.Alias == "" {
		proj.Alias = proj.Name
	}
	if proj.ProjectType == "" {
		proj.ProjectType = "General"
	}
	if proj.Path == "" {
		proj.Path = "./"
	}

	// Find latest file modification date in the project directory
	latestFileTime, err := GetLatestFileModTime(proj.Path)
	if err != nil {
		return nil, fmt.Errorf("error getting latest file mod time: %v", err)
	}

	// If it's newer than what's recorded, update DateModified
	if latestFileTime.After(proj.DateModified) {
		fmt.Printf("Updating DateModified for project '%s' to %v\n", proj.Name, latestFileTime)
		proj.DateModified = latestFileTime

		// Save the updated project info back to the file
		if err := saveProjectInfo(filename, &proj); err != nil {
			return nil, fmt.Errorf("failed to save updated project info: %w", err)
		}
	}

	return &proj, nil
}

// GetLatestFileModTime recursively walks through the project directory
// and returns the latest modification time among all files.
func GetLatestFileModTime(projectDir string) (time.Time, error) {
	var latestModTime time.Time

	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip files that can't be accessed.
			return nil
		}
		if !info.IsDir() {
			if info.ModTime().After(latestModTime) {
				latestModTime = info.ModTime()
			}
		}
		return nil
	})
	if err != nil {
		return time.Time{}, fmt.Errorf("error walking project directory: %v", err)
	}

	return latestModTime, nil
}

func saveProjectInfo(filename string, proj *Project) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(proj); err != nil {
		return fmt.Errorf("failed to encode project info: %v", err)
	}

	return nil
}
