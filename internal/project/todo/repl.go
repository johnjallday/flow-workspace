package todo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ----- Minimal definitions duplicated from workspace package ----- //

// Project is a minimal duplicate of the project type used in projects.toml.
type Project struct {
	Name string `toml:"name"`
	Path string `toml:"path"`
	// Add other fields if needed.
}

// StartTodoREPL is the interactive REPL for a single todo.md file.
func StartTodoREPL(todoFilePath string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("TODO REPL started for file: %s\n", filepath.Base(todoFilePath))
	printTodoHelp()
	todos, err := LoadAllTodos(todoFilePath)
	if err != nil {
		fmt.Printf("Error loading tasks from '%s': %v\n", todoFilePath, err)
	}
	DisplayTodos(todos)

	for {
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
		case "help":
			printTodoHelp()
		case "list":
			todos, err := LoadAllTodos(todoFilePath)
			if err != nil {
				fmt.Printf("Error loading tasks from '%s': %v\n", todoFilePath, err)
				continue
			}
			DisplayTodos(todos)
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
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func printTodoHelp() {
	fmt.Println(`Available commands (TODO REPL):
  list              - List all tasks
  add               - Add a new task
  complete          - Mark a task as completed
  delete            - Delete a task
  edit              - Edit a task
  weekly            - Run the weekly review for this TODO file
  help              - Show this help message
  exit              - Exit the TODO REPL
`)
}

// -------------------------------------------------------------------------
// The following functions wrap StartTodoREPL for different scopes.
// They duplicate (or independently implement) the logic so that package
// todo does not import workspace (thus avoiding an import cycle).
// -------------------------------------------------------------------------

// StartProjectTodoREPL launches the TODO REPL for a project directory.
// It ensures that a todo.md exists in the given project directory.
func StartProjectTodoREPL(projectDir string) {
	todoFile := filepath.Join(projectDir, "todo.md")
	if _, err := os.Stat(todoFile); os.IsNotExist(err) {
		fmt.Printf("todo.md not found in '%s'. Creating one...\n", projectDir)
		initialContent := "# todo\n\n"
		if err := os.WriteFile(todoFile, []byte(initialContent), 0644); err != nil {
			fmt.Printf("Failed to create todo.md: %v\n", err)
			return
		}
	}
	StartTodoREPL(todoFile)
}
