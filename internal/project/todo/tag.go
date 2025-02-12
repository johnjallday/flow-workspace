package todo

import (
	"fmt"
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

// tagWorkspace checks two folders above the projectPath for a "ws_info.toml" file.
// If the file is found, it returns the base name of that directory (i.e. the workspace name),
// otherwise it returns an empty string.
func tagWorkspace(projectPath string) string {
	// Get the parent folder of projectPath.
	workspacePath := filepath.Dir(projectPath)
	// Get the folder one level above parentPath (i.e. two levels above projectPath).
	fmt.Println(filepath.Base(workspacePath))
	if _, err := os.Stat(workspacePath); err == nil {
		return filepath.Base(workspacePath)
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
