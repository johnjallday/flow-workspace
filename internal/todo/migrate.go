package todo

import (
	"database/sql"
	"fmt"
	"time"
)

// InsertTodo inserts a single todo entry into the database.
func InsertTodo(db *sql.DB, t Todo) error {
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

func MigrateFinishedTodos(todoPath string, db *sql.DB) error {
	// Load all todos from the file.
	todos, err := LoadAllTodos(todoPath)
	if err != nil {
		return fmt.Errorf("error loading todos: %w", err)
	}

	var remainingTodos []Todo
	var migratedCount int

	// Calculate the cutoff date (7 days ago)
	cutoffDate := time.Now().AddDate(0, 0, -7)

	fmt.Println("Checking for completed todos older than 7 days...")

	for _, t := range todos {
		fmt.Println("Checking:", t.Description)

		// Ensure CompletedDate is valid before checking
		if t.CompletedDate.IsZero() {
			fmt.Println("Skipping, no completion date set")
			remainingTodos = append(remainingTodos, t)
			continue
		}

		if t.CompletedDate.Before(cutoffDate) {
			fmt.Println("completed")
			fmt.Println("Migrating:", t.Description)
			// add logic for insert to db
			InsertTodo(db, t)
			migratedCount++
		}
		// If any todos were migrated, update the todo file.
		if migratedCount > 0 {
			if err := SaveTodos(todoPath, remainingTodos); err != nil {
				return fmt.Errorf("error saving todos: %w", err)
			}
			fmt.Printf("Migrated %d finished todo(s) to the database and removed them from the todo file.\n", migratedCount)
		} else {
			fmt.Println("No finished todos older than 7 days found for migration.")
		}
	}

	return nil
}
