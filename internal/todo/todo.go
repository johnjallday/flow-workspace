package todo

import "time"

// Todo represents a single task.
type Todo struct {
	Description   string
	CompletedDate time.Time // non-zero means the task is complete
	CreatedDate   time.Time
	DueDate       time.Time
	ProjectName   string
	WorkspaceName string
}

// Manager defines the operations you expect on todos.
type Manager interface {
	LoadAllTodos() ([]Todo, error)
	SaveTodos(todos []Todo) error
	AddTodo(todo Todo) error
	DeleteTodo(index int) error
}
