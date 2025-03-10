package project

import (
	"fmt"

	"github.com/johnjallday/flow-workspace/internal/todo"
)

func AddTodoToProject(projectDir, description, dueDate string) error {
	if projectDir == "" || description == "" {
		return fmt.Errorf("project-dir and description are required")
	}

	todoFile := projectDir + "/todo.md"
	service := todo.NewFileTodoService(todoFile)

	if err := service.AddTodo(description, dueDate); err != nil {
		return fmt.Errorf("failed to add todo: %w", err)
	}

	return nil
}
