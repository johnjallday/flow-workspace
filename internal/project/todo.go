// internal/project/todo.go

package project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Todo represents a single TODO task within todo.md.
type Todo struct {
	Description   string
	Completed     bool
	CreatedDate   time.Time
	DueDate       time.Time
	ProjectName   string
	WorkspaceName string
}

// Regular expression to extract tags in the format #key:value
var tagRegex = regexp.MustCompile(`#(\w+):([^\s#]+)`)

// HandleTodoInput manages the input of a new task within the specified workspace.
func HandleTodoInput(workspacePath string) {
	todoFile := filepath.Join(workspacePath, "todo.md")

	// If 'todo.md' does not exist, create it with a '# todo' header
	if _, err := os.Stat(todoFile); os.IsNotExist(err) {
		initialContent := "# todo\n\n"
		if err := WriteFileContent(todoFile, initialContent); err != nil {
			fmt.Printf("Failed to create 'todo.md': %v\n", err)
			return
		}
	}

	// Clean up 'todo.md' to ensure proper organization
	if err := CleanUpTodoFile(todoFile); err != nil {
		fmt.Println("Error cleaning up 'todo.md':", err)
		return
	}

	// Initialize REPL reader for user input
	reader := bufio.NewReader(os.Stdin)

	// Prompt user for task details
	fmt.Print("Enter task description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)
	if description == "" {
		fmt.Println("Task description cannot be empty.")
		return
	}

	fmt.Print("Enter due date (YYYY-MM-DD) or leave empty: ")
	dueDateStr, _ := reader.ReadString('\n')
	dueDateStr = strings.TrimSpace(dueDateStr)
	var dueDate string
	if dueDateStr != "" {
		// Validate date format
		_, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			fmt.Println("Invalid date format. Please use YYYY-MM-DD.")
			return
		}
		dueDate = dueDateStr
	}

	fmt.Print("Enter project name or leave empty: ")
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)

	fmt.Print("Enter workspace name or leave empty: ")
	workspaceName, _ := reader.ReadString('\n')
	workspaceName = strings.TrimSpace(workspaceName)

	// Construct the task line with tags
	taskLine := "- [ ] " + description
	if dueDate != "" {
		taskLine += " #due:" + dueDate
	}
	if projectName != "" {
		taskLine += " #project:" + projectName
	}
	if workspaceName != "" {
		taskLine += " #workspace:" + workspaceName
	}
	taskLine += " #created:" + time.Now().Format("2006-01-02")

	// Read existing content
	content, err := ReadFileContent(todoFile)
	if err != nil {
		fmt.Println("Error reading 'todo.md':", err)
		return
	}

	// Split content into lines
	lines := strings.Split(content, "\n")

	// Find the index after the '# todo' header
	insertIndex := 0
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "# todo") {
			insertIndex = i + 1
			break
		}
	}

	// Ensure there's at least one blank line after the header
	for insertIndex < len(lines) && strings.TrimSpace(lines[insertIndex]) == "" {
		insertIndex++
	}

	// Insert the new task at the determined index
	lines = append(lines[:insertIndex], append([]string{taskLine}, lines[insertIndex:]...)...)

	// Reconstruct the content
	updatedContent := strings.Join(lines, "\n")

	// Write back to 'todo.md'
	if err := WriteFileContent(todoFile, updatedContent); err != nil {
		fmt.Println("Error writing to 'todo.md':", err)
		return
	}

	fmt.Println("Task added successfully.")
}

