package todo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StartTodoREPL starts an interactive REPL for a given todo file.
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
			todos, err = LoadAllTodos(todoFilePath)
			if err != nil {
				fmt.Printf("Error loading tasks from '%s': %v\n", todoFilePath, err)
				continue
			}
			DisplayTodos(todos)
		case "add":
			if err := addTask(filepath.Dir(todoFilePath), reader); err != nil {
				fmt.Printf("Error adding task: %v\n", err)
			}
		// (You would similarly refactor the complete, delete, and edit commands.)
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func printTodoHelp() {
	fmt.Println(`Available commands (TODO REPL):
  list      - List all tasks
  add       - Add a new task
  complete  - Mark a task as completed
  delete    - Delete a task
  edit      - Edit a task
  help      - Show this help message
  exit      - Exit the TODO REPL`)
}

// addTask prompts for details and adds a new task.
func addTask(projectPath string, reader *bufio.Reader) error {
	todoFile := filepath.Join(projectPath, "todo.md")
	if _, err := os.Stat(todoFile); os.IsNotExist(err) {
		initialContent := "# todo\n\n"
		if err := WriteFileContent(todoFile, initialContent); err != nil {
			return fmt.Errorf("failed to create 'todo.md': %w", err)
		}
	}
	// (Optional: call a cleanup function to organize the file.)

	fmt.Print("Enter task description: ")
	desc, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read task description: %w", err)
	}
	desc = strings.TrimSpace(desc)
	if desc == "" {
		return fmt.Errorf("task description cannot be empty")
	}

	fmt.Print("Enter due date (YYYY-MM-DD) or leave empty: ")
	dueStr, _ := reader.ReadString('\n')
	dueStr = strings.TrimSpace(dueStr)
	if dueStr != "" {
		if _, err := time.Parse("2006-01-02", dueStr); err != nil {
			return fmt.Errorf("invalid date format; please use YYYY-MM-DD")
		}
	}

	projectName := TagProject(projectPath)
	workspaceName := TagWorkspace(projectPath)
	taskLine := BuildTaskLine(desc, dueStr, projectName, workspaceName)

	// Insert the new task after the "# todo" header.
	content, err := ReadFileContent(todoFile)
	if err != nil {
		return fmt.Errorf("failed to read '%s': %w", todoFile, err)
	}
	lines := strings.Split(content, "\n")
	insertIndex := 0
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "# todo") {
			insertIndex = i + 1
			break
		}
	}
	for insertIndex < len(lines) && strings.TrimSpace(lines[insertIndex]) == "" {
		insertIndex++
	}
	newLines := append(lines[:insertIndex], append([]string{taskLine}, lines[insertIndex:]...)...)
	updatedContent := strings.Join(newLines, "\n")

	if err := WriteFileContent(todoFile, updatedContent); err != nil {
		return fmt.Errorf("failed to write to '%s': %w", todoFile, err)
	}

	fmt.Println("Task added successfully.")
	return nil
}
