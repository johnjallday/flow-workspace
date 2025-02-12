package project

import (
	"fmt"
	"os"
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

	return &proj, nil
}
