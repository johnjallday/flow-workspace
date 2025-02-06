package todo

import (
	"os"
	"path/filepath"
	"time"
)

// tagProject checks whether a "project_info.toml" file exists in projectPath.
// If it does, it returns the base name of projectPath (i.e. the project folder name),
// otherwise it returns an empty string.
func tagProject(projectPath string) string {
	projectInfoPath := filepath.Join(projectPath, "project_info.toml")
	if _, err := os.Stat(projectInfoPath); err == nil {
		return filepath.Base(projectPath)
	}
	return ""
}

// otherwise it returns an empty string.
func tagWorkspace(projectPath string) string {
	wsInfoPath := filepath.Join(projectPath, "ws_info.toml")
	if _, err := os.Stat(wsInfoPath); err == nil {
		return filepath.Base(projectPath)
	}
	return ""
}

// buildTaskLine constructs and returns the complete task line including tags.
// It appends any non-empty tags (due date, project, and workspace) along with
// a creation date tag using the current date.
func buildTaskLine(description, dueDate, projectName, workspaceName string) string {
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
	// Append a created date tag using the current time.
	tags += " #created:" + time.Now().Format("2006-01-02")

	return "- [ ] " + description + tags
}
