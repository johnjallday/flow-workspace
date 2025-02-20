package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/johnjallday/flow-workspace/internal/todo"
	"github.com/johnjallday/flow-workspace/internal/workspace"
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

// MigrateOldFinishedTodos aggregates all todos from every workspace under rootDir,
// filters the ones that are completed more than a week ago, and inserts them into the database.
func MigrateOldFinishedTodos(rootDir string, db *sql.DB) error {
	var aggregated []todo.Todo
	// Directories to skip at the root level.
	skip := map[string]bool{
		".config":     true,
		".spacedrive": true,
		".DS_Store":   true,
		".TagStudio":  true,
	}

	// List all entries in the root directory.
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return fmt.Errorf("failed to read root directory '%s': %w", rootDir, err)
	}

	// For each subdirectory (workspace), try to load its projects.toml.
	for _, entry := range entries {
		if !entry.IsDir() || skip[entry.Name()] {
			continue
		}
		workspaceDir := filepath.Join(rootDir, entry.Name())
		projectsTomlPath := filepath.Join(workspaceDir, "projects.toml")
		if _, err := os.Stat(projectsTomlPath); os.IsNotExist(err) {
			// Skip this workspace if no projects.toml exists.
			continue
		}

		// Use workspace.LoadProjectsToml to load the projects.
		projs, err := workspace.LoadProjectsToml(workspaceDir)
		if err != nil {
			fmt.Printf("Error loading projects from '%s': %v\n", workspaceDir, err)
			continue
		}

		// For each project entry, determine its directory and load its todo.md.
		for _, proj := range projs.Projects {
			var projectDir string
			if proj.Path != "" && proj.Path != "./" {
				if filepath.IsAbs(proj.Path) {
					projectDir = filepath.Clean(proj.Path)
				} else {
					projectDir = filepath.Join(workspaceDir, proj.Path)
				}
			} else {
				// Default: assume the project is in a folder named after the project.
				projectDir = filepath.Join(workspaceDir, proj.Name)
			}

			todoFile := filepath.Join(projectDir, "todo.md")
			if _, err := os.Stat(todoFile); os.IsNotExist(err) {
				// Skip if todo.md does not exist.
				continue
			}

			// Load todos from the file.
			todos, err := todo.LoadAllTodos(todoFile)
			if err != nil {
				fmt.Printf("Error loading todos from '%s': %v\n", todoFile, err)
				continue
			}

			aggregated = append(aggregated, todos...)
		}
	}

	// Filter finished todos older than 7 days.
	threshold := time.Hour * 24 * 7
	now := time.Now()
	var toMigrate []todo.Todo
	for _, t := range aggregated {
		if !t.CompletedDate.IsZero() && now.Sub(t.CompletedDate) > threshold {
			toMigrate = append(toMigrate, t)
		}
	}

	// Migrate each filtered todo to the database.
	for _, t := range toMigrate {
		if err := InsertTodo(db, t); err != nil {
			fmt.Printf("Error migrating todo '%s': %v\n", t.Description, err)
		} else {
			fmt.Printf("Migrated todo: %s\n", t.Description)
		}
	}

	return nil
}