// HandleCompleteTodo marks a task as completed based on its description and metadata.
func HandleCompleteTodo(todoFile string, reader *bufio.Reader) {
	todos, err := LoadAllTodos(todoFile)
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No tasks available to complete.")
		return
	}

	// Display all incomplete tasks with indices
	fmt.Println("Incomplete Tasks:")
	var incompleteTasks []Todo
	for _, todo := range todos {
		if !todo.Completed {
			incompleteTasks = append(incompleteTasks, todo)
		}
	}

	if len(incompleteTasks) == 0 {
		fmt.Println("All tasks are already completed.")
		return
	}

	for i, todo := range incompleteTasks {
		dueDate := "No due date"
		if !todo.DueDate.IsZero() {
			dueDate = todo.DueDate.Format("2006-01-02")
		}
		project := "No project"
		if todo.ProjectName != "" {
			project = todo.ProjectName
		}
		workspace := "No workspace"
		if todo.WorkspaceName != "" {
			workspace = todo.WorkspaceName
		}
		fmt.Printf("%d. %s (Due: %s, Project: %s, Workspace: %s)\n",
			i+1, todo.Description, dueDate, project, workspace)
	}

	fmt.Print("Enter the number of the task to mark as completed: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Convert input to integer
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(incompleteTasks) {
		fmt.Println("Invalid task number.")
		return
	}

	// Find and mark the selected task as completed
	selectedTask := incompleteTasks[index-1]
	for i, todo := range todos {
		if todo.Description == selectedTask.Description &&
			todo.CreatedDate.Equal(selectedTask.CreatedDate) &&
			todo.DueDate.Equal(selectedTask.DueDate) &&
			todo.ProjectName == selectedTask.ProjectName &&
			todo.WorkspaceName == selectedTask.WorkspaceName &&
			!todo.Completed {
			todos[i].Completed = true
			break
		}
	}

	// Reconstruct the content
	var updatedContent string
	updatedContent += "# todo\n\n"
	for _, todo := range todos {
		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}
		taskLine := fmt.Sprintf("- %s %s", status, todo.Description)
		if !todo.CreatedDate.IsZero() {
			taskLine += fmt.Sprintf(" #created:%s", todo.CreatedDate.Format("2006-01-02"))
		}
		if !todo.DueDate.IsZero() {
			taskLine += fmt.Sprintf(" #due:%s", todo.DueDate.Format("2006-01-02"))
		}
		if todo.ProjectName != "" {
			taskLine += fmt.Sprintf(" #project:%s", todo.ProjectName)
		}
		if todo.WorkspaceName != "" {
			taskLine += fmt.Sprintf(" #workspace:%s", todo.WorkspaceName)
		}
		taskLine += "\n"
		updatedContent += taskLine
	}

	// Backup 'todo.md' before writing changes
	if err := BackupFile(todoFile); err != nil {
		fmt.Printf("Failed to backup 'todo.md': %v\n", err)
	}

	// Write back to 'todo.md'
	if err := WriteFileContent(todoFile, updatedContent); err != nil {
		fmt.Println("Error writing to 'todo.md':", err)
		return
	}

	fmt.Println("Task marked as completed successfully.")
}

// HandleDeleteTodo allows users to delete a task based on its number.
func HandleDeleteTodo(todoFile string, reader *bufio.Reader) {
	todos, err := LoadAllTodos(todoFile)
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No tasks available to delete.")
		return
	}

	// Display all tasks with indices
	fmt.Println("All Tasks:")
	for i, todo := range todos {
		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}
		fmt.Printf("%d. %s %s\n", i+1, status, todo.Description)
	}

	fmt.Print("Enter the number of the task to delete: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Convert input to integer
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(todos) {
		fmt.Println("Invalid task number.")
		return
	}

	// Remove the selected task
	todos = append(todos[:index-1], todos[index:]...)

	// Reconstruct the content
	var updatedContent string
	updatedContent += "# todo\n\n"
	for _, todo := range todos {
		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}
		taskLine := fmt.Sprintf("- %s %s", status, todo.Description)
		if !todo.CreatedDate.IsZero() {
			taskLine += fmt.Sprintf(" #created:%s", todo.CreatedDate.Format("2006-01-02"))
		}
		if !todo.DueDate.IsZero() {
			taskLine += fmt.Sprintf(" #due:%s", todo.DueDate.Format("2006-01-02"))
		}
		if todo.ProjectName != "" {
			taskLine += fmt.Sprintf(" #project:%s", todo.ProjectName)
		}
		if todo.WorkspaceName != "" {
			taskLine += fmt.Sprintf(" #workspace:%s", todo.WorkspaceName)
		}
		taskLine += "\n"
		updatedContent += taskLine
	}

	// Backup 'todo.md' before writing changes
	if err := BackupFile(todoFile); err != nil {
		fmt.Printf("Failed to backup 'todo.md': %v\n", err)
	}

	// Write back to 'todo.md'
	if err := WriteFileContent(todoFile, updatedContent); err != nil {
		fmt.Println("Error writing to 'todo.md':", err)
		return
	}

	fmt.Println("Task deleted successfully.")
}

