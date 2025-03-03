package db

import (
	"database/sql"
	"fmt"
)

// InitDB connects to the SQLite database at dbPath.
func InitDB(dbPath string) (*sql.DB, error) {
	// Adjust the driver and connection string if you are using a different database.
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Verify the connection is working.
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Database connection successful.")
	return conn, nil
}

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

// MigrateFinishedTodos migrates todos that are completed and more than 7 days old into the database.
// After migration, the migrated entries are removed from the todo file.
