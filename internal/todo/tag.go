package todo

import (
	"os"
	"path/filepath"
	"time"
)

// TagProject returns a tag for the project if a "project_info.toml" exists in projectPath.
func TagProject(projectPath string) string {
	projectInfoPath := filepath.Join(projectPath, "project_info.toml")
	if _, err := os.Stat(projectInfoPath); err == nil {
		return filepath.Base(projectPath)
	}
	return ""
}

// TagWorkspace returns a workspace tag based on finding a "ws_info.toml"
// two directories above the given projectPath.
func TagWorkspace(projectPath string) string {
	parentPath := filepath.Dir(projectPath)
	workspacePath := filepath.Dir(parentPath)
	wsInfoPath := filepath.Join(workspacePath, "ws_info.toml")
	if _, err := os.Stat(wsInfoPath); err == nil {
		return filepath.Base(workspacePath)
	}
	return ""
}

// BuildTaskLine constructs a complete task line (including tags) to be added to a TODO file.
func BuildTaskLine(description, dueDate, projectName, workspaceName string) string {
	tags := ""
	if dueDate != "" {
		tags += " #due:" + dueDate
	}
	if projectName != "" {
		tags += " #project:" + projectName
	}
	if workspaceName != "" {
		tags += " #workspace:" + workspaceName
	}
	tags += " #created:" + time.Now().Format("2006-01-02")
	return "- [ ] " + description + tags
}