// HandleFilterTodoByProject displays tasks filtered by a specific project name.
func HandleFilterTodoByProject(todoFile string, projectName string) {
	todos, err := LoadAllTodos(todoFile)
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No tasks available.")
		return
	}

	filteredTasks := []Todo{}
	for _, todo := range todos {
		if strings.EqualFold(todo.ProjectName, projectName) {
			filteredTasks = append(filteredTasks, todo)
		}
	}

	if len(filteredTasks) == 0 {
		fmt.Printf("No tasks found for project '%s'.\n", projectName)
		return
	}

	fmt.Printf("Tasks for project '%s':\n", projectName)
	for i, todo := range filteredTasks {
		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}
		fmt.Printf("%d. %s %s\n", i+1, status, todo.Description)
	}
}

// HandleEditTodo allows users to edit a task based on its number.
func HandleEditTodo(todoFile string, reader *bufio.Reader) {
	todos, err := LoadAllTodos(todoFile)
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No tasks available to edit.")
		return
	}

	// Display all tasks with indices
	fmt.Println("All Tasks:")
	for i, todo := range todos {
		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}
		fmt.Printf("%d. %s %s\n", i+1, status, todo.Description)
	}

	fmt.Print("Enter the number of the task to edit: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Convert input to integer
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(todos) {
		fmt.Println("Invalid task number.")
		return
	}

	selectedTask := todos[index-1]

	// Prompt for new description
	fmt.Printf("Enter new description (leave empty to keep '%s'): ", selectedTask.Description)
	newDescription, _ := reader.ReadString('\n')
	newDescription = strings.TrimSpace(newDescription)
	if newDescription == "" {
		newDescription = selectedTask.Description
	}

	// Prompt for new due date
	fmt.Printf("Enter new due date (YYYY-MM-DD) or leave empty to keep '%s': ", selectedTask.DueDate.Format("2006-01-02"))
	newDueDateStr, _ := reader.ReadString('\n')
	newDueDateStr = strings.TrimSpace(newDueDateStr)
	var newDueDate time.Time
	if newDueDateStr != "" {
		newDueDate, err = time.Parse("2006-01-02", newDueDateStr)
		if err != nil {
			fmt.Println("Invalid date format. Please use YYYY-MM-DD.")
			return
		}
	} else {
		newDueDate = selectedTask.DueDate
	}

	// Update task
	todos[index-1].Description = newDescription
	todos[index-1].DueDate = newDueDate

	// Reconstruct the content
	var updatedContent string
	updatedContent += "# todo\n\n"
	for _, todo := range todos {
		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}
		taskLine := fmt.Sprintf("- %s %s", status, todo.Description)
		if !todo.CreatedDate.IsZero() {
			taskLine += fmt.Sprintf(" #created:%s", todo.CreatedDate.Format("2006-01-02"))
		}
		if !todo.DueDate.IsZero() {
			taskLine += fmt.Sprintf(" #due:%s", todo.DueDate.Format("2006-01-02"))
		}
		if todo.ProjectName != "" {
			taskLine += fmt.Sprintf(" #project:%s", todo.ProjectName)
		}
		if todo.WorkspaceName != "" {
			taskLine += fmt.Sprintf(" #workspace:%s", todo.WorkspaceName)
		}
		taskLine += "\n"
		updatedContent += taskLine
	}

	// Backup 'todo.md' before writing changes
	if err := BackupFile(todoFile); err != nil {
		fmt.Printf("Failed to backup 'todo.md': %v\n", err)
	}

	// Write back to 'todo.md'
	if err := WriteFileContent(todoFile, updatedContent); err != nil {
		fmt.Println("Error writing to 'todo.md':", err)
		return
	}

	fmt.Println("Task edited successfully.")
}

