package todo

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Todo represents a single task.
type Todo struct {
	Description   string
	Completed     bool
	CreatedDate   time.Time
	DueDate       time.Time
	ProjectName   string
	WorkspaceName string
}

// tagRegex is used to extract tags from a task line.
var tagRegex = regexp.MustCompile(`#(\w+):([^\s#]+)`)

// LoadAllTodos reads a todo file (e.g. "todo.md") and parses its tasks.
func LoadAllTodos(filename string) ([]Todo, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var todos []Todo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip blank lines or comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		t, err := parseTodo(line)
		if err != nil {
			fmt.Printf("Skipping invalid task line: %s\n", line)
			continue
		}
		todos = append(todos, t)
	}

	return todos, nil
}

// parseTodo converts a single todo line into a Todo struct.
func parseTodo(line string) (Todo, error) {
	var t Todo

	if strings.HasPrefix(line, "- [x]") {
		t.Completed = true
		line = strings.TrimPrefix(line, "- [x]")
	} else if strings.HasPrefix(line, "- [ ]") {
		t.Completed = false
		line = strings.TrimPrefix(line, "- [ ]")
	} else {
		return t, fmt.Errorf("invalid task format")
	}

	line = strings.TrimSpace(line)
	tags := tagRegex.FindAllStringSubmatch(line, -1)
	for _, tag := range tags {
		if len(tag) == 3 {
			key := strings.ToLower(tag[1])
			value := tag[2]
			switch key {
			case "created":
				d, err := time.Parse("2006-01-02", value)
				if err != nil {
					return t, fmt.Errorf("invalid created_date format")
				}
				t.CreatedDate = d
			case "due":
				d, err := time.Parse("2006-01-02", value)
				if err != nil {
					return t, fmt.Errorf("invalid due_date format")
				}
				t.DueDate = d
			case "project":
				t.ProjectName = value
			case "workspace":
				t.WorkspaceName = value
			}
		}
	}

	// Remove the tags from the description.
	desc := tagRegex.ReplaceAllString(line, "")
	t.Description = strings.TrimSpace(desc)
	return t, nil
}

// WriteFileContent writes content to a file with 0644 permissions.
func WriteFileContent(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

// ReadFileContent reads and returns the content of a file.
func ReadFileContent(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// BackupFile creates a timestamped backup of the given file.
func BackupFile(filename string) error {
	timestamp := time.Now().Format("20060102_150405")
	backupFilename := fmt.Sprintf("%s_backup_%s", filename, timestamp)
	input, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}
	if err := os.WriteFile(backupFilename, input, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}
	return nil
}
