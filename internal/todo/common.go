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
	CompletedDate time.Time // non-zero means the task is complete
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
		// If the task is completed but doesn't have a #completed tag, append it.
		if !strings.Contains(line, "#completed:") {
			line = line + fmt.Sprintf(" #completed:%s", time.Now().Format("2006-01-02"))
		}
		// Mark as completed.
		t.CompletedDate = time.Now()
		line = strings.TrimPrefix(line, "- [x]")
	} else if strings.HasPrefix(line, "- [ ]") {
		// Not completed, leave CompletedDate as zero value.
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
			case "completed":
				// If a #completed tag is present, use its date.
				d, err := time.Parse("2006-01-02", value)
				if err != nil {
					return t, fmt.Errorf("invalid completed_date format")
				}
				t.CompletedDate = d
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

// SaveTodos writes the given slice of todos back to the specified file.
// Each todo is formatted as a single line in the todo file.
func SaveTodos(filename string, todos []Todo) error {
	var lines []string

	// Iterate over todos and build a string for each one.
	for _, t := range todos {
		var lineBuilder strings.Builder

		// Mark as completed if CompletedDate is set.
		if t.CompletedDate.IsZero() {
			lineBuilder.WriteString("- [ ] ")
		} else {
			lineBuilder.WriteString("- [x] ")
		}

		// Write the description.
		lineBuilder.WriteString(t.Description)

		// Append tags for created, due, project, and workspace if available.
		if !t.CreatedDate.IsZero() {
			lineBuilder.WriteString(" #created:")
			lineBuilder.WriteString(t.CreatedDate.Format("2006-01-02"))
		}
		if !t.DueDate.IsZero() {
			lineBuilder.WriteString(" #due:")
			lineBuilder.WriteString(t.DueDate.Format("2006-01-02"))
		}
		if t.ProjectName != "" {
			lineBuilder.WriteString(" #project:")
			lineBuilder.WriteString(t.ProjectName)
		}
		if t.WorkspaceName != "" {
			lineBuilder.WriteString(" #workspace:")
			lineBuilder.WriteString(t.WorkspaceName)
		}

		lines = append(lines, lineBuilder.String())
	}

	// Join all lines into a single string with newline separation.
	content := strings.Join(lines, "\n")
	return WriteFileContent(filename, content)
}
