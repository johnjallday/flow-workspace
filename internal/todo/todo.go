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
