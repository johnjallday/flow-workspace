package startup

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	dbtodo "github.com/johnjallday/flow-workspace/internal/db/todo"
	_ "github.com/mattn/go-sqlite3"
)

// createConfig creates a "config" table and a default config if one doesn't exist.
func createConfig(db *sql.DB, username string) error {
	query := `
    CREATE TABLE IF NOT EXISTS config (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	// Check if a config entry exists.
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM config").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// Insert default configuration.
		_, err = db.Exec("INSERT INTO config (username) VALUES (?)", username)
		if err != nil {
			return err
		}
		log.Printf("Default config created with username: %s", username)
	}
	return nil
}

// createDB opens (and creates, if necessary) the SQLite database file,
// creates the todos table, and creates a default configuration if not already present.
func createDB(dbPath string, username string) {
	// Create the database file if it doesn't exist.
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		log.Printf("Database created at %s", dbPath)
	} else {
		log.Printf("Database already exists at %s", dbPath)
	}

	// Open a connection to the SQLite database.
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the todos table.
	if err := dbtodo.CreateTodoTable(db); err != nil {
		log.Fatalf("failed to create todos table: %v", err)
	}

	// Create the config table and default config if not exists.
	if err := createConfig(db, username); err != nil {
		log.Fatalf("failed to create config: %v", err)
	}
}

// StartDB determines the binary's directory, checks for an existing SQLite file,
// and either prompts for a username to create one or opens an existing file.
func StartDB() {
	// Determine the directory of the binary.
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	binDir := filepath.Dir(exePath)
	log.Println("Binary directory:", binDir)

	// Look for existing SQLite files with the pattern fw_*.sqlite.
	pattern := filepath.Join(binDir, "fw_*.sqlite")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}

	var dbPath string
	if len(matches) == 0 {
		// No existing SQLite file; prompt for username.
		var username string
		fmt.Print("Enter username: ")
		_, err := fmt.Scanln(&username)
		if err != nil || username == "" {
			log.Fatal("username is required")
		}
		dbPath = filepath.Join(binDir, "fw_"+username+".sqlite")
		createDB(dbPath, username)
		fmt.Printf("Welcome %s!\n", username)
	} else {
		// Use the first matching SQLite file.
		dbPath = matches[0]
		// Open the database and retrieve the username from the config table.
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var storedUsername string
		err = db.QueryRow("SELECT username FROM config LIMIT 1").Scan(&storedUsername)
		if err != nil {
			log.Fatalf("failed to retrieve username from config: %v", err)
		}
		fmt.Printf("Welcome %s!\n", storedUsername)
	}
}
