package db

import (
	"database/sql"
	"fmt"

	"github.com/johnjallday/flow-workspace/internal/todo"
)

// CreateTodoTable creates the "todos" table if it does not already exist.
func CreateTodoTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		completed_date DATETIME,
		created_date DATETIME NOT NULL,
		due_date DATETIME,
		project_name TEXT,
		workspace_name TEXT
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating todos table: %w", err)
	}
	fmt.Println("Todos table created or already exists.")
	return nil
}

// InsertTodo inserts a single todo entry into the database.
func InsertTodo(db *sql.DB, t todo.Todo) error {
	query := `
	INSERT INTO todos (description, completed_date, created_date, due_date, project_name, workspace_name)
	VALUES (?, ?, ?, ?, ?, ?);
	`
	// For date fields, if a zero value is present, we pass nil.
	var completedDate interface{}
	if !t.CompletedDate.IsZero() {
		completedDate = t.CompletedDate
	} else {
		completedDate = nil
	}

	var dueDate interface{}
	if !t.DueDate.IsZero() {
		dueDate = t.DueDate
	} else {
		dueDate = nil
	}

	_, err := db.Exec(query,
		t.Description,
		completedDate,
		t.CreatedDate,
		dueDate,
		t.ProjectName,
		t.WorkspaceName,
	)
	if err != nil {
		return fmt.Errorf("error inserting todo: %w", err)
	}
	return nil
}

// MoveTodosToDB reads the todo entries from the provided todoFile and inserts them into the database.
func MoveTodosToDB(db *sql.DB, todoFile string) error {
	// Load todos from the file using your existing function.
	todos, err := todo.LoadAllTodos(todoFile)
	if err != nil {
		return fmt.Errorf("error loading todos from file '%s': %w", todoFile, err)
	}

	// Loop over each todo and insert it into the DB.
	for _, t := range todos {
		if err := InsertTodo(db, t); err != nil {
			return fmt.Errorf("error inserting todo '%s': %w", t.Description, err)
		}
	}

	fmt.Printf("Successfully moved %d todos to the database.\n", len(todos))
	return nil
}
