package todo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Add creates a new task and appends it to the todo.md file located in projectPath.
func Add(projectPath string) error {
	todoFile := filepath.Join(projectPath, "todo.md")

	// If 'todo.md' does not exist, create it with a "# todo" header.
	if _, err := os.Stat(todoFile); os.IsNotExist(err) {
		initialContent := "# todo\n\n"
		if err := WriteFileContent(todoFile, initialContent); err != nil {
			return fmt.Errorf("failed to create 'todo.md': %w", err)
		}
	}

	// Clean up 'todo.md' to ensure proper organization.
	if err := CleanUpTodoFile(todoFile); err != nil {
		return fmt.Errorf("error cleaning up 'todo.md': %w", err)
	}

	reader := bufio.NewReader(os.Stdin)

	// Prompt user for task description.
	fmt.Print("Enter task description: ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read task description: %w", err)
	}
	description = strings.TrimSpace(description)
	if description == "" {
		return fmt.Errorf("task description cannot be empty")
	}

	// Prompt for due date.
	fmt.Print("Enter due date (YYYY-MM-DD) or leave empty: ")
	dueDateStr, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read due date: %w", err)
	}
	dueDateStr = strings.TrimSpace(dueDateStr)
	var dueDate string
	if dueDateStr != "" {
		if _, err := time.Parse("2006-01-02", dueDateStr); err != nil {
			return fmt.Errorf("invalid date format; please use YYYY-MM-DD")
		}
		dueDate = dueDateStr
	}

	// Automatically determine the project and workspace tags.
	projectName := tagProject(projectPath)
	workspaceName := tagWorkspace(projectPath)

	// Construct the task line with tags.
	taskLine := buildTaskLine(description, dueDate, projectName, workspaceName)

	// Read existing content.
	content, err := ReadFileContent(todoFile)
	if err != nil {
		return fmt.Errorf("failed to read '%s': %w", todoFile, err)
	}

	// Determine where to insert the new task (after the "# todo" header).
	lines := strings.Split(content, "\n")
	insertIndex := 0
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "# todo") {
			insertIndex = i + 1
			break
		}
	}
	for insertIndex < len(lines) && strings.TrimSpace(lines[insertIndex]) == "" {
		insertIndex++
	}

	// Insert the new task.
	newLines := append(lines[:insertIndex], append([]string{taskLine}, lines[insertIndex:]...)...)
	updatedContent := strings.Join(newLines, "\n")

	// Write updated content back to the file.
	if err := WriteFileContent(todoFile, updatedContent); err != nil {
		return fmt.Errorf("failed to write to '%s': %w", todoFile, err)
	}

	fmt.Println("Task added successfully.")
	return nil
}