// HandleTodoView manages the 'todo' view flow within the specified workspace.
func HandleTodoView(workspacePath string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nTodo View Options:")
	fmt.Println("1. View All TODOs (Root Level)")
	fmt.Println("2. View Workspace TODOs")
	fmt.Println("3. View Project TODOs")
	fmt.Print("Select an option (1, 2, or 3): ")

	option, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading option:", err)
		return
	}
	option = strings.TrimSpace(option)

	switch option {
	case "1":
		// Root Level
		todos, err := GatherTodosRoot(workspacePath)
		if err != nil {
			fmt.Println("Error gathering root TODOs:", err)
			return
		}
		DisplayTodos(todos)
	case "2":
		// Workspace Level
		todos, err := GatherTodosWorkspace(workspacePath)
		if err != nil {
			fmt.Println("Error gathering workspace TODOs:", err)
			return
		}
		DisplayTodos(todos)
	case "3":
		// Project Level
		todos, err := GatherTodosProject(workspacePath)
		if err != nil {
			fmt.Println("Error gathering project TODOs:", err)
			return
		}
		DisplayTodos(todos)
	default:
		fmt.Println("Invalid option. Please select 1, 2, or 3.")
	}
}

// DisplayTodos prints the list of TODOs in a formatted manner.
func DisplayTodos(todos []Todo) {
	if len(todos) == 0 {
		fmt.Println("No TODOs found.")
		return
	}

	fmt.Println("\n--- TODO List ---")
	for i, todo := range todos {
		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}
		fmt.Printf("%d. %s %s\n", i+1, status, todo.Description)
		if !todo.DueDate.IsZero() {
			fmt.Printf("   Due Date: %s\n", todo.DueDate.Format("2006-01-02"))
		}
		if !todo.CreatedDate.IsZero() {
			fmt.Printf("   Created: %s\n", todo.CreatedDate.Format("2006-01-02"))
		}
		if todo.ProjectName != "" {
			fmt.Printf("   Project: %s\n", todo.ProjectName)
		}
		if todo.WorkspaceName != "" {
			fmt.Printf("   Workspace: %s\n", todo.WorkspaceName)
		}
	}
	fmt.Println("-----------------\n")
}

