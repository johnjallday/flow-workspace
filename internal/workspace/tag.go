package workspace

import (
	"os"

	"github.com/BurntSushi/toml"
)

// WorkspaceInfo represents the TOML structure
type WorkspaceInfo struct {
	Accounts []string `toml:"accounts"`
	Aliases  []string `toml:"aliases"`
	Tags     []string `toml:"tags"`
	Projects []string `toml:"projects"`
}

// LoadTOML reads a TOML file and decodes it into WorkspaceInfo
func LoadTOML(filename string) (*WorkspaceInfo, error) {
	var workspace WorkspaceInfo
	if _, err := toml.DecodeFile(filename, &workspace); err != nil {
		return nil, err
	}
	return &workspace, nil
}

// SaveTOML saves the WorkspaceInfo struct into a TOML file
func SaveTOML(filename string, workspace *WorkspaceInfo) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(workspace)
}
