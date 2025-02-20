package todo

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	todocommon "github.com/johnjallday/flow-workspace/internal/todo"
)

// reloadAndDisplay reloads todos from the given file and displays them.
func reloadAndDisplay(todoFilePath string) {
	todos, err := todocommon.LoadAllTodos(todoFilePath)
	if err != nil {
		fmt.Printf("Error reloading tasks from '%s': %v\n", todoFilePath, err)
		return
	}
	todocommon.DisplayTodos(todos)
}

// StartTodoREPL is the interactive REPL for a single todo.md file.
func StartTodoREPL(todoFilePath string) {
	clearScreen()
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("TODO REPL started for file: %s\n", filepath.Base(todoFilePath))
	// Initial load and display of todos.
	printHelp()
	reloadAndDisplay(todoFilePath)

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
			clearScreen()
			printHelp()
		case "list":
			clearScreen()
			printHelp()
			reloadAndDisplay(todoFilePath)
		case "add":
			clearScreen()
			Add(filepath.Dir(todoFilePath))
			printHelp()
			reloadAndDisplay(todoFilePath)
		case "complete":
			clearScreen()
			Complete(todoFilePath, reader)
			printHelp()
			reloadAndDisplay(todoFilePath)
		case "delete":
			clearScreen()
			Delete(todoFilePath, reader)
			printHelp()
			reloadAndDisplay(todoFilePath)
		case "edit":
			clearScreen()
			Edit(todoFilePath, reader)
			printHelp()
			reloadAndDisplay(todoFilePath)
		case "weekly":
			clearScreen()
			fmt.Println("Running weekly review...")
			// (Add any weekly review functionality here)
			printHelp()
			reloadAndDisplay(todoFilePath)
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func printHelp() {
	fmt.Println(`Available commands (TODO REPL):
  list              - List all tasks
  add               - Add a new task
  complete          - Mark a task as completed
  delete            - Delete a task
  edit              - Edit a task
  weekly            - Run the weekly review for this TODO file
  help              - Show this help message
  exit              - Exit the TODO REPL`)
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("clear") // Linux
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls") // Windows
	case "darwin":
		cmd = exec.Command("clear") // macOS
	default:
		fmt.Println("CLS for", runtime.GOOS, "not implemented")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