// CleanUpTodoFile ensures that all tasks are under the '# todo' header.
// It moves any tasks found before the header to below it and removes duplicate headers.
func CleanUpTodoFile(filename string) error {
	// Read the existing content of 'todo.md'
	content, err := ReadFileContent(filename)
	if err != nil {
		return fmt.Errorf("error reading '%s': %w", filename, err)
	}

	lines := strings.Split(content, "\n")
	var tasks []string
	var headers []int // To keep track of header positions

	// Iterate through lines to collect tasks and identify headers
	for idx, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# todo") {
			headers = append(headers, idx)
			continue // Skip the header line
		}
		if strings.HasPrefix(trimmed, "- [ ]") || strings.HasPrefix(trimmed, "- [x]") {
			tasks = append(tasks, line)
		}
	}

	// If no '# todo' header exists, add it at the top
	if len(headers) == 0 {
		lines = append([]string{"# todo", ""}, lines...)
	} else {
		// Remove all headers except the first one
		if len(headers) > 1 {
			// Remove duplicate headers starting from the second one
			for i := 1; i < len(headers); i++ {
				headerIdx := headers[i] - i // Adjust index due to previous removals
				if headerIdx >= 0 && headerIdx < len(lines) {
					lines = append(lines[:headerIdx], lines[headerIdx+1:]...)
				}
			}
		}

		// Remove any tasks that might still be before the header
		// Re-extract tasks after removing duplicate headers
		linesAfterHeader := lines[0 : headers[0]+1]
		var newTasks []string
		for _, line := range linesAfterHeader {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "- [ ]") || strings.HasPrefix(trimmed, "- [x]") {
				newTasks = append(newTasks, line)
			}
		}
		lines = linesAfterHeader
		tasks = append(newTasks, tasks...)
	}

	// Reconstruct the content: header followed by tasks
	var updatedLines []string
	updatedLines = append(updatedLines, "# todo", "") // Ensure header and a blank line

	updatedLines = append(updatedLines, tasks...)

	// Optionally, add other content below tasks if necessary
	// For example, if you have sections after tasks, handle them here

	// Join all lines into a single string
	updatedContent := strings.Join(updatedLines, "\n")

	// Write the cleaned content back to 'todo.md'
	if err := WriteFileContent(filename, updatedContent); err != nil {
		return fmt.Errorf("error writing to '%s': %w", filename, err)
	}

	return nil
}

// BackupFile creates a timestamped backup of the specified file.
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

// GatherTodosRoot scans all workspaces and their projects to collect TODOs.
func GatherTodosRoot(rootDir string) ([]Todo, error) {
	var todos []Todo

	// Walk through the root directory
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Failed to access path %s: %v\n", path, err)
			return nil // Continue walking
		}

		// Check if the file is named 'todo.md'
		if !info.IsDir() && info.Name() == "todo.md" {
			fmt.Printf("Found todo.md: %s\n", path)
			t, err := LoadTodos(path)
			if err != nil {
				fmt.Printf("Failed to load TODOs from '%s': %v\n", path, err)
				return nil // Continue walking
			}
			todos = append(todos, t...)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the path '%s': %w", rootDir, err)
	}

	if len(todos) == 0 {
		fmt.Println("No todo.md files found in the root directory.")
	}

	return todos, nil
}

// GatherTodosWorkspace scans all projects within a workspace to collect TODOs.
func GatherTodosWorkspace(workspacePath string) ([]Todo, error) {
	var todos []Todo

	projectsTomlPath := filepath.Join(workspacePath, "projects.toml")
	projs, err := LoadProjectsInfo(projectsTomlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load 'projects.toml': %w", err)
	}

	for _, proj := range projs.Projects {
		todoPath := filepath.Join(workspacePath, proj.Path, "todo.md")
		if _, err := os.Stat(todoPath); os.IsNotExist(err) {
			fmt.Printf("todo.md not found for project '%s' at '%s'\n", proj.Name, todoPath)
			continue
		}

		fmt.Printf("Found todo.md: %s\n", todoPath)
		t, err := LoadTodos(todoPath)
		if err != nil {
			fmt.Printf("Failed to load TODOs from '%s': %v\n", todoPath, err)
			continue
		}
		todos = append(todos, t...)
	}

	if len(todos) == 0 {
		fmt.Println("No todo.md files found in the selected workspace.")
	}

	return todos, nil
}

// GatherTodosProject retrieves TODOs from a specific project's todo.md.
func GatherTodosProject(projectPath string) ([]Todo, error) {
	var todos []Todo

	todoPath := filepath.Join(projectPath, "todo.md")
	if _, err := os.Stat(todoPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("todo.md does not exist at '%s'", todoPath)
	}

	fmt.Printf("Found todo.md: %s\n", todoPath)
	t, err := LoadTodos(todoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load TODOs from '%s': %w", todoPath, err)
	}

	todos = append(todos, t...)

	if len(todos) == 0 {
		fmt.Println("No TODOs found in the project.")
	}

	return todos, nil
}

