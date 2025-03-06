package todo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	db "github.com/johnjallday/flow-workspace/internal/db/todo"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// reloadAndDisplay reloads todos from the given file and displays them.
func reloadAndDisplay(todoFilePath string) {
	todos, err := LoadAllTodos(todoFilePath)
	if err != nil {
		fmt.Printf("Error reloading tasks from '%s': %v\n", todoFilePath, err)
		return
	}
	PrintTodos(todos)
}

// StartTodoREPL is the interactive REPL for a single todo.md file.
func StartTodoREPL(dbPath string, todoFilePath string) {
	reader := bufio.NewReader(os.Stdin)

	mydb, err := db.InitDB(dbPath)
	// Use our local initDB function to connect to the database.
	if err != nil {
		fmt.Println("Error connecting to db:", err)
		return
	}

	MigrateFinishedTodos(todoFilePath, mydb)
	// Initial load and display of todos.
	for {

		clearScreen()

		fmt.Println("dbPath:", dbPath)

		// Migrate finished todos using the provided function from the db package.
		printHelp()
		reloadAndDisplay(todoFilePath)
		fmt.Printf("\n[todo:%s] >> ", filepath.Base(todoFilePath))
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		command := strings.ToLower(parts[0])

		switch command {
		case "exit":
			fmt.Println("Exiting TODO REPL. Goodbye!")
			return
		case "add":
			Add(filepath.Dir(todoFilePath))
		case "complete":
			Complete(todoFilePath, reader)
		case "delete":
			Delete(todoFilePath, reader)
		case "edit":
			Edit(todoFilePath, reader)
		case "weekly":
			fmt.Println("Running weekly review...")
			// (Add any weekly review functionality here)
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func printHelp() {
	fmt.Println(`Available commands (TODO REPL):
  add               - Add a new task
  complete          - Mark a task as completed
  delete            - Delete a task
  edit              - Edit a task
  weekly            - Run the weekly review for this TODO file
  exit              - Exit the TODO REPL`)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
