package todo

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// TodoService defines the operations for managing TODO items.
type TodoService interface {
	AddTodo(description, dueDate string) error
	EditTodo(index int, newDescription, newDueDate string) error
	DeleteTodo(index int) error
	CompleteTodo(index int) error
	ListTodos() ([]Todo, error)
}

// FileTodoService is a concrete implementation of TodoService that uses a todo.md file.
type FileTodoService struct {
	todoFilePath string
}

// NewFileTodoService creates a new FileTodoService with the given todo file path.
func NewFileTodoService(todoFilePath string) *FileTodoService {
	return &FileTodoService{todoFilePath: todoFilePath}
}

// AddTodo creates a new task and appends it to the todo file.
func (s *FileTodoService) AddTodo(description, dueDate string) error {
	// Determine the project path from the todo file path.
	projectPath := filepath.Dir(s.todoFilePath)

	// Automatically determine project and workspace tags.
	projectName := tagProject(projectPath)
	workspaceName := tagWorkspace(projectPath)

	// Build the task line.
	taskLine := buildTaskLine(description, dueDate, projectName, workspaceName)

	// Read the existing file content.
	content, err := ReadFileContent(s.todoFilePath)
	if err != nil {
		return fmt.Errorf("failed to read '%s': %w", s.todoFilePath, err)
	}

	// Insert new task after the "# todo" header.
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
	newLines := append(lines[:insertIndex], append([]string{taskLine}, lines[insertIndex:]...)...)
	updatedContent := strings.Join(newLines, "\n")

	return WriteFileContent(s.todoFilePath, updatedContent)
}

// ListTodos returns the list of todos from the file.
func (s *FileTodoService) ListTodos() ([]Todo, error) {
	return LoadAllTodos(s.todoFilePath)
}

// EditTodo updates the description and/or due date of the task at the given index.
func (s *FileTodoService) EditTodo(index int, newDescription, newDueDate string) error {
	todos, err := LoadAllTodos(s.todoFilePath)
	if err != nil {
		return fmt.Errorf("error loading tasks: %w", err)
	}

	if index < 0 || index >= len(todos) {
		return fmt.Errorf("invalid task index")
	}

	selectedTask := todos[index]
	if newDescription != "" {
		selectedTask.Description = newDescription
	}
	if newDueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", newDueDate)
		if err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}
		selectedTask.DueDate = parsedDate
	}
	todos[index] = selectedTask

	return SaveTodos(s.todoFilePath, todos)
}

// DeleteTodo removes the task at the given index.
func (s *FileTodoService) DeleteTodo(index int) error {
	todos, err := LoadAllTodos(s.todoFilePath)
	if err != nil {
		return fmt.Errorf("error loading tasks: %w", err)
	}

	if index < 0 || index >= len(todos) {
		return fmt.Errorf("invalid task index")
	}

	todos = append(todos[:index], todos[index+1:]...)
	return SaveTodos(s.todoFilePath, todos)
}

// CompleteTodo marks the task at the given index as completed.
func (s *FileTodoService) CompleteTodo(index int) error {
	todos, err := LoadAllTodos(s.todoFilePath)
	if err != nil {
		return fmt.Errorf("error loading tasks: %w", err)
	}

	if index < 0 || index >= len(todos) {
		return fmt.Errorf("invalid task index")
	}

	// Ensure the task is not already completed.
	if !todos[index].CompletedDate.IsZero() {
		return fmt.Errorf("task already completed")
	}

	todos[index].CompletedDate = time.Now()
	return SaveTodos(s.todoFilePath, todos)
}