// LoadTodos reads and parses a todo.md file into a slice of Todo structs.
func LoadTodos(filename string) ([]Todo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open '%s': %w", filename, err)
	}
	defer file.Close()

	var todos []Todo
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and headers
		}
		if strings.HasPrefix(line, "- [ ]") || strings.HasPrefix(line, "- [x]") {
			todo, err := parseTodoLine(line)
			if err != nil {
				fmt.Printf("Skipping invalid task line: %s\n", line)
				continue
			}
			todos = append(todos, todo)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading '%s': %w", filename, err)
	}

	return todos, nil
}

// LoadAllTodos parses all tasks from todo.md into a slice of Todo structs.
func LoadAllTodos(filename string) ([]Todo, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var todos []Todo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and headers
		}
		todo, err := parseTodo(line)
		if err != nil {
			fmt.Printf("Skipping invalid task line: %s\n", line)
			continue
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// parseTodo parses a single task line into a Todo struct.
func parseTodo(line string) (Todo, error) {
	var todo Todo

	// Check if the task is completed
	if strings.HasPrefix(line, "- [x]") {
		todo.Completed = true
		line = strings.TrimPrefix(line, "- [x]")
	} else if strings.HasPrefix(line, "- [ ]") {
		todo.Completed = false
		line = strings.TrimPrefix(line, "- [ ]")
	} else {
		return todo, fmt.Errorf("invalid task format")
	}

	line = strings.TrimSpace(line)

	// Extract tags
	tags := tagRegex.FindAllStringSubmatch(line, -1)

	for _, tag := range tags {
		if len(tag) == 3 {
			key := strings.ToLower(tag[1])
			value := tag[2]
			switch key {
			case "created":
				createdDate, err := time.Parse("2006-01-02", value)
				if err != nil {
					return todo, fmt.Errorf("invalid created_date format")
				}
				todo.CreatedDate = createdDate
			case "due":
				dueDate, err := time.Parse("2006-01-02", value)
				if err != nil {
					return todo, fmt.Errorf("invalid due_date format")
				}
				todo.DueDate = dueDate
			case "project":
				todo.ProjectName = value
			case "workspace":
				todo.WorkspaceName = value
			}
		}
	}

	// Remove tags from description
	description := tagRegex.ReplaceAllString(line, "")
	todo.Description = strings.TrimSpace(description)

	return todo, nil
}

// WriteFileContent writes the provided content to the specified file.
func WriteFileContent(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

// ReadFileContent reads the entire content of the specified file and returns it as a string.
func ReadFileContent(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// parseTodoLine parses a single line from todo.md into a Todo struct.
func parseTodoLine(line string) (Todo, error) {
	var todo Todo

	// Check if the task is completed
	if strings.HasPrefix(line, "- [x]") {
		todo.Completed = true
		line = strings.TrimPrefix(line, "- [x]")
	} else if strings.HasPrefix(line, "- [ ]") {
		todo.Completed = false
		line = strings.TrimPrefix(line, "- [ ]")
	} else {
		return todo, fmt.Errorf("invalid task format")
	}

	line = strings.TrimSpace(line)

	// Extract tags
	tags := tagRegex.FindAllStringSubmatch(line, -1)
	for _, tag := range tags {
		if len(tag) == 3 {
			key := strings.ToLower(tag[1])
			value := tag[2]
			switch key {
			case "created":
				createdDate, err := time.Parse("2006-01-02", value)
				if err != nil {
					return todo, fmt.Errorf("invalid created_date format")
				}
				todo.CreatedDate = createdDate
			case "due":
				dueDate, err := time.Parse("2006-01-02", value)
				if err != nil {
					return todo, fmt.Errorf("invalid due_date format")
				}
				todo.DueDate = dueDate
			case "project":
				todo.ProjectName = value
			case "workspace":
				todo.WorkspaceName = value
			}
		}
	}

	// Remove tags from description
	description := tagRegex.ReplaceAllString(line, "")
	todo.Description = strings.TrimSpace(description)

	return todo, nil
}
